from email import encoders
from email.header import Header
from email.mime.base import MIMEBase
from email.mime.multipart import MIMEMultipart
from sqlalchemy import false
import databaseAccess as da
import timeCard as tc
import smtplib, sys, os



def print_usage():
    print("Usage: generateTimeCards.py <period start date> <period end date>")
    print("Date format: YYYY-MM-DD")
    sys.exit()

# check the commandline arguments
if len(sys.argv) != 3:
    print_usage()

start = sys.argv[1]
end = sys.argv[2]
payday = false
timecards = {}

# get all the rows from the database between the start and end dates
print(f"Getting all hours from {start} to {end}")
result = da.get_time_cards(start, end)
if result:
    for r in result:
        if r.id in timecards:
            timecards[r.id].add_hours(str(r.date), r.time)
        else:
            timecards[r.id] = tc.TimeCard(r.id, start, end)
            timecards[r.id].add_hours(str(r.date), r.time)

    # setup smtp to send emails
    print("Connecting to SMTP server...")
    msg = MIMEMultipart()
    msg['From'] = f"{Header('TimeBot').encode()} <{os.environ['SMTP_USERNAME']}>"
    with smtplib.SMTP_SSL("smtp.gmail.com", 465) as smtp:
        smtp.login(os.environ['SMTP_USERNAME'], os.environ['SMTP_PASSWORD'])

        # loop through the timecards and send them to their recipients
        tmpPath = './timeCard.txt'
        for t in timecards.values():
            print(f"Sending time card to {t.name}...")

            # grab the payday off the first timecard
            if payday == false:
                payday = t.payday

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
                print(e.with_traceback)

            # clean up old headers before next iteration
            del msg['Subject']
            del msg['To']
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
        body += f"Subject: PCC Payroll totals for the week of {start}\n\n"
        body += f"Pay period: {start}  <->  {end}\nPayday: {payday}\n\n"

        for t in timecards.values():
            hours_sum += t.total_hours
            cost_sum += t.total_hours * t.wage
            body += f"{t.name}  -->  {t.total_hours}\n"
            
        body += f"\nTotal Hours  -->  {hours_sum}"
        body += f"\nTotal Cost  -->  {cost_sum}"

        smtp.sendmail(fro, to, body)
