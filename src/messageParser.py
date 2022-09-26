import src.databaseAccess as da
import src.twilioActions as ta
from src.exceptions import *
import datetime
from src.timeCalc import calculate_time
from datetime import timedelta



time_error = "Error. Time formatted incorrectly."


def process_time(message: str, fro: str) -> str:
    mess = message.split()
    if len(mess) < 6:
        return f"{time_error} Too few parameters."
    if len(mess) > 7:
        return f"{time_error} Too many parameters."


    # get the employee id
    employee_id = da.get_employee_id(mess[1], mess[2])
    if not employee_id:
        return "Error. Employee not found."

    # get the start time, end time, break time, and extra time
    start = mess[3]
    end = mess[4]
    less = mess[5]
    more = mess[6] if len(mess) == 7 else "0"

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
    submission = da.submit_time(employee_id, time, message)
    if more == "":
        more = 0
    result =  f"{today}\n"
    result += f"{mess[1].title()} {mess[2].title()}\n"
    result += f"Start: {start}\n"
    result += f"End: {end}\n"
    if less > 0:
        result += f"Lunch hours: {less}\n"
    if more > 0:
        result += f"Extra hours: {more}\n"
    result += f"{submission}"
    confirmation = ta.confirm_submission(employee_id, result, fro)

    if confirmation:
        return result
    else:
        return f"{result}\n{confirmation}"


def process_message(message: str, fro: str):
    mess = message.lower()

    # handle a time submission
    if mess.startswith("time") or mess.startswith("hours"):
        if "help" in mess:
            return "Help"
        return process_time(message, fro)

    # ignore the message, it isn't meant for us
    else:
        return False
