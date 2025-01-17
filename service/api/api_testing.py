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

def test_login(username):
    url = f"{BASE_URL}/session"
    data = {
        "username": f"{username}",
    }
    response = requests.post(url, json=data)
    if response.status_code == 200 or response.status_code == 201:
        print("Login Test Passed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)
    else:
        print("Login Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)
    return response.json()["identifier"]

def test_getConversations(userID):
    url = f"{BASE_URL}/conversations"
    data = {
        "Authorization" : "Bearer " + userID,
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

    user1 = test_login("testuser")
    user2 = test_login("besugo")
    user3 = test_login("diuk")
    user4 = test_login("isacco")
    user5 = test_login("bombolino")

    test_getConversations(user1)
    test_getConversations(user2)
    test_getConversations(user3)
    test_getConversations(user4)
    test_getConversations(user5)