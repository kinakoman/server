# データベースのテーブルモデルを定義

from sqlalchemy import Column,Integer,String
from sqlalchemy.ext.declarative import declarative_base

Base=declarative_base()

class User(Base):
    __tablename__='python_test'
    id=Column(Integer,primary_key=True,index=True, autoincrement=True)
    name=Column(String,index=True)
    age=Column(Integer)