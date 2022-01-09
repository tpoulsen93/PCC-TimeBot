import sys, os
import messageParser

from fastapi import FastAPI, Form, Response, Request, HTTPException
from twilio.twiml.messaging_response import MessagingResponse
from twilio.request_validator import RequestValidator  

app = FastAPI()


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
    msg = response.message(f"Hi {From}, you said: {Body}")
    # msg = messageParser.process_message(Body)
    # if not message:
    #     print("We received a message but it wasn't for us")
    #     print(Body)
    # msg = response.message()
    print(msg)
    sys.stdout.flush()
    return Response(content=str(response), media_type="application/xml")
