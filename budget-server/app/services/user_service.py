# エンドポイントに設定するハンドラを定義する
from app.schemas import UserCreate,UserResponse
from app.crud import user_crud
from fastapi import HTTPException

def create_user(user: UserCreate)->UserResponse:
    db_user=user_crud.get_user_by_name(user.name)
    if db_user:
        raise HTTPException(status_code=400, detail="User already exists")
    
    new_user=user_crud.create_user(user)
    return UserResponse(id=new_user.id, name=new_user.name, age=new_user.age)