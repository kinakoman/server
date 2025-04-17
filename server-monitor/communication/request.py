# httpリクエストを提供するモジュール
import requests
import os 
from requests.auth import HTTPBasicAuth
from dotenv import load_dotenv

load_dotenv()
# 環境変数の読み込み
server_scheme=os.environ['SERVER_SCHEME']
server_host=os.environ['MAIN_SERVER_HOST']
server_port=os.environ['MAIN_SERVER_PORT']
auth_user=os.environ['LOCAL_AUTH_USER']
auth_pass=os.environ['LOCAL_AUTH_PASSWORD']

# サーバー状態を取得する関数
# 戻り値
# ture:サーバー稼働
# false:サーバー未稼働
def ServerStatus()->bool:
    # /status/にアクセス
    url=f"{server_scheme}://{server_host}:{server_port}/status/"
    try:
        # Basic認証でサーバーにアクセス
        res=requests.post(url,auth=HTTPBasicAuth(auth_user,auth_pass))
        res.raise_for_status()
        # サーバーと接続出来たらTrueを返す
        return True
    except requests.exceptions.ConnectionError as e:
        # サーバーと接続に失敗したらFalseを返す
        return False 
    except requests.exceptions.HTTPError as e:
        return False
    