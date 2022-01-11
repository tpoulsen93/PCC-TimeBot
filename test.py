from datetime import date, datetime, timedelta
import datetime, pytz

d = (datetime.datetime.today() + timedelta(hours=5)).date()

print(datetime.datetime.today())
print(d)
#print(datetime.datetime.today().date())


#print(datetime.datetime.today().astimezone(pytz.timezone('est')).date())
#print(datetime.datetime.today().astimezone(pytz.timezone('mst')))



