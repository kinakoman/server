# データベースの操作関数を定義するファイル
from app.models import User
from app.schemas import UserCreate
from app.database import SessionLocal

def get_user_by_name(name: str):
    db= SessionLocal()
    try:
        user = db.query(User).filter(User.name == name).first()
        return user
    finally:
        db.close()

def create_user( user: UserCreate):
    db= SessionLocal()
    try:
        new_user=User(name=user.name, age=user.age)
        db.add(new_user)
        db.commit()
        db.refresh(new_user)
        return new_user
    finally:
        db.close()
        
