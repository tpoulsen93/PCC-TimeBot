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
