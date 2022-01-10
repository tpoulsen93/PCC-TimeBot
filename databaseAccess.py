import json
import os
import psycopg2
from typing_extensions import Required
from sqlalchemy.sql.expression import false, select, update
from sqlalchemy.sql.schema import ForeignKey
from sqlalchemy.sql.sqltypes import Date
from sqlalchemy import create_engine, MetaData, Table, Column, String, Integer, Float, Date
from sqlalchemy.orm import session, sessionmaker
from sqlalchemy import create_engine, insert
from datetime import date


# engine = create_engine('sqlite:///database.db')
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
    Column('time', Float),
    Column('draw', Float),
    Column('date', Date),
    Column('msg', String)
)

meta.create_all(engine)


def insert_time(id, time, msg):
    stmt = insert(payroll).values(id=id, time=time, date=date.today(), msg=msg)
    with engine.connect() as conn:
        conn.execute(stmt)

# def insert_draw(id, amount, msg):
#     stmt = insert(payroll).values(id=id, draw=amount, date=date.today(), msg=msg)

#     with engine.connect() as conn:
#         conn.execute(stmt)


# return true if the employee exists in the database, else return false
def get_employee_id(first: str, last: str):
    stmt = employees.select().where(
        employees.c.first_name == first and employees.c.last_name == last
    )
    with engine.connect() as conn:
        result = conn.execute(stmt)
        print(f"result: {result}")
        print(f"result[0] {result[0]}")
    return result[0]

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
        

# if __name__ == "__main__":
    # insert_employee(first_name="Taylor", last_name="Poulsen", wage="20.00",  \
    #                 phone_number="432-276-1331", email="DanielMBogden@gmail.com")
    # print("Successfull")

    # Session = sessionmaker(bind = engine)
    # session = Session()
    # members = session.query(employees).all()

    # with Session() as session:
    #     t = session.query(employees).filter(employees.first_name=="Taylor").all()
    #     print(t)
    
