import src.exceptions as exceptions
from datetime import timedelta


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
