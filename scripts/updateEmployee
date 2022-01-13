#!/bin/python3

import requests
import sys


def print_usage():
    print("Usage: ./updateEmployee.py <first> <last> <target> <value>")
    print("Targets -> [wage, email, phone]")
    sys.exit()

# check commandline arguments
if len(sys.argv) != 5:
    print_usage()

base = "https://pcc-time-bot.herokuapp.com/updateEmployee"
first = sys.argv[1]
last = sys.argv[2]
target = sys.argv[3]
value = sys.argv[4]

response = requests.get(f"{base}/{first}/{last}/{target}/{value}")

print(response.text)