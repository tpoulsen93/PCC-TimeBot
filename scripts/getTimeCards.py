#!/bin/python3

import requests
import sys


def print_usage():
    print("Usage: ./getTimeCards <start date> <end date>")
    print("Date format:    YYYY-MM-DD")
    sys.exit()

# check commandline arguments
if len(sys.argv) != 3:
    print_usage()

base = "https://pcc-time-bot.herokuapp.com/getTimeCards"
start = sys.argv[1]
end = sys.argv[2]

response = requests.get(f"{base}/{start}/{end}")

print(response.text)