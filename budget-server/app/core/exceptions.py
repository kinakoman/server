from fastapi.exceptions import RequestValidationError
from fastapi import Request
from fastapi.responses import JSONResponse

async def validation_exception_handler(request:Request,exception:RequestValidationError):
    return JSONResponse(
        status_code=422,
        content={
            "message": "リクエストのバリデーションに失敗",
            "errors": exception.errors(),
            "body": exception.body,
        },
    )