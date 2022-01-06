# Download the helper library from https://www.twilio.com/docs/python/install
import json
import os
from twilio.rest import Client
from dotenv import load_dotenv
from datetime import date

from fastapi import FastAPI, Request, Form
from pydantic import BaseModel

from twilio.twiml.messaging_response import MessagingResponse

import databaseAccess
import exceptions
import messageParser

load_dotenv()
account = os.environ.get('account_sid')
token = os.environ.get('auth_token')
client = Client(account, token)

app = FastAPI()

users = {
    "+12083500006": "Taylor Poulsen",
    "+14322761331": "Daniel Bogden"
}



# message = client.messages \
#                 .create(
#                      body="Join Earth's mightiest heroes. Like Kevin Bacon.",
#                      from_='+18324971734',
#                      to='+12083500006'
#                  )

# print(message.sid)


@app.post("/sms")
async def response(From: str = Form(...), Body: str = Form(...)) -> str:
    # print(f"Mesage from: {From}")
    # print(f"Body: {Body}")

    try:
        # body = messageParser.process_message(Body)
        print(messageParser.process_message(Body))   
    except exceptions.DrawException:
        # body = "Draw formatted incorrectly" 
        print("Draw formatted incorrectly")
    except :
        print("there was an exception")



    mess = client.messages \
                .create(
                     body = "Join Earth's mightiest heroes. Like Kevin Bacon.",
                     from_ = '+18324971734',
                     to = From
                 )

    response = MessagingResponse()
    response.message(mess)
    return response
    



