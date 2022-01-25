import databaseAccess as da


# get all the inputs
first =     input("Enter employee first name:       ")
last =      input("Enter employee last name:        ")
date =      input("Enter date:     <YYYY-MM-DD>     ")
time =      input("Enter time:                      ")


# validate the date
d = date.split('-')
if len(d) != 3:
    print("Bad date.")
if int(d[0]) < 2022 or int(d[0] > 9999):
    print("Bad date. <year>")
if int(d[1]) < 1 or int(d[0] > 12):
    print("Bad date. <month>")
if int(d[2]) < 1 or int(d[2] > 31):
    print("Bad date. <day>")

# double check everything
print("\n")
print(f"name:   {first.title()} {last.title()}")
print(f"date:   {date}")
print(f"time:   {time}")

good = input("Submit? (y/n)   ")
if good.startswith("y"):
    print(da.add_time(first, last, date, time))
else:
    print("Cancelled. See you in the next life...")
