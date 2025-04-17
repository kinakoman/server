# データベースとのアクセスを提供するモジュール
import mysql.connector
from mysql.connector.errors import Error
import os
import pandas as pd
from typing import Tuple
from dotenv import load_dotenv

load_dotenv()
# 環境変数の読み込み
user_name=os.environ['DB_USER']
user_password=os.environ['DB_PASSWORD']
db_host=os.environ['DB_HOST']
db_name=os.environ['DB_NAME']
db_table=os.environ['DB_TABLE']

# データベースサーバーとのコネクションを作成
# 戻り値
# connection:データベースとのコネクション
# err:コネクションエラー、
def createServerConnection()->Tuple[mysql.connector.connection_cext.CMySQLConnection,Error]:
    try:
        # サーバーとのコネクション
        connection=mysql.connector.connect(
            host=db_host,
            user=user_name,
            password=user_password,
            database=db_name)
        return connection,None
    # コネクションエラー処理
    except Error as err:
        return None,err
    
# server_logのデータを取得
# 引数
# connection:mysqlとのコネクション
# err:mysqlとのコネクションエラー
# rowNum:GUIに表示するデータテーブルの行数
# GUIcolumns:GUIの設定カラムラベル
# 戻り値
# table:データベースから取得したサーバーログ
def getServerLog(connection:mysql.connector.connection_cext.CMySQLConnection,rowNum:int,GUIcolumns:list,err:Error)->pd.DataFrame:
    if err!=None:
        return pd.DataFrame([[""]*len(GUIcolumns)]*rowNum,columns=GUIcolumns)
    
    try:  
        # カーソルを実体化
        cursor=connection.cursor()
        # server_logのカラム名を取得
        cursor.execute(f"DESC {db_table}")
        # カラム名をリストに変換
        table_column_label=[row[0] for row in cursor.fetchall()]
        # server_logのデータを取得
        cursor.execute(f"SELECT * FROM {db_table} ORDER BY timestamp DESC LIMIT {rowNum}")
        # データとカラムをpandasデータフレームに
        table= pd.DataFrame(data=cursor.fetchall(),columns=table_column_label)
        # カーソルの終了
        cursor.close()
    except Error as e:
        return pd.DataFrame([[""]*len(GUIcolumns)]*rowNum)
    
    # データが表示行数より少ないときの処理
    missing_rows = rowNum - len(table) # 不足している行数
    if missing_rows > 0:
        # 空白行の DataFrame を作成（""で埋める）
        empty_data = pd.DataFrame([[""] * len(table_column_label)] * missing_rows, columns=table_column_label)
        # 既存データに追加
        table = pd.concat([table, empty_data], ignore_index=True)
   
    return table







    
    
    
