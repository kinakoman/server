# データベースのテーブルモデルを定義

from sqlalchemy import Column,Integer,String,Boolean
from sqlalchemy.ext.declarative import declarative_base


Base=declarative_base()

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
    