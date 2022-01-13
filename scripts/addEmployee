#!/bin/python3

import requests
import sys


def print_usage():
    print("Usage: ./addEmployee <first> <last> <wage> <email> <phone>")
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

print(response.text)
