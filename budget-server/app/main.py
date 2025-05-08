from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from app.api.routes import router
from dotenv import load_dotenv
from app.core.exceptions import validation_exception_handler

load_dotenv("../../.env")

app= FastAPI()

app.include_router(router)
app.add_exception_handler(RequestValidationError, validation_exception_handler)

@app.get("/")
def read_root():
    return {"Hello": "World"}