import os
import datetime

from datetime import timedelta

from sqlalchemy import MetaData, Table, Column, String, Integer, Float, Date
from sqlalchemy import create_engine, insert, text, update, delete, true, ForeignKey



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
    Column('phone', String, unique=True),
    Column('email', String)
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
def insert_time(id, time, msg) -> str:
    # heroku uses utc time and we need mountain time so this is my hacky conversion
    today = (datetime.datetime.today() - timedelta(hours=7)).date()

    dupe = duplicate_submission(id)
    if not dupe:
        stmt = insert(payroll).values(id = id, time = time, date = today, message = msg)
        result = f"Submitted {str(time)} hours"
    else:
        stmt = text("UPDATE payroll SET time = :t, msg = :m WHERE id = :i AND date = :d")
        result = f"Updated submission from {str(dupe)} to {str(time)}"
    with engine.connect() as conn:
        conn.execute(stmt, t = time, m = msg, i = id, d = today)
    return result


# return true if the employee exists in the database, else return false
def get_employee_id(first: str, last: str):
    stmt = text("SELECT employees.id FROM employees WHERE \
        employees.first_name LIKE :f AND employees.last_name LIKE :l")
    with engine.connect() as conn:
        result = conn.execute(stmt, f = first, l = last).first()
    return result


# add a new employee to the table
def insert_employee(first_name, last_name, wage, email = "", phone = ""):
    stmt = insert(employees).values(
        first_name = first_name,
        last_name = last_name,
        wage = wage,
        phone = phone if phone != "" else None,
        email = email if email != "" else None
    )

    with engine.connect() as conn:
        conn.execute(stmt)


# update wage, email, or phone for an employee
def update_employee(first_name, last_name, wage, email = "", phone = ""):
    id = get_employee_id(first_name, last_name)
    with engine.connect() as conn:
        if wage != 0:
            stmt = update(employees).values(wage = wage).where(id = id)
            conn.execute(stmt)
        if email != "":
            stmt = update(employees).values(email = email).where(id = id)
            conn.execute(stmt)
        if phone != "":
            stmt = update(employees).values(phone = phone).where(id = id)
            conn.execute(stmt)    
