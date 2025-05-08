# データベースのテーブルモデルを定義

from sqlalchemy import Column,Integer,String,Date,Boolean
from sqlalchemy.ext.declarative import declarative_base


Base=declarative_base()

# class User(Base):
#     __tablename__='python_test'
#     id=Column(Integer,primary_key=True,index=True, autoincrement=True)
#     name=Column(String,index=True)
#     age=Column(Integer)
    
# class Expense(Base):
#     __tablename__='expenses'
#     id=Column(Integer,primary_key=True,index=True, autoincrement=True)
#     user=Column(String(50),index=True)
#     expense_date=Column(Date,index=True)  # '年月日'
#     amount=Column(Integer)
#     settled=Column(Integer,default=False)  # 清算完了

class Expenses(Base):
    __tablename__='expenses'
    id=Column(Integer,primary_key=True,autoincrement=True)
    user=Column(String(50))
    expense=Column(Integer)
    year=Column(Integer)
    month=Column(Integer)
    day=Column(Integer)
    description=Column(String(100))
    calculation=Column(Boolean,default=True)
    settlement=Column(Boolean,default=False)
    fixed=Column(Boolean,default=False)
    