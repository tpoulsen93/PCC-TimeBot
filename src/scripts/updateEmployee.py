import databaseAccess as da


# get all the inputs
first =     input("Enter employee first name:       ")
last =      input("Enter employee last name:        ")
target =    input("Enter target: <wage|email|phone> ")
value =     input("Enter new value:                 ")

# double check everything
print("\n")
print(f"name:   {first.title()} {last.title()}")
print(f"target: {target}")
print(f"value:  {value}")

good = input("Submit? (y/n)   ")
if good.startswith("y"):
    print(da.update_employee(first, last, target, value))
else:
    print("Cancelled. See you in the next life...")