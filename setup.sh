#!/bin/bash

#Author: Taylor Poulsen
#Date: January 2021

echo "compile java code..."
echo
javac *.java
echo

echo "install xlsx2csv to convert excel spreadsheet to csv...(Ubuntu)"
echo
sudo apt-get install xlsx2csv
echo

echo "install ssmtp to send texts and emails"
echo
sudo apt-get install ssmtp
echo

echo
echo "Look up how to configure ssmtp and do it."
echo "Create the employee database using this format:"
echo "firstname [sheetId]"
echo "Then you should be ready to rumble.... I hope..."
