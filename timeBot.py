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
number = os.environ['test_phone']
account = os.environ['account_sid']
token = os.environ['auth_token']

client = Client(account, token)

app = FastAPI()

users = {
    "+12083500006": "Taylor Poulsen",
    "+14322761331": "Daniel Bogden"
}


@app.post("/sms")
async def response(From: str = Form(...), Body: str = Form(...)) -> str:
    # print(f"Mesage from: {From}")
    # print(f"Body: {Body}")

    try:
        # body = message
        result = messageParser.process_message(Body)
        print(result)   
    except exceptions.DrawException:
        # body = "Draw formatted incorrectly" 
        print("Draw formatted incorrectly")
    except Exception as e:
        print("there was an exception")
        print(e.with_traceback)



    if not response:
        mess = client.messages.create(
            body = "Join Earth's mightiest heroes. Like Kevin Bacon.",
            from_ = '+18324971734',
            to = From
        )
    else:
        mess = client.messages.create(
            body = response,
            from_ = '+18324971734',
            to = From
        )

    # response = MessagingResponse()
    # response.message(mess)
    # return response
    



