from src.exceptions import *
from datetime import timedelta


# calculate hours for the day and return them
def calculate_time(start: str, end: str, less: str, more: str) -> float:
    # length of times should be 6 or 7  -->  00:00xm or 0:00xm
    if len(start) < 3 or len(start) > 7 or len(end) < 3 or len(end) > 7:
        raise TimeFormatException

    # build clock-in time
    if len(start > 4):  # 9:00am
        arr = start.split(":")

        # validate the start hours and meridiem
        startHours = int(arr[0])
        if startHours < 1 or startHours > 12:
            raise HoursException

        if arr[1].endswith("am"):
            arr[1] = arr[1].replace("am", "")
        elif arr[1].endswith("pm"):
            startHours += 12
            arr[1] = arr[1].replace("pm", "")
        else: # neither am nor pm detected
            raise MeridiemException

        # validate the start minutes
        startMinutes = int(arr[1])
        if startMinutes < 0 or startMinutes > 59:
            raise MinutesException

    else:   # 9am
        startMinutes = 0
        if start.endswith("am"):
            start = start.replace("am", "")
            startHours = int(start)
        elif start.endswith("pm"):
            start = start.replace("pm", "")
            startHours = int(start) + 12
        else: # neither am nor pm detected
            raise MeridiemException

        if startHours < 1 or startHours > 12:
            raise HoursException

    startTime = timedelta(hours=startHours, minutes=startMinutes)

            

    # build clock-out time
    if len(end) > 4:    # 9:00am
        arr = end.split(":")

        # validate end hours and meridiem
        endHours = int(arr[0])
        if endHours < 1 or endHours > 12:
            raise HoursException

        if arr[1].endswith("am"):
            arr[1] = arr[1].replace("am", "")
        elif arr[1].endswith("pm"):
            endHours += 12
            arr[1] = arr[1].replace("pm", "")
        else: # neither am nor pm detected
            raise MeridiemException

        # validate end minutes
        endMinutes = int(arr[1])
        if endMinutes < 0 or endMinutes > 59:
            raise MinutesException

    else:   # 9am
        endMinutes = 0
        if end.endswith("am"):
            end = end.replace("am", "")
            endHours = int(end)
        elif end.endswith("pm"):
            end = end.replace("pm", "")
            endHours = int(end) + 12
        else: # neither am nor pm detected
            raise MeridiemException

        if endHours < 1 or endHours > 12:
            raise HoursException
    
    endTime = timedelta(hours=endHours, minutes=endMinutes)



    # check for any other exceptions in the message
    if endTime < startTime:
        raise IllegalTimeException
    
    try:
        subtract = float(less)
    except:
        raise LunchException

    try:
        add = float(more) if more != "" else 0
    except:
        raise ExtraException

    # compute hours for the day and return as a float rounded to 2 decimal places
    hours = endTime - startTime - timedelta(hours=subtract) + timedelta(hours=add)
    return round(hours / timedelta(hours=1), 2)
