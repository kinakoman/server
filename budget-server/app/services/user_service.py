# エンドポイントに設定するハンドラを定義する
from app.schemas import CreateExpenseRequest,CreateExpenseResponse,GetExpenseRequest,GetExpenseResponse
from app.crud.user_crud import create_expense,get_expense_by_year
from fastapi import HTTPException

# def get_budget(date:BudgetGet) -> BudgetResponse:
#     budget = get_budget_by_month(date)
#     if not budget:
#         raise HTTPException(status_code=404, detail="Budget not found")
    
#     # return BudgetResponse(month=budget.month, amount=budget.amount)
#     return budget

# def record(data:ExpenseRequest)->ExpenseOut:
#     new_record=record_budget(data)
    
#     if not new_record:
#          raise HTTPException(status_code=404, detail="New Record not found")
    
#     print(new_record.id)
#     # return ExpenseOut(id=new_record.id)
#     return new_record

def get_expense_by_year_service(request:GetExpenseRequest) -> GetExpenseResponse:
     print("now")
     expenses = get_expense_by_year(request)
     
     if not expenses:
          raise HTTPException(status_code=404, detail="Expenses not found")
     
     return expenses

def create_expense_service(requset:CreateExpenseRequest)->CreateExpenseResponse:
    new_record=create_expense(requset)
    
    if not new_record:
         raise HTTPException(status_code=404, detail="New Record not found")
    
    print(new_record.id)
    # return ExpenseOut(id=new_record.id)
    return new_record