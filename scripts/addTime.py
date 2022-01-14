#!/bin/python3

from datetime import date
import requests
import sys


def print_usage():
    print("Usage: ./addHours <first> <last> <date> <time>")
    print("Date format:    YYYY-MM-DD")
    sys.exit()

# check commandline arguments
if len(sys.argv) != 6:
    print_usage()

base = "https://pcc-time-bot.herokuapp.com/addHours"
first = sys.argv[1]
last = sys.argv[2]
date = sys.argv[3]
time = sys.argv[4]

response = requests.get(f"{base}/{first}/{last}/{date}/{time}")

print(response.text)
