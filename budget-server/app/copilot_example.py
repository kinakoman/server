import os
from datetime import date
from fastapi import FastAPI, HTTPException, Depends, Request, status
from pydantic import BaseModel
from sqlalchemy import create_engine, Column, Integer, String, Float, Date, Boolean, and_
from sqlalchemy.orm import sessionmaker, declarative_base, Session

# Environment variables for authentication
API_KEY_NAME = os.getenv("PROXY_SECRET_HEADER", "X-Proxy-Secret")
API_KEY = os.getenv("PROXY_SECRET_KEY", "defaultsecret")

# Database connection (adjust the connection string as needed)
DATABASE_URL = os.getenv("DATABASE_URL", "mysql+mysqlconnector://user:password@localhost/budget_db")
engine = create_engine(DATABASE_URL, echo=True)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

# Database Model
class ExpenseDB(Base):
    __tablename__ = "expenses"
    id = Column(Integer, primary_key=True, index=True)
    user = Column(String(50), index=True)
    expense_date = Column(Date, index=True)  # '年月日'
    amount = Column(Float)
    settled = Column(Boolean, default=False)  # 清算完了

# Pydantic Schemas
class ExpenseCreate(BaseModel):
    user: str
    expense_date: date
    amount: float

class ExpenseOut(ExpenseCreate):
    id: int
    settled: bool

# Dependency to get DB session
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# Dependency to verify secret header
async def verify_api_key(request: Request):
    header_value = request.headers.get(API_KEY_NAME)
    if header_value != API_KEY:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")

# Create the FastAPI app
app = FastAPI()

# Create tables when the app starts
@app.on_event("startup")
def startup():
    Base.metadata.create_all(bind=engine)

# Endpoint to retrieve expenses for a given year and month
@app.get("/budget/get", response_model=list[ExpenseOut], dependencies=[Depends(verify_api_key)])
def get_expenses(year: int, month: int, db: Session = Depends(get_db)):
    try:
        start_date = date(year, month, 1)
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid year or month")
    # Determine the first day of the next month
    next_month = month + 1 if month < 12 else 1
    next_year = year if month < 12 else year + 1
    end_date = date(next_year, next_month, 1)
    
    expenses = db.query(ExpenseDB).filter(
        and_(ExpenseDB.expense_date >= start_date, ExpenseDB.expense_date < end_date)
    ).all()
    return expenses

# Endpoint to record a new expense
@app.post("/budget/record", response_model=ExpenseOut, dependencies=[Depends(verify_api_key)])
def record_expense(expense: ExpenseCreate, db: Session = Depends(get_db)):
    new_expense = ExpenseDB(
        user=expense.user,
        expense_date=expense.expense_date,
        amount=expense.amount,
        settled=False  # default value
    )
    db.add(new_expense)
    db.commit()
    db.refresh(new_expense)
    return new_expense

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)