import databaseAccess as da


# get all the inputs
first = input("Enter employee first name:       ")
last = input("Enter employee last name:        ")

# indeces for result
PHONE = 4
EMAIL = 5
SUPERVISOR_ID = 6

id = da.get_employee_id(first, last)
employee = da.get_employee(id)

print(f"Name:       {first.title()} {last.title()}")
print(f"ID:         {id}")
print(f"Phone:      {employee[PHONE]}")
print(f"Email:      {employee[EMAIL]}")
print(f"Supervisor: {da.get_employee_name(employee[SUPERVISOR_ID])}")
