import src.databaseAccess as databaseAccess
from src.exceptions import *
import datetime, sys, os

from twilio.rest import Client
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
    employee_id = databaseAccess.get_employee_id(mess[1], mess[2])
    if not employee_id:
        return "Error. Employee not found."

    # get the start time, end time, break time, and extra time
    start = mess[3]
    end = mess[4]
    less = mess[5]
    more = mess[6] if len(mess) == 7 else ""

    # calculate the hours for the day or return an error message
    try:
        time = calculate_time(start, end, less, more)
    except HoursException:
        return f"{time_error} Hours spot is wrong."
    except MeridiemException:
        return f"{time_error} Meridiem is wrong. (am/pm)"
    except MinutesException:
        return f"{time_error} Minutes spot is wrong."
    except IllegalTimeException:
        return f"{time_error} End time is before start time..."
    except LunchException:
        return "Error. Subtracted hours formatted incorrectly."
    except ExtraException:
        return "Error. Additional hours formatted incorrectly."
    except TimeFormatException:
        return time_error

    # heroku uses utc time and we need mountain time so this is my hacky conversion
    today = (datetime.datetime.today() - timedelta(hours=7)).date()

    # add the hours to the database
    submission = databaseAccess.submit_time(employee_id, time, message)
    result = f"{today}\n{submission} for {mess[1].title()} {mess[2].title()}"

    # send the submission to the supervisor and myself
    supervisor_id = databaseAccess.get_super_id(employee_id)
    if supervisor_id != None:
        supervisor_phone = databaseAccess.get_employee_phone(supervisor_id)
        if not supervisor_phone:
            return "Error. Supervisor phone not found."

    tp_id = databaseAccess.get_employee_id('taylor', 'poulsen')
    if not tp_id:
        return "Error. TP not found..."
    tp_phone = databaseAccess.get_employee_phone(tp_id)
    if not tp_phone:
        return "Error. TP phone not found."

    twilio = os.environ['TWILIO_PHONE']
    client = Client(
        os.environ['TWILIO_ACCOUNT_SID'],
        os.environ['TWILIO_AUTH_TOKEN']
    )
    client.messages.create(
        from_=f"+1{twilio}",
        to=f"+1{tp_phone}",
        body=result
    )
    if supervisor_id:
        client.messages.create(
            from_=f"+1{twilio}",
            to=f"+1{supervisor_phone}",
            body=result
        )

    return result


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
