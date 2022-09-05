#!/bin/bash

# list of all possible actions whose index is a primary key for the scripts
actions=(
    "Send time cards"
    "Send updated time card"
    "Add a new employee"
    "Update an existing employee"
    "Make a new time submission"
    "Get employee summary"
)

# list of all the scripts that correspond to the actions above
scripts=(
    "admin.sendTimeCards"
    "admin.resendTimeCard"
    "admin.addEmployee"
    "admin.updateEmployee"
    "admin.addTime"
    "admin.getEmployee"
)

# print out all the possible script options
echo "Admin actions:"
echo

for i in ${!actions[@]}
do
    echo "[$i]  ${actions[$i]}"
done

echo
read -n 1 -p "Enter the number of the action to execute -->  " i

# change directory attempt to run the selected script
if [[ $i -ge 0 && $i -lt ${#actions[@]} ]]; then
    echo
    echo "Initializing..."; echo
    cd src; python3 -m ${scripts[$i]}
else
    echo "Bad choice... Game over."
fi
