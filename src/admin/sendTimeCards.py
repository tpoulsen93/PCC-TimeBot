from email import encoders
from email.header import Header
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
import databaseAccess as da
import timeCard as tc
import smtplib, sys, os

# def print_usage():
#     print("Usage: generateTimeCards.py <period start date> <period end date>")
#     print("Date format: YYYY-MM-DD")
#     sys.exit()

# # check the commandline arguments
# if len(sys.argv) != 3:
#     print_usage()

# # get arguments from commandline
# start = sys.argv[1]
# end = sys.argv[2]

# ask for necessary inputs
print("Date format:  YYYY-MM-DD")
start = input("Enter pay period start date:  ")
end   = input("Enter pay period end date:    ")

payday = False
timecards = {}

# get all the rows from the database between the start and end dates
print(f"\nGetting all hours submitted between {start} and {end}...")
result = da.get_time_cards(start, end)

print("Building time cards...")
if result:
    for r in result:
        if r.id in timecards:
            timecards[r.id].add_hours(str(r.date), r.time)
        else:
            timecards[r.id] = tc.TimeCard(r.id, start, end)
            timecards[r.id].add_hours(str(r.date), r.time)
            # grab the payday off the first timecard
            if not payday:
                payday = timecards[r.id].payday


    # setup smtp to send emails
    print("Connecting to SMTP server...")
    with smtplib.SMTP_SSL("smtp.gmail.com", 465) as smtp:
        smtp.login(os.environ['SMTP_USERNAME'], os.environ['SMTP_PASSWORD'])

        # loop through the timecards and send them to their recipients
        tmpPath = './timeCard.txt'
        for t in timecards.values():
            msg = MIMEMultipart()
            msg['From'] = f"{Header('TimeBot').encode()} <{os.environ['SMTP_USERNAME']}>"
            print(f"Sending time card to {t.name}...")

            msg['Subject'] = f"Time Card for payday: {payday}"
            msg['To'] = t.email

            # write the current timecard to a file
            f = open(tmpPath, "w")
            f.write(t.to_string())
            f.close()

            # add the timecard file as an attachment to the email
            f = open(tmpPath, "rb")
            card = MIMEBase('application', 'octet-stream')
            card.set_payload(f.read())
            encoders.encode_base64(card)
            card.add_header('Content-Disposition', 'attachment; filename="TimeCard.txt"')
            msg.attach(card)

            try:
                smtp.send_message(msg)
            except Exception as e:
                print(f"Failed to send timecard to {t.name}")
                print(f"\t{e.__cause__}")
                print(f"\t{e.with_traceback}")

            # clean up the old message before next iteration
            del msg
            f.close()

        # send all the results to TP for payroll submission
        print("Sending payroll totals to admin...")
        fro = os.environ['SMTP_USERNAME']
        to = da.get_employee_email(da.get_employee_id('taylor', 'poulsen'))
        
        # build the message body from all the timecards
        hours_sum = 0
        cost_sum = 0
        body =  f"From: TimeBot <{fro}>\n"
        body += f"To: TP <{to}>\n"
        body += f"Subject: PCC Payroll totals for payday -> {payday}\n\n"
        body += f"Pay period: {start}  <->  {end}\nPayday: {payday}\n\n"

        for t in timecards.values():
            hours_sum += t.total_hours
            cost_sum += t.total_hours * t.wage
            body += f"{t.name}  -->  {round(t.total_hours, 2)}\n"
            
        body += f"\nTotal Hours  -->  {hours_sum}"
        body += f"\nTotal Cost  -->  {cost_sum}"

        smtp.sendmail(fro, to, body)

    print("Mission accomplished")

else:
    print("No hours found for indicated dates...")

#print everything to the logs
sys.stdout.flush()