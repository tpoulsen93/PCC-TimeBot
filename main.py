from typing import Optional

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()


# class Item(BaseModel):
#     name: str
#     price: float
#     is_offer: Optional[bool] = None


@app.get("/")
def read_root():
    return {"Hello": "World"}


# @app.post("/sms")
# async def parse_message(From: str = Form(...), Body: str = Form(...)) -> str:

#     return {"item_name": item.name, "item_id": item_id}
