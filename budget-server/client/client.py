import requests

res=requests.post('http://localhost:8000/user/', json={"name":"Tanaka", "age":40})

print(res.json())