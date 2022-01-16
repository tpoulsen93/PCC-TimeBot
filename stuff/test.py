from twilio.rest import Client

from sqlalchemy import create_engine

import sys, os

sys.path.insert(0, '/Users/tpoulsen/Code/PCC-TimeBot')
from src.databaseAccess import *


url = os.environ['DATABASE_URL']
url = url.replace("postgres", "postgresql")
engine = create_engine(url)

id = get_employee_id("taylor", "poulsen")
print(f"my id: {id}")

phone = get_employee_phone(id)
print(f"my phone: {phone}")


supervisor_id = get_super_id('1')
if not supervisor_id:
    print("Error. Supervisor not found.")
print(f"super id: {supervisor_id}")
