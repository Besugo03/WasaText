import requests
import json

# URL del server
BASE_URL = "http://localhost:6969"

def test_hello_world():
    url = f"{BASE_URL}/"
    response = requests.get(url)
    if response.status_code == 200:
        print("Hello World Test Passed")
        print("Response:", response.text)
    else:
        print("Hello World Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_login():
    url = f"{BASE_URL}/session"
    data = {
        "name": "testuser"
    }
    response = requests.post(url, json=data)
    if response.status_code == 200:
        print("Login Test Passed")
        print("Response:", response.text)
    else:
        print("Login Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)
    return response.json()["identifier"]

def test_getConversations():
    url = f"{BASE_URL}/users/testuser/conversations"
    data = {
        "Authorization" : userID,
    }
    response = requests.get(url, headers=data)
    if response.status_code == 200:
        print("Get Conversations Test Passed")
        print("Response:", response.text)
    else:
        print("Get Conversations Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)
    return response.text

if __name__ == "__main__":
    # Test "Hello World" endpoint
    test_hello_world()

    userID = test_login()

    test_getConversations()