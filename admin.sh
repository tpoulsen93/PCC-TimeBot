#!/bin/bash

# list of all possible actions whose index is a primary key for the scripts
actions=(
    "Add a new employee"
    "Update an existing employee"
    "Make a new time submission"
)

# list of all the scripts that correspond to the actions above
scripts=(
    "admin.addEmployee"
    "admin.updateEmployee"
    "admin.addTime"
)

# print out all the possible script options
echo "Admin actions:"
echo

for i in ${!actions[@]}
do
    echo "[$i]  ${actions[$i]}"
done

echo
read -p "Enter the number of the action to execute -->  " i

# change directory attempt to run the selected script
if [[ $i -ge 0 && $i -lt ${#actions[@]} ]]; then
    echo "Initializing..."; echo
    cd src; python3 -m ${scripts[$i]}
else
    echo "Bad choice... Game over."
fi

