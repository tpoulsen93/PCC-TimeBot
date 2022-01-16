from twilio.rest import Client

from sqlalchemy import create_engine

import sys, os

sys.path.insert(0, '/Users/tpoulsen/Code/PCC-TimeBot')
from src.databaseAccess import *


url = os.environ['DATABASE_URL']
url = url.replace("postgres", "postgresql")
engine = create_engine(url)

result = get_employee_id("taylor", "poulsen")
id = result[0]
print(f"my id: {id}")

result = get_employee_phone(id).scalar()
#print(result.scalar())
print(f"my phone: {result}")


supervisor_id = databaseAccess.get_super_id('1')
if not supervisor_id:
    return "Error. Supervisor not found."
supervisor_id = int(supervisor_id[0])
print(supervisor_id)
