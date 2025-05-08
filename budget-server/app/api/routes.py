# APIのエンドポイントを設定
from fastapi import APIRouter,HTTPException
from app.schemas import CreateExpenseRequest,CreateExpenseResponse
from app.services.user_service import create_expense_service
# from fastapi import Query

router = APIRouter()

# 家計簿の取得関数
# @router.get("/budget/get",response_model=BudgetResponse)
# def get_budget(year: int=Query(...,description="year"), month: int=Query(...,description="month"), day: int=Query(...,description="day")):
#     """
#     Get budget for a specific year and month.
#     """
#     date_model= BudgetGet(year=year, month=month, day=day)
    
#     return user_service.get_budget(date_model)


# @router.post("/budget/record",response_model=ExpenseOut)
# def record(data:ExpenseRequest):
#     # print(data)
#     return user_service.record(data)

@router.post("/budget/create-expense",response_model=CreateExpenseResponse)
def create_expense_method(request:CreateExpenseRequest):
    try:
        return create_expense_service(request)
    except HTTPException as e:
        raise e
    except Exception as e:
        raise HTTPException(status_code=500, detail="Internal Server Error")