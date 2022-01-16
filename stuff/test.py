from twilio.rest import Client
import os


client = Client(
    os.environ['TWILIO_ACCOUNT_SID'],
    os.environ['TWILIO_AUTH_TOKEN']
)

twil = os.environ['TWILIO_PHONE']
tp = os.environ['TP_PHONE']

client.messages.create(from_=f"+1{twil}",
                       to=f"+1{tp}",
                       body='Ahoy from Twilio!')




