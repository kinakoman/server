# httpリクエストのデータ型を定義
from pydantic import BaseModel
from typing import Optional


    
class CreateExpenseRequest(BaseModel):
    user: str
    expense:int
    year:int
    month:int
    day:int
    description:str
    calculation:Optional[bool] = True
    settlement:Optional[bool] = False
    fixed:Optional[bool] = False

class CreateExpenseResponse(BaseModel):
    id :int