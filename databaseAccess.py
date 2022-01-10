import os

from sqlalchemy import MetaData, Table, Column, String, Integer, Float, Date
from sqlalchemy import create_engine, insert, text, update, ForeignKey
from datetime import date

from sqlalchemy.sql.expression import delete, true


url = os.environ['DATABASE_URL']
url = url.replace("postgres", "postgresql") # sqlalchemy deprecated postgres so this is my hacky solution...
engine = create_engine(url)

meta = MetaData()

employees = Table(
    'employees', meta, 
    Column('id', Integer, autoincrement=True, primary_key=True),
    Column('first_name', String),
    Column('last_name', String),
    Column('wage', Float),
    Column('phone_number', String, unique=True),
    Column('email', String)
)

payroll = Table(
    'payroll', meta,
    Column('id', ForeignKey('employees.id')),
    Column('transaction_id', Integer, autoincrement=True, primary_key=True),
    Column('time', Float),
    Column('draw', Float),
    Column('date', Date),
    Column('msg', String)
)

meta.create_all(engine)


def duplicate_submission(id):
    stmt = text("SELECT payroll.time FROM payroll WHERE \
        payroll.id LIKE :i AND payroll.date LIKE :d")
    with engine.connect() as conn:
        result = conn.execute(stmt, i = id, d = date.today()).first()
    if not result:
        return False
    return result[0]


# submit hours for an employee
def insert_time(xid, time, msg) -> str:
    dupe = duplicate_submission(id)
    if not dupe:
        stmt = insert(payroll).values(xid=id, time=time, date=date.today(), msg=msg)
        result = f"Submitted {str(time)} hours"
    else:
        stmt = update(payroll).values(time = time).where(id is Integer(xid) and date is date.today())
        result = f"Updated hours submission from {str(dupe)} to {str(time)}"
    with engine.connect() as conn:
        conn.execute(stmt)
    return result
    


# submit a draw for an employee
def insert_draw(id, amount, msg) -> str:
    stmt = insert(payroll).values(id=id, draw=amount, date=date.today(), msg=msg)
    with engine.connect() as conn:
        conn.execute(stmt)


# return true if the employee exists in the database, else return false
def get_employee_id(first: str, last: str):
    stmt = text("SELECT employees.id FROM employees WHERE \
        employees.first_name LIKE :f AND employees.last_name LIKE :l")
    with engine.connect() as conn:
        result = conn.execute(stmt, f = first, l = last).first()
    return result


# add a new employee to the table
def insert_employee(first_name, last_name, wage, email="", phone=""):
    stmt = insert(employees).values(
        first_name = first_name,
        last_name = last_name,
        wage = wage,
        phone_number = phone if phone != "" else None,
        email = email if email != "" else None
    )

    with engine.connect() as conn:
        conn.execute(stmt)


# update wage, email, or phone for an employee
def update_employee(first_name, last_name, wage, email="", phone=""):
    id = get_employee_id(first_name, last_name)
    with engine.connect() as conn:
        if wage != 0:
            stmt = update(employees).values(wage = wage).where(id = id)
            conn.execute(stmt)
        if email != "":
            stmt = update(employees).values(email = email).where(id = id)
            conn.execute(stmt)
        if phone != "":
            stmt = update(employees).values(phone_number = phone).where(id = id)
            conn.execute(stmt)    
