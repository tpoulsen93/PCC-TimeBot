from sqlalchemy import create_engine

import sys, os

sys.path.insert(0, '/Users/tpoulsen/Code/PCC-TimeBot')
from databaseAccess import *
from timeCard import *


url = os.environ['DATABASE_URL']
url = url.replace("postgres", "postgresql")
engine = create_engine(url)


bob = TimeCard(3, "2022-01-24", "2022-01-30")
#bob.days['2022-01-25'] = 10
bob.add_hours("2022-01-25", 5)
print(bob.to_string())
