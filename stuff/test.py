from twilio.rest import Client

from sqlalchemy import create_engine, insert, text, ForeignKey

import sys, os

sys.path.insert(0, '/Users/tpoulsen/Code/PCC-TimeBot')
from src.databaseAccess import *


url = os.environ['DATABASE_URL']
url = url.replace("postgres", "postgresql")
engine = create_engine(url)

result = get_employee_id("taylor", "poulsen")
print(result)
id = result.scalar()
print(id)

result = get_employee_phone(id)
print(result)
print(result.scalar())

