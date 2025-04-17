import tkinter as tk
from tkinter import ttk
import connection.connector as con
import communication.request as req
import threading
from dotenv import load_dotenv
import pandas as pd

# フレームのクラス
class App(tk.Frame):
    def __init__(self,root:tk.Tk):
        # フレームをrootの子として配置
        super().__init__(root)
        self.pack(padx=10,fill="both")
        
        # クラスの初期化
        # ウィジェットを初期化
        self.createWidget()
        # GUIの更新を開始
        self.updateGUI()
        # GUI情報の更新を開始
        self.updateData()
        
    # ウィジェットの初期化関数
    def createWidget(self):
        # ウィジェットの列は3つの領域に分轄
        for i in range(3):
            self.columnconfigure(i,weight=1)
  
        # サーバーステータス
        self.createServerStatus()
        # MySQLステータス
        self.createMySQLStatus()
        # サーバーログの生成
        self.createServerLogTable()
    
    # サーバーの稼働状態のモニターを作成
    def createServerStatus(self):
        # タイトルラベル
        self.serverStatusTitle=tk.Label(self,text="Server Status")
        self.serverStatusTitle.grid(row=0, column=0, sticky="w", padx=5, pady=5)
        # ステータス
        self.serverStatus="" # デフォルトの表示
        self.serverStatusText=tk.Label(self,text=self.serverStatus)
        self.serverStatusText.grid(row=0, column=1, sticky="w", padx=5, pady=5)
    
    # MySQLの稼働状態のモニターを作成
    def createMySQLStatus(self):
        # タイトルラベル
        self.MySQLStatusTitle=tk.Label(self,text="MySQL Status")
        self.MySQLStatusTitle.grid(row=1, column=0, sticky="w", padx=5, pady=5)
        # ステータス
        self.MySQLStatus="" # デフォルトの表示
        self.MySQLStatusText=tk.Label(self,text=self.MySQLStatus,font=("Arial",10))
        self.MySQLStatusText.grid(row=1, column=1, sticky="w", padx=5, pady=5)
        # MySQLのテーブルカラムとGUIのテーブルカラムの一致状況を通知
        self.MySQLColumnMatch="" #デフォルトの表示
        self.MySQLColumnMatchText=tk.Label(self,text=self.MySQLColumnMatch)
        self.MySQLColumnMatchText.grid(row=1, column=2, sticky="w", padx=5, pady=5)
    
    # サーバーログのテーブルを作成
    def createServerLogTable(self):
        self.tableColumns=["id","path","method","timestamp"] #カラムラベル
        self.columnsMinWidth=[40,200,60,150] # カラムの設定幅
        self.tableRowNum=15 # 表示するテーブル列数
        self.serverLogTable=pd.DataFrame([[""] * len(self.tableColumns)] * self.tableRowNum) # サーバーログの初期表示
        # テーブルの作成
        self.tree = ttk.Treeview(self,columns=self.tableColumns,show="headings",height=self.tableRowNum)
        # ヘッダーの追加
        for column,width in zip(self.tableColumns,self.columnsMinWidth):
            # カラム情報の追加
            self.tree.heading(column,text=column)
            self.tree.column(column,minwidth=width,width=width)
        self.tree.grid(row=2, column=0, columnspan=3, sticky="we", padx=5, pady=10)
        # 初期テーブルの作成
        for index, row in self.serverLogTable.iterrows():
            self.tree.insert("", "end", iid=index, values=tuple(row))

    # GUI表示の更新を実行する関数(描画を行う)
    def updateGUI(self):
        # サーバーステータス情報を更新
        self.serverStatusText.config(text=self.serverStatus)
        # MySQLのステータス情報を更新
        self.MySQLStatusText.config(text=self.MySQLStatus)
        self.MySQLColumnMatchText.config(text=self.MySQLColumnMatch)
        # テーブル情報を更新
        self.tree.delete(*self.tree.get_children()) # 既存テーブルの削除
        for index, row in self.serverLogTable.iterrows(): # 新規追加
            self.tree.insert("", "end", iid=index, values=tuple(row))
        # 次の更新もスケジュール
        self.after(3000,self.updateGUI)
    
    # GUIに表示するデータの更新を実行する関数
    def updateData(self):
        # サーバーステータス情報の更新
        self.updateServerStatus()
        # サーバーログテーブルとMySQLステータスの更新
        self.updateServerLog()

    # サーバーステータスの更新制御関数
    def updateServerStatus(self):
        # スレッドで非同期でサーバーステータスを取得
        threading.Thread(target=self.fetchLiveStatus,daemon=True).start()
    
    # ステータスを非同期で取得する関数
    def fetchLiveStatus(self):
        # サーバーのステータス情報を更新
        self.serverStatus="Server is running" if req.ServerStatus() else "Server is stopped"
        # 次の更新もスケジュール
        self.after(3000,self.updateServerStatus)

    # サーバーログテーブルとMySQLステータスの更新制御関数     
    def updateServerLog(self):
        # スレッドで非同期で最新情報を取得
        threading.Thread(target=self.fetchServerLog, daemon=True).start()

    def fetchServerLog(self):
        # MySQLとのコネクションを取得
        connection, err = con.createServerConnection()
        # MySQLステータスの更新
        MySQLStatus = "MySQL is stopped" if err != None else "MySQL is running"
        # 最新のサーバーログの取得
        serverLogTable = con.getServerLog(connection=connection, rowNum=self.tableRowNum, GUIcolumns=self.tableColumns, err=err)
        # カラム一致確認
        if serverLogTable.columns.to_list() != self.tableColumns:
            # 最新のテーブルカラムと既存のテーブルカラムが一致しなければサーバーログの表示を停止
            serverLogTable = pd.DataFrame([[""] * len(self.tableColumns)] * self.tableRowNum)
            # テーブルのカラム一致状況を更新
            MySQLColumnMatch = "COLUMN : MISMATCH"
        else:
            MySQLColumnMatch = "COLUMN : MATCH"
        # コネクションの停止
        if connection != None:
            connection.close()

        # 取得した最新情報をGUIに反映(描画はしない)
        self.serverLogTable=serverLogTable
        self.MySQLStatus=MySQLStatus
        self.MySQLColumnMatch=MySQLColumnMatch
        # 次の更新もスケジュール
        self.after(3000, self.updateServerLog)

        
# GUIの生成
def initGUI():
    # GUIの初期化
    root=tk.Tk()
    root.geometry("500x500")
    root.title("サーバー監視アプリ")
    # フレームの追加
    app=App(root=root)
    app.mainloop()

# .envファイルのロード
load_dotenv("../.env")
initGUI()