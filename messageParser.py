
from datetime import timedelta, date
from sqlalchemy.sql.sqltypes import String
import databaseAccess
import exceptions

time_error = "Error. Time formatted incorrectly."
draw_error = "Error. Draw formatted incorrectly."


# calculate hours for the day and return them
def calculate_time(start: str, end: str, less: str, more: str) -> float:
    # length of times should be 6 or 7  -->  00:00xm or 0:00xm
    if len(start) < 6 or len(start) > 7 or len(end) < 6 or len(end) > 7:
        raise exceptions.TimeFormatException

    # build clock-in time
    arr = start.split(":")

    # validate the start hours and meridiem
    startHours = int(arr[0])
    if startHours < 1 or startHours > 12:
        raise exceptions.HoursException

    if arr[1].endswith("am"):
        arr[1] = arr[1].replace("am", "")
    elif arr[1].endswith("pm"):
        startHours += 12
        arr[1] = arr[1].replace("pm", "")
    else: # neither am nor pm detected
        raise exceptions.MeridiemException

    # validate the start minutes
    startMinutes = int(arr[1])
    if startMinutes < 0 or startMinutes > 59:
        raise exceptions.MinutesException

    startTime = timedelta(hours=startHours, minutes=startMinutes)


    # build clock-out time
    arr = end.split(":")

    # validate end hours and meridiem
    endHours = int(arr[0])
    if endHours < 1 or endHours > 12:
        raise exceptions.HoursException

    if arr[1].endswith("am"):
        arr[1] = arr[1].replace("am", "")
    elif arr[1].endswith("pm"):
        endHours += 12
        arr[1] = arr[1].replace("pm", "")
    else: # neither am nor pm detected
        raise exceptions.MeridiemException

    # validate end minutes
    endMinutes = int(arr[1])
    if endMinutes < 0 or endMinutes > 59:
        raise exceptions.MinutesException

    endTime = timedelta(hours=endHours, minutes=endMinutes)

    # check for any other exceptions in the message
    if endTime < startTime:
        raise exceptions.IllegalTimeException
    
    try:
        subtract = float(less)
    except:
        raise exceptions.LunchException

    try:
        add = float(more) if more != "" else 0
    except:
        raise exceptions.ExtraException

    # compute hours for the day and return as a float rounded to 2 decimal places
    hours = endTime - startTime - timedelta(hours=subtract) + timedelta(hours=add)
    return round(hours / timedelta(hours=1), 2)



def process_draw(message: str) -> str:
    mess = message.split()
    if len(mess) < 4:
        return f"{draw_error} Too few parameters"
    if len(mess) > 4:
        return f"{draw_error} Too many parameters"

    # get the employee id or return False if they don't exist
    employeeId = databaseAccess.get_employee_id(mess[1].lower(), mess[2].lower())
    if not employeeId:
        return "Error. Employee not found."
    employeeId = int(employeeId[0])

    if "$" in mess[3]:
        mess[3].replace("$", "")

    try: # cast the dollar amount to a float    
        draw = float(mess[3])
    except:
        raise exceptions.DrawException

    # add the draw to the database and return the message to be texted back
    databaseAccess.insert_draw(employeeId, draw, message)
    return "A $" + str(draw) + " draw was submitted for " + mess[1].title() + " " + mess[2].title()



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

    # add the hours to the database and return the message to be texted back
    submission = databaseAccess.insert_time(employeeId, time, message)
    return f"{submission} for {mess[1].title()} {mess[2].title()} for {date.today()}"



def process_message(message: str):
    mess = message.lower()

    # handle a time submission
    if mess.startswith("time") or mess.startswith("hours"):
        if "help" in mess:
            return "Help"
        return process_time(message)

    # handle a draw submission
    elif message.lower().startswith("draw"):
        return process_draw(message)

    # ignore the message, it isn't meant for us
    else:
        return False
