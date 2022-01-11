#!/bin/python3

import requests
import pprint
import json
import sys


def print_usage():
    print("Usage: add_employee.py <first> <last> <wage> <email> <phone>")
    sys.exit()

# check commandline arguments
if len(sys.argv) != 6:
    print_usage()

base = "https://pcc-time-bot.herokuapp.com/addEmployee"
first = sys.argv[1]
last = sys.argv[2]
wage = sys.argv[3]
email = sys.argv[4]
phone = sys.argv[5]

response = requests.get(f"{base}/{first}/{last}/{wage}/{email}/{phone}")

# pprint.pprint(json.loads(response.content))
print(response.text)
