# httpリクエストのデータ型を定義
from pydantic import BaseModel
from typing import Optional

# /budget/create-expenseのリクエストデータ型
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

# /budget/create-expenseのレスポンスデータ型
class CreateExpenseResponse(BaseModel):
    id :int
    
# /budget/get-expense-by-yearのリクエストデータ型
class GetExpenseRequest(BaseModel):
    year: int

# /budget/get-expense-by-yearのレスポンスデータ型
class GetExpenseResponse(BaseModel):
    id: int
    user:str
    expense: int
    year: int
    month: int
    day: int
    description: str
    calculation: bool
    settlement: bool
    fixed: bool
    
    
