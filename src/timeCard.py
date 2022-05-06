from datetime import date, datetime, timedelta
import databaseAccess as da



class TimeCard:
    def __init__(self, id, start_date, end_date):
        self.id = id
        self.days = {}
        self.total_hours = 0

        start = datetime.strptime(start_date, '%Y-%m-%d').date()
        end = datetime.strptime(end_date, '%Y-%m-%d').date()
        self.payday = end + timedelta(days=12)

        current = start
        while current <= end:
            self.days[str(current)] = 0
            current += timedelta(days=1)

        employee = da.get_employee(id)
        for e in employee:
            self.name = f"{e.first_name.title()} {e.last_name.title()}"
            self.email = e.email
            self.phone = e.phone
            self.wage = float(e.wage)


    def add_hours(self, date: str, hours: float):
        self.days[date] = round(hours, 2)
        self.total_hours += hours


    def build_day_line(self, date: str) -> str:
        return "{:^10s} | {:^4s}| {:^5s}\n".format(
            date,
            datetime.strptime(date, '%Y-%m-%d').strftime('%A')[:3],
            f"{round(self.days[date], 2):g}"
        )


    def to_string(self) -> str:
        result =  f"{self.name}\n\n\n"
        result +=  "{:^11s}|{:^5s}|{:>6s}\n".format("Date", "Day", "Hours")
        result += "{:-^11}+{:-^5}+{:->6}\n".format("","","")

        for d in self.days.keys():
            result += self.build_day_line(d)

        # result += f"\n\nWage:  ${round(self.wage, 2)}/hr\n"
        result += f"Total hours:  {round(self.total_hours, 2)}\n"
        # result +=  "Gross pay:  ${:0,.2f}\n".format(self.wage * self.total_hours)
        result += f"Payday:  {self.payday}\n"

        return result

    





