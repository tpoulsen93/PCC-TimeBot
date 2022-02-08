from sqlalchemy import MetaData, Table, Column, String, Integer, Float, Date
from sqlalchemy import create_engine, insert, text, ForeignKey, TIMESTAMP
from datetime import timedelta, datetime
import os


url = os.environ['DATABASE_URL']
url = url.replace("postgres", "postgresql") 
engine = create_engine(url)
meta = MetaData()

employees = Table(
    'employees', meta, 
    Column('id', Integer, autoincrement=True, primary_key=True),
    Column('first_name', String),
    Column('last_name', String),
    Column('wage', Float),
    Column('phone', String),
    Column('email', String),
    Column('supervisor_id', Integer),
    Column('timestamp', TIMESTAMP(timezone='america/boise'), nullable=False, default=datetime.now())
)

payroll = Table(
    'payroll', meta,
    Column('id', ForeignKey('employees.id')),
    Column('transaction_id', Integer, autoincrement=True, primary_key=True),
    Column('time', Float),
    Column('date', Date),
    Column('message', String),
    Column('timestamp', TIMESTAMP(timezone='america/boise'), nullable=False, default=datetime.now())
)

meta.create_all(engine)



# get the supervisor id of the employee
def get_super_id(employee_id) -> int:
    stmt = text("SELECT supervisor_id FROM employees \
        WHERE id = :i")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = employee_id)
    return result if not result else int(result.scalar())


# return the id number of the employee if they exist
def get_employee_id(first: str, last: str) -> int:
    stmt = text("SELECT id FROM employees \
        WHERE first_name = :f AND last_name = :l")
    with engine.connect() as conn:
        result = conn.execute(stmt, f = first.lower(), l = last.lower()).first()
    return result if not result else int(result[0])


# return the first and last name of the employee associated with the id
def get_employee_name(id) -> str:
    stmt = text("SELECT first_name, last_name FROM employees \
        WHERE id = :i")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = id)
    
    first = result[0]
    last = result[1]
    return f"{first} {last}"


def get_employee_phone(id) -> str:
    stmt = text("SELECT phone FROM employees \
        WHERE id = :i")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = id)
    return result if not result else result.scalar()


def get_employee_email(id) -> str:
    stmt = text("SELECT email FROM employees \
        WHERE id = :i")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = id)
    return result if not result else result.scalar()


def get_employee(id):
    stmt = text("SELECT * FROM employees WHERE id = :i")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = id)
    return result


def duplicate_submission(id, date):
    stmt = text("SELECT time FROM payroll WHERE id = :i AND date = :d")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = id, d = date).first()
    if not result:
        return False
    return result[0]


# submit hours for an employee
def submit_time(id, time, msg) -> str:
    # heroku uses utc time and we need mountain time so this is my hacky conversion
    today = (datetime.today() - timedelta(hours=7)).date()

    dupe = duplicate_submission(id, today)
    if not dupe:
        stmt = insert(payroll).values(id = id, time = time, date = today, message = msg)
        result = f"Submitted hours: {time:g}"
    else:
        stmt = text("UPDATE payroll SET time = :t, message = :m WHERE id = :i AND date = :d")
        result = f"Updated hours: {dupe:g} to {time:g}"
    with engine.connect() as conn:
        conn.execute(stmt, t = time, m = msg, i = id, d = today)
    return result
    

# add hours for an employee on specific date
def add_time(first, last, date, time):
    today = (datetime.today() - timedelta(hours=7)).date()

    id = get_employee_id(first, last)
    dupe = duplicate_submission(id, date)
    if not dupe:
        msg = f"Added {first} {last} manually on {today}"
        stmt = insert(payroll).values(id = id, time = time, date = date, message = msg)
        result = f"Submitted {time:g} hours for {first.title()} {last.title()} on {date}"
    else:
        msg = f"Updated {first} {last} manually on {today}"
        stmt = text("UPDATE payroll SET time = :t, message = :m WHERE id = :i AND date = :d")
        result = f"Updated submission for {first.title()} {last.title()} from {dupe:g} to {time:g} hours on {date}"
    with engine.connect() as conn:
        conn.execute(stmt, t = time, m = msg, i = id, d = date)
    return result


def add_employee(first, last, wage, email = "", phone = "", super_first = "", super_last = ""):
    if super_first != "" and super_last != "":
        super_id = get_employee_id(super_first, super_last)
        if not super_id:
            return "Error. Failed to find supervisor id"
    else:
        super_id = None

    stmt = insert(employees).values(
        first_name = first.lower(),
        last_name = last.lower(),
        supervisor_id = super_id,
        wage = wage,
        phone = phone if phone != "" else None,
        email = email.lower() if email != "" else None
    )
    with engine.connect() as conn:
        conn.execute(stmt)
    return f"{first.title()} {last.title()} was successfully added"


# update wage, email, or phone for an employee
def update_employee(first, last, target, value):
    id = get_employee_id(first, last)
    stmt = text(f"UPDATE employees SET {target} = :v WHERE id = :i")
    with engine.connect() as conn:
        conn.execute(stmt, v = value, i = id)
    return f"{first.title()} {last.title()}'s {target} was changed to {value}"


# get all the information for the indicated dates from the database
def get_time_cards(start, end):
    stmt = text(f"SELECT id, time, date FROM payroll \
        WHERE date >= '{start}'::date AND date <= '{end}'::date ORDER BY id")
    with engine.connect() as conn:
        result = conn.execute(stmt)
    return result
