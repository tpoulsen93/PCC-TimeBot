from fastapi import FastAPI, Form, Response, Request, HTTPException
from twilio.twiml.messaging_response import MessagingResponse
from twilio.request_validator import RequestValidator  

app = FastAPI()


@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.get("/test/{string}")
# async def parse_message(From: str = Form(...), Body: str = Form(...)) -> str:
def read_sms(string):
    return {"this is the string": string}


@app.post("/sms")
async def parse_message(request: Request, From: str = Form(...), Body: str = Form(...)):
    # make sure the request is from Twillio not a rando
    validator = RequestValidator(os.environ["TWILIO_AUTH_TOKEN"])
    form_ = await request.form()
    if not validator.validate(str(request.url), form_, request.headers.get("X-Twilio-Signature", "")):
        raise HTTPException(status_code=400, detail="Error in Twilio Signature")

    response = MessagingResponse() 
    msg = response.message(f"Hi {From}, you said: {Body}")
    print(msg)
    return Response(content=str(response), media_type="application/xml")
