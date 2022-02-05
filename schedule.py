import smtplib, os

with smtplib.SMTP_SSL("smtp.gmail.com", 465) as smtp:
    smtp.login(os.environ['SMTP_USERNAME'], os.environ['SMTP_PASSWORD'])
    fro = os.environ['SMTP_USERNAME']
    to = "t-poulsen@hotmail.com"

    body =  f"From: TimeBot <poulsent.23@gmail.com>"
        
    body += f"\n\n hiiiiiii"

    smtp.sendmail(fro, to, body)
