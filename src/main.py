import src.messageParser as messageParser
import src.databaseAccess as databaseAccess
import json, sys, os

from fastapi import FastAPI, Form, Response, Request, HTTPException
from twilio.twiml.messaging_response import MessagingResponse
from twilio.request_validator import RequestValidator  


app = FastAPI()


def text_usage(response: MessagingResponse):
    response.message(
        "Usage: <time/draw> <first name> <last name> <start time> <end time> <subtracted hours(lunch)>\
            [<additional hours(drive time)>]\nHere is an example...")
    response.message("Time Taylor Poulsen 11:46am 5:04pm 1.25 3.6")

@app.get("/")
def read_root():
    return {"PCC": "Poulsen Concrete Contractors Inc."}


@app.get("/addHours/{first}/{last}/{date}/{time}")
def add_hours(first, last, date, time):
    try:
        result = databaseAccess.add_hours(first, last, date, time)
    except:
        return f"Failed to add hours... yikes" 
    return f"Submitted {time} hours on {date} for {first} {last}"


@app.get("/getTimeCards/{start}/{end}")
def get_time_cards(start, end):
    try:
        result = databaseAccess.get_time_cards(start, end)
    except:
        return f"Failed to get time cards... whoopsies :("
    return json.dumps(result)



@app.get("/updateEmployee/{first}/{last}/{target}/{value}")
def update_employee(first, last, target, value) -> str:
    try:
        databaseAccess.update_employee(first, last, target, value)
    except:
        return "Something bad happened... Failed to update employee"
    return f"Employee successfully updated:\n\
                Name:   {first} {last}\n\
                {target.title()}:   {value}"


@app.get("/addEmployee/{first}/{last}/{wage}/{email}/{phone}")
def add_employee(first, last, wage, email, phone) -> str:
    try:
        databaseAccess.add_employee(first, last, wage, email, phone)
    except:
        return "Something bad happened... Failed to add employee"
    return f"Employee successfully added:\n\
                Name:   {first} {last}\n\
                Wage:   ${wage}/hr\n\
                Email:  {email}\n\
                Phone:  {phone}"


@app.post("/sms")
async def parse_message(request: Request, From: str = Form(...), Body: str = Form(...)):
    # make sure the request is from Twillio not a rando
    validator = RequestValidator(os.environ["TWILIO_AUTH_TOKEN"])
    form_ = await request.form()
    if not validator.validate(str(request.url), form_, request.headers.get("X-Twilio-Signature", "")):
        print("unexpected user encountered")
        raise HTTPException(status_code=400, detail="Error in Twilio Signature")

    # process the message
    response = MessagingResponse()     
    try:
        msg = messageParser.process_message(Body)
    except Exception as e:
        response.message("Encountered an unexpected error. Check your format and try again.")
        print("Encountered unexpected error in message:")
        print(f"[{Body}]")
        print(e)
        print(e.with_traceback)
        return Response(content=str(response), media_type="application/xml")

    if not msg:
        print("Ignored message:")
        print(f"[{Body}]")
    else:
        if msg.startswith("Help"):
            text_usage(response)
        elif msg.startswith("Error"):
            response.message(msg)
            text_usage(response)
        else:
            response.message(msg)

        print("Processed message:")
        print(f"[{Body}]")
        print("Responded:")
        print(f"[{msg}]")
        

    sys.stdout.flush()
    return Response(content=str(response), media_type="application/xml")