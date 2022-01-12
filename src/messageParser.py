import src.databaseAccess as databaseAccess
import src.exceptions as exceptions
import datetime

from src.timeCalc import calculate_time
from datetime import timedelta

time_error = "Error. Time formatted incorrectly."


def process_time(message: str) -> str:
    mess = message.split()
    if len(mess) < 6:
        return f"{time_error} Too few parameters."
    if len(mess) > 7:
        return f"{time_error} Too many parameters."

    # get the employee id
    employeeId = databaseAccess.get_employee_id(mess[1].lower(), mess[2].lower())
    if not employeeId:
        return "Error. Employee not found."
    employeeId = int(employeeId[0])

    # get the start time, end time, break time, and extra time
    start = mess[3]
    end = mess[4]
    less = mess[5]
    more = mess[6] if len(mess) == 7 else ""

    # calculate the hours for the day or return an error message
    try:
        time = calculate_time(start, end, less, more)
    except exceptions.HoursException:
        return f"{time_error} Hours spot is wrong."
    except exceptions.MeridiemException:
        return f"{time_error} Meridiem is wrong. (am/pm)"
    except exceptions.MinutesException:
        return f"{time_error} Minutes spot is wrong."
    except exceptions.IllegalTimeException:
        return f"{time_error} End time is before start time..."
    except exceptions.LunchException:
        return "Error. Subtracted hours formatted incorrectly."
    except exceptions.ExtraException:
        return "Error. Additional hours formatted incorrectly."
    except exceptions.TimeFormatException:
        return time_error

    # heroku uses utc time and we need mountain time so this is my hacky conversion
    today = (datetime.datetime.today() - timedelta(hours=7)).date()

    # add the hours to the database and return the message to be texted back
    submission = databaseAccess.submit_time(employeeId, time, message)
    return f"{submission} for {mess[1].title()} {mess[2].title()} for {today}"



def process_message(message: str):
    mess = message.lower()

    # handle a time submission
    if mess.startswith("time") or mess.startswith("hours"):
        if "help" in mess:
            return "Help"
        return process_time(message)

    # ignore the message, it isn't meant for us
    else:
        return False
