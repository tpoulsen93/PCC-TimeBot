from src.timeCardGeneration import buildTimeCards
import datetime, os

from datetime import timedelta

from sqlalchemy import MetaData, Table, Column, String, Integer, Float, Date
from sqlalchemy import create_engine, insert, text, ForeignKey



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
    Column('phone', String, unique=True),
    Column('email', String),
    Column('supervisor', String) # <first> <last>
)

payroll = Table(
    'payroll', meta,
    Column('id', ForeignKey('employees.id')),
    Column('transaction_id', Integer, autoincrement=True, primary_key=True),
    Column('time', Float),
    Column('date', Date),
    Column('message', String)
)

meta.create_all(engine)


def duplicate_submission(id):
    # heroku uses utc time and we need mountain time so this is my hacky conversion
    today = (datetime.datetime.today() - timedelta(hours=7)).date()

    stmt = text("SELECT time FROM payroll WHERE id = :i AND date = :d")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = id, d = today).first()
    if not result:
        return False
    return result[0]


# submit hours for an employee
def submit_time(id, time, msg) -> str:
    # heroku uses utc time and we need mountain time so this is my hacky conversion
    today = (datetime.datetime.today() - timedelta(hours=7)).date()

    dupe = duplicate_submission(id)
    if not dupe:
        stmt = insert(payroll).values(id = id, time = time, date = today, message = msg)
        result = f"Submitted {str(time)} hours"
    else:
        stmt = text("UPDATE payroll SET time = :t, message = :m \
            WHERE id = :i AND date = :d")
        result = f"Updated submission from {str(dupe)} to {str(time)} hours"
    with engine.connect() as conn:
        conn.execute(stmt, t = time, m = msg, i = id, d = today)
    return result


# return the id number of the employee if they exist
def get_employee_id(first: str, last: str):
    stmt = text("SELECT id FROM employees \
        WHERE first_name = :f AND last_name = :l")
    with engine.connect() as conn:
        result = conn.execute(stmt, f = first, l = last).first()
    return result


# return the first and last name of the employee associated with the id
def get_employee_name(id) -> str:
    stmt = text("SELECT first_name, last_name FROM employees \
        WHERE id = :i")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = id)
    
    first = result[0]
    last = result[1]
    return f"{first} {last}"


# add a new employee to the table
def add_employee(first, last, wage, email = "", phone = ""):
    stmt = insert(employees).values(
        first_name = first.lower(),
        last_name = last.lower(),
        wage = wage,
        phone = phone if phone != "" else None,
        email = email if email != "" else None
    )
    with engine.connect() as conn:
        conn.execute(stmt)


# update wage, email, or phone for an employee
def update_employee(first, last, target, value):
    id = get_employee_id(first, last)
    stmt = text("UPDATE employees SET :t = :v WHERE id = :i")
    with engine.connect() as conn:
        conn.execute(stmt, t = target, v = value, i = id)


# get all the information for the indicated dates from the database and send them
# to the parser to build the json object of time cards
def get_time_cards(start, end):
    stmt = text("SELECT id, time, date FROM payroll \
        WHERE date >= :s AND date <= :e GROUP BY id ORDER BY date")
    with engine.connect() as conn:
        result = conn.execute(stmt, s = start, e = end)

    return buildTimeCards(result, start, end)


