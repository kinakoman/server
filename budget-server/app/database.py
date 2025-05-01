# データベースの接続設定を行う

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
import os
from dotenv import load_dotenv
load_dotenv()

DIALECT = "mysql+mysqlconnector"
USERNAME = os.environ.get("DB_USER")  # 環境変数から取得
PASSWORD = os.environ.get("DB_PASSWORD")  # 環境変数から取得
HOSTNAME = os.environ.get("DB_HOST")  # 環境変数から取得
PORT = os.environ.get("DB_PORT")  # 環境変数から取得
DATABASE_NAME = os.environ.get("DB_NAME")  # 環境変数から取得、デフォルトはtest_db



DATABASE_URL = f"{DIALECT}://{USERNAME}:{PASSWORD}@{HOSTNAME}:{PORT}/{DATABASE_NAME}"

engine = create_engine(DATABASE_URL)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

