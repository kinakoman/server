from fastapi import FastAPI
from app.api.routes import router
from dotenv import load_dotenv

load_dotenv("../../.env")

app= FastAPI()

app.include_router(router)


@app.get("/")
def read_root():
    return {"Hello": "World"}