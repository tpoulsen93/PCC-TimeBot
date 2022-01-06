#!/bin/bash

# Author: Taylor Poulsen
# Date: January 2021
# Purpose: Automate some of my weekly payroll tasks for Poulsen Concrete Contractors



# Print usage if incorrect # of arguments provided
if [ $# -lt 3 ] || [ $# -gt 3 ]; then
    echo Usage:   $0 '<start date> <payday> <# of days in pay period>'
    echo "Date format: mm-dd-yy"
    exit 0
fi



# add the actual address of the excel file here, will have to be updated once per year
spreadSheet=../Payroll_2021.xlsx

startDate=$1        #the start date of the payroll to be computed
payDay=$2           #the date of the payday
payPeriodLength=$3  #the length of the pay period

# make a directory for copies of all of the TimeCards
directory=../TimeCards/$payDay
if [ ! -e $directory ]; then
    mkdir $directory
fi

csv=$directory/tmp.csv  #temporary csv file that will be overwritten over and over during program execution

b64temp=base64_template.txt #temporary b64 template that will also be overwritten over and over during execution

# get timecard recipients from user input
read -p "Enter timecard recipients seperated by spaces: " -a employeeArray

# loop through the sheets creating temporary csv's, generating time cards, and sending them
for employee in ${employeeArray[@]}
do
    # get sheet id for each person and initialize their timecard
    sheetID=$(grep -oP "(?<=$employee \[).*(?=\])" employeeSheetIDs.txt)
    timeCard=$directory/TimeCard_$employee.txt

    # generate the template for the base64
    if [ ! -e $b64temp ]; then
        echo "Content-Type: application;" >> $b64temp
        echo "Content-Transfer-Encoding: base64" >> $b64temp
        echo "Content-Disposition: attachment; filename=\"TimeCard.txt\"" >> $b64temp
    fi

    # convert the spreadsheet to a csv
    xlsx2csv -s $sheetID $spreadSheet $csv

    # generate timecard from csv
    if [ -e $csv ]; then
        java GenerateTimeCards $csv $startDate $payDay $payPeriodLength > $timeCard
        rm -f $csv
    else
        echo "Error: No csv found for $employee."
        echo "Nothing was generated for $employee."
        continue
    fi

    # get info from the first few lines of the timecard
    phone=$(head -n 1 $timeCard)
    email=$(head -n 2 $timeCard | tail -1)
    name=$(head -n 3 $timeCard | tail -1)

    # remove the first 2 lines from the timecard
    sed -i '1,2d' $timeCard

    # generate the template for the base64 file
    echo "Subject: Time Card for $name" >> $b64temp
    echo "Content-Type: application;" >> $b64temp
    echo "Content-Transfer-Encoding: base64" >> $b64temp
    echo "Content-Disposition: attachment; filename=\"TimeCard.txt\"" >> $b64temp

    # convert the timecard to base64 and append to the template to be sent as an attachment in the email
    base64 $timeCard >> $b64temp

    # send the timecard and delete the base64 file
    ssmtp $email < $b64temp
    rm -f $b64temp
    echo "A time card was sent to $name"

    # send text message alerting recipient that their timecard has been sent
    #echo "Your time card for the pay period starting on $startDate has been sent to $email" | ssmtp $phone
    #echo "Direct Deposit will be scheduled for $payDay" | ssmtp $phone
    #echo "If you see any errors, or have any questions, please text me at 208-350-0006" | ssmtp $phone
    #echo "$name was alerted via text message"
done
