from datetime import date, datetime, timedelta
import databaseAccess as da


class Day:
    def __init__(self, day: date):
        self.date = day
        self.weekday = day.strftime('%A')[:3]
        self.hours = 0

    def to_string(self):
        return "{:^10s} | {:^4s}| {:^5s}\n".format(str(self.date), self.weekday, str(self.hours))




class TimeCard:
    def __init__(self, id, start_date, end_date):
        self.id = id
        self.days = []
        self.total_hours = 0

        start = datetime.strptime(start_date, '%Y-%m-%d').date()
        end = datetime.strptime(end_date, '%Y-%m-%d').date()
        self.payday = end + timedelta(days=12)

        current = start
        while current <= end:
            self.days.append(Day(current))
            current += timedelta(days=1)

        employee = da.get_employee(id)
        for e in employee:
            self.name = f"{e.first_name.title()} {e.last_name.title()}"
            self.email = e.email
            self.phone = e.phone
            self.wage = float(e.wage)


    def add_hours(self, date: str, hours: float):
        for d in self.days:
            if str(d.date) == date:
                d.hours = hours
                self.total_hours += hours


    def to_string(self) -> str:
        result =  f"{self.name}\n\n\n"
        result +=  "{:^11s}|{:^5s}|{:>6s}\n".format("Date", "Day", "Hours")
        result += "{:-^11}+{:-^5}+{:->6}\n".format("","","")

        for d in self.days:
            result += d.to_string()

        result += f"\n\nWage:  ${self.wage}/hr\n"
        result += f"Total hours:  {self.total_hours}\n"
        result +=  "Gross pay:  ${:0,.2f}\n".format(self.wage * self.total_hours)
        result += f"Payday:  {self.payday}"

        return result

    





