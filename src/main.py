import os
import sys

from fastapi import FastAPI, Form, HTTPException, Request, Response
from twilio.request_validator import RequestValidator
from twilio.twiml.messaging_response import MessagingResponse

import src.databaseAccess as databaseAccess
import src.messageParser as messageParser
import src.twilioActions as twilioActions

app = FastAPI()


def text_usage(response: MessagingResponse):
    response.message(
        "Usage: Time <first name> <last name> <start time> <end time> <subtracted hours(lunch)>\
            [<additional hours(drive time)>]\nExample:")
    response.message("Time Taylor Poulsen 11:46am 5:04pm 1.25 3.6")


@app.get("/")
def read_root():
    return {"PCC": "Poulsen Concrete Contractors Inc."}


@app.get("/submitHours")
async def submit_hours(id: int, date: str, hours: float):
    result = databaseAccess.add_time(id, date, hours)
    twilioActions.confirm_submission(id, result, "TimeApp")
    return result


@app.post("/sms")
async def parse_message(request: Request, From: str = Form(...), Body: str = Form(...)):
    # make sure the request is from Twilio not a rando
    validator = RequestValidator(os.environ["TWILIO_AUTH_TOKEN"])
    form_ = await request.form()
    if not validator.validate(str(request.url), form_, request.headers.get("X-Twilio-Signature", "")):
        print("unexpected user encountered")
        raise HTTPException(
            status_code=400, detail="Error in Twilio Signature")

    # process the message
    response = MessagingResponse()
    try:
        msg = messageParser.process_message(Body, From)
    except Exception as e:
        response.message(
            "Encountered an unexpected error. Check your format and try again.")
        print("Encountered unexpected error in message:")
        print(f"[{Body}]")
        print(e)
        print(e.with_traceback)
        return Response(content=str(response), media_type="application/xml")

    if not msg:
        print(f"Ignored message from {From}:")
        print(f"[{Body}]")
    else:
        if msg.startswith("Help"):
            text_usage(response)
        elif msg.startswith("Error"):
            response.message(msg)
            text_usage(response)
        else:
            response.message(msg)

        print(f"Processed message from {From}:")
        print(f"[{Body}]")
        print("Responded:")
        print(f"[{msg}]")

    sys.stdout.flush()
    return Response(content=str(response), media_type="application/xml")
