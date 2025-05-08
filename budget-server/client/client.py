import requests

# res=requests.post('http://localhost:8000/budget/get', json={"user":"Tanaka", "amount":2000})
data={
    "user":"Tayama",
    "expense":2000,
    "year":2025,
    "month":12,
    "day":25,
    "description":"test",
    "calculation":True,
    "settlement":False,
    "fixed":False
}

res=requests.post('http://localhost:8000/budget/create-expense', json=data)

res=requests.get('http://localhost:8000/budget/get-expense-by-year', params={"year":2025})

json=res.json()
print(len(json))
print(json[0])
print(json[1])