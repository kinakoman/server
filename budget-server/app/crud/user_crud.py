# データベースの操作関数を定義するファイル
from app.schemas import CreateExpenseRequest
from app.database import SessionLocal
from app.models import Expenses


# def get_user_by_name(name: str):
#     db= SessionLocal()
#     try:
#         user = db.query(User).filter(User.name == name).first()
#         return user
#     finally:
#         db.close()

# def create_user( user: UserCreate):
#     db= SessionLocal()
#     try:
#         new_user=User(name=user.name, age=user.age)
#         db.add(new_user)
#         db.commit()
#         db.refresh(new_user)
#         return new_user
#     finally:
#         db.close()

# def get_budget_by_month(Date:BudgetGet):
#     db = SessionLocal()
#     try:
#         budgets = db.query(Expense).filter(Expense.expense_date==date(year=Date.year,month=Date.month,day=Date.day)).all()
#         return budgets
#     finally:
#         db.close()

# def record_budget(data:ExpenseRequest):
#     Date=date(year=2025,month=12,day=25)
    
#     db=SessionLocal()
#     try:
#         new_expense=Expense(user=data.user,expense_date=Date,amount=data.amount)
#         db.add(new_expense)
#         db.commit()
#         db.refresh(new_expense)
#         return new_expense
#     finally:
#         db.close()

                
def create_expense(data:CreateExpenseRequest):
    db=SessionLocal()
    try:
        new_expense=Expenses(
            user=data.user,
            expense=data.expense,
            year=data.year,
            month=data.month,
            day=data.day,
            description=data.description,
            calculation=data.calculation,
            settlement=data.settlement,
            fixed=data.fixed
        )
        db.add(new_expense)
        db.commit()
        db.refresh(new_expense)
        return new_expense
    finally:
        db.close()