# Author: Taylor Poulsen


# parent exception class for errors caught in messageParser
class MessageParseException(Exception):
    pass


# error found in time formatting
class TimeException(MessageParseException):
    pass


class TimeFormatException(TimeException):
    pass


class IllegalTimeException(TimeException):
    pass


class HoursException(TimeException):
    pass


class MinutesException(TimeException):
    pass


class MeridiemException(TimeException):
    pass


# error found in dollar amount
class DollarException(MessageParseException):
    pass


class DrawException(DollarException):
    pass


# error found in float format
class FloatException(MessageParseException):
    pass


class LunchException(FloatException):
    pass


class ExtraException(FloatException):
    pass


# user not found in database
class NoSuchUserException(MessageParseException):
    pass
