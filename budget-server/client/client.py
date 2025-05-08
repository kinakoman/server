import requests

# res=requests.post('http://localhost:8000/budget/get', json={"user":"Tanaka", "amount":2000})
res=requests.get('http://localhost:8000/budget/get?year=2025&month=12&day=25')


print(res.json())