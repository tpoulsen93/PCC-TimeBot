import sys, os
import exceptions
import messageParser

from fastapi import FastAPI, Form, Response, Request, HTTPException
from twilio.twiml.messaging_response import MessagingResponse
from twilio.request_validator import RequestValidator  


app = FastAPI()


def text_usage(response: MessagingResponse):
    response.message(
        "Usage: <time/draw> <first name> <last name> <start time> <last time> <lunch> [<extra>]\n\
            Here is an example...")
    response.message("Time Taylor Poulsen 11:46am 5:04pm 1.25 3.6")

@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.get("/test/{string}")
def read_sms(string):
    return {"this is the string": string}


@app.post("/sms")
async def parse_message(request: Request, From: str = Form(...), Body: str = Form(...)):
    # make sure the request is from Twillio not a rando
    validator = RequestValidator(os.environ["TWILIO_AUTH_TOKEN"])
    form_ = await request.form()
    if not validator.validate(str(request.url), form_, request.headers.get("X-Twilio-Signature", "")):
        raise HTTPException(status_code=400, detail="Error in Twilio Signature")

    # process the message
    response = MessagingResponse()     
    try:
        msg = messageParser.process_message(Body)
    except Exception as e:
        print(e.with_traceback)
        return Response(content=str(response), media_type="application/xml")

    if not msg:
        print("Ignored message:")
        print(f"[{Body}]")
    else:
        response.message(msg)
        if msg.startswith("Error"):
            text_usage(response)
        print("Processed message:")
        print(f"[{Body}]")
        print("Response:")
        print(f"[{msg}]")
        

    sys.stdout.flush()
    return Response(content=str(response), media_type="application/xml")