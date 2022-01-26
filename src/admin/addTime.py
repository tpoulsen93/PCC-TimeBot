import databaseAccess as da


# get all the inputs
first =     input("Enter employee first name:       ")
last =      input("Enter employee last name:        ")
date =      input("Enter date:     <YYYY-MM-DD>     ")
time =      input("Enter hours:                     ")

# double check everything
print("\n")
print(f"name:   {first.title()} {last.title()}")
print(f"date:   {date}")
print(f"hours:  {time}")

good = input("\nSubmit? (y/n)   ")
print()
if good.startswith("y"):
    print(da.add_time(first, last, date, time))
else:
    print("Cancelled. See you in the next life...")
