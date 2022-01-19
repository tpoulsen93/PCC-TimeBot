import databaseAccess as da


# get all the inputs
first =     input("Enter employee first name:       ")
last =      input("Enter employee last name:        ")
wage =      input("Enter employee wage:             ")
email =     input("Enter employee email address:    ")
phone =     input("Enter employee phone number:     ")
s_first =   input("Enter supervisor first name:     ")
s_last =    input("Enter supervisor last name:      ")


# double check everything
print("\n")
print(f"name:       {first.title()} {last.title()}")
print(f"wage:       {wage}")
print(f"email:      {email}")
print(f"phone:      {phone}")
print(f"supervisor: {s_first.title()} {s_last.title()}")

good = input("Submit? (y/n)   ")
if good.startswith("y"):
    print(da.add_employee(first, last, wage, email, phone, s_first, s_last))
else:
    print("Cancelled. See you in the next life...")
