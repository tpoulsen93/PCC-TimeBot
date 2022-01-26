from src.exceptions import *
from datetime import timedelta


def build_time_delta(time: str) -> timedelta:
    if ":" in time:  # 9:00am
        if len(time) < 6 or len(time) > 7:
            raise TimeException

        arr = time.lower().split(":")

        # validate the hours and meridiem
        hours = int(arr[0])
        if hours < 1 or hours > 12:
            raise HoursException

        if arr[1].endswith("am"):
            if hours == 12:
                hours = 0
            arr[1] = arr[1].replace("am", "")
        elif arr[1].endswith("pm"):
            if hours != 12:
                hours += 12
            arr[1] = arr[1].replace("pm", "")
        else: # neither am nor pm detected
            raise MeridiemException

        # validate the minutes
        minutes = int(arr[1])
        if minutes < 0 or minutes > 59:
            raise MinutesException

    else:   # 9am
        if len(time) < 3 or len(time) > 4:
            raise TimeException

        minutes = 0
        if time.endswith("am"):
            time = time.replace("am", "")
            if time.startswith("12"):
                hours = int(time) - 12
            else:
                hours = int(time)
        elif time.endswith("pm"):
            time = time.replace("pm", "")
            if time.startswith("12"):
                hours = int(time)
            else:
                hours = int(time) + 12
        else: # neither am nor pm detected
            raise MeridiemException

        if hours < 1 or hours > 24:
            raise HoursException

    return timedelta(hours=hours, minutes=minutes)

# calculate hours for the day and return them
def calculate_time(start: str, end: str, less: str, more: str) -> float:
    # build clock-in and clock-out time
    startTime = build_time_delta(start)
    endTime = build_time_delta(end)

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
