import databaseAccess as da
import timeCard as tc
import sys, os

def print_usage():
    print("Usage: generateTimeCards.py <period start date> <period end date>")
    print("Date format: YYYY-MM-DD")
    sys.exit()

# check the commandline arguments
if len(sys.argv) != 3:
    print_usage()

start = sys.argv[1]
end = sys.argv[2]

timecards = {}
result = da.get_time_cards(start, end)
for r in result:
    if r.id in timecards:
        tcard = timecards[r.id]
        # print(type(tcard))
        # print(tcard)
        # tc.add_hours(r.date, r.time)
        # print(f"id: {r.id}   date: {r.date}    time: {r.time}")
        # print(timecards[r.id].to_string())
    else:
        timecards[r.id] = tc.TimeCard(r.id, start, end)
        tcard = timecards[r.id]
        tcard.days[0].hours = 8
        # print(tcard.to_string())

for t in timecards.values():
    print(f"{t.to_string()}\n\n----------------\n\n")


