# APIのエンドポイントを設定
from fastapi import APIRouter
from app.schemas import UserCreate, UserResponse
from app.services import user_service

router = APIRouter()
@router.post("/user/", response_model=UserResponse)
def create_user(user: UserCreate):
    """
    Create a new user.
    """
    return user_service.create_user(user)