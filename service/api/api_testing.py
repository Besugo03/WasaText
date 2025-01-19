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

def test_getConversationInfo(userID, conversationID):
    url = f"{BASE_URL}/conversations/info/{conversationID}"
    data = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.get(url, headers=data)
    if response.status_code == 200:
        print("Get Conversation Info Test Passed")
        print("Response:", response.text)
    else:
        print("Get Conversation Info Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_createConversation_withDiuk(userID):
    url = f"{BASE_URL}/conversations/create"
    data = {
        "participants": ["diuk"],
        }
    headers = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.post(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Create Conversation Test Passed")
        print("Response:", response.text)
    else:
        print("Create Conversation Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)
    return response.text

def test_sendMessage(userID, conversationID, content, is_forwarded):
    url = f"{BASE_URL}/conversations/message"
    data = {
        "conversation_id": conversationID,
        "content": content,
        "is_forwarded": is_forwarded,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.post(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Send Message Test Passed")
        print("Response:", response.text)
    else:
        print("Send Message Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_addParticipant(userID, conversationID, participantID):
    url = f"{BASE_URL}/conversations/adduser"
    data = {
        "conversation_id": conversationID,
        "username" : participantID,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.put(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201 or response.status_code == 204:
        print("Add Participant Test Passed")
        print("Response:", response.text)
    else:
        print("Add Participant Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_removeParticipant(userID, conversationID, participantID):
    url = f"{BASE_URL}/conversations/kickuser"
    data = {
        "conversation_id": conversationID,
        "username" : participantID,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.delete(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Remove Participant Test Passed")
        print("Response:", response.text)
    else:
        print("Remove Participant Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_forwardMessage(userID, conversationID, messageID):
    url = f"{BASE_URL}/conversations/message/forward"
    data = {
        "conversation_id": conversationID,
        "message_id": messageID,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.post(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Forward Message Test Passed")
        print("Response:", response.text)
    else:
        print("Forward Message Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_deleteMessage(userID, messageID):
    url = f"{BASE_URL}/conversations/message"
    data = {
        "message_id": messageID,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.delete(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201 or response.status_code == 204:
        print("Delete Message Test Passed")
        print("Response:", response.text)
    else:
        print("Delete Message Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_changePicture(userID, newPicture):
    photourl = f"{BASE_URL}/user/photo"
    data = {
        "photo": newPicture,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    response = requests.put(photourl, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Change Username and Picture Test Passed")
        print("Response:", response.text)
    else:
        print("Change Username and Picture Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_changeUsername(userID, newUsername):
    url = f"{BASE_URL}/user/name"
    data = {
        "username": newUsername,
    }
    headers ={
        "Authorization" : "Bearer " + userID,
    }

    response = requests.put(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Change Username and Picture Test Passed")
        print("Response:", response.text)
    else:
        print("Change Username and Picture Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_changeGroupPhoto(userID, conversationID, newPhoto):
    url = f"{BASE_URL}/conversations/photo"
    data = {
        "conversation_id": conversationID,
        "photo": newPhoto,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    response = requests.put(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Change Group Photo Test Passed")
        print("Response:", response.text)
    else:
        print("Change Group Photo Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_changeGroupName(userID, conversationID, newName):
    url = f"{BASE_URL}/conversations/name"
    data = {
        "conversation_id": conversationID,
        "name": newName,
    }
    headers ={   
        "Authorization" : "Bearer " + userID,
    }

    response = requests.put(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Change Group Name Test Passed")
        print("Response:", response.text)
    else:
        print("Change Group Name Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_addCommentToMessage(userID, messageID, comment):
    url = f"{BASE_URL}/conversations/message/comment"
    data = {
        "message_id": messageID,
        "content": comment,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    response = requests.post(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Add Comment to Message Test Passed")
        print("Response:", response.text)
    else:
        print("Add Comment to Message Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_getCommentInfo(userID, commentID):
    url = f"{BASE_URL}/conversations/message/comment"
    json = {
        "comment_id": commentID,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    response = requests.get(url, headers=headers, json=json)
    if response.status_code == 200:
        print("Get Comment Info Test Passed")
        print("Response:", response.text)
    else:
        print("Get Comment Info Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_deleteComment(userID, commentID):
    url = f"{BASE_URL}/conversations/message/comment"
    data = {
        "comment_id": commentID,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    response = requests.delete(url, json=data, headers=headers)
    if response.status_code == 200 or response.status_code == 201:
        print("Delete Comment Test Passed")
        print("Response:", response.text)
    else:
        print("Delete Comment Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_getMessageInfo(userID, messageID):
    url = f"{BASE_URL}/conversations/message/info"
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    body = {
        "message_id": messageID,
    }

    response = requests.get(url, headers=headers, json=body)
    if response.status_code == 200:
        print("Get Message Info Test Passed")
        print("Response:", response.text)
    else:
        print("Get Message Info Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

if __name__ == "__main__":
    # Test "Hello World" endpoint
    test_hello_world()

    user1 = test_login("newTestUser")
    user2 = test_login("besugo")
    user3 = test_login("diuk")

    # test_getConversations(user1)
    # test_getConversations(user2)
    # test_getConversations(user3)

    # test_createConversation_withDiuk(user1)
    # test_getConversationInfo(user1, "1")

    # test_sendMessage(user1, "1", "message four", False)
    # test_sendMessage(user2, "1", "Hola, soy besugo", False)

    # test_removeParticipant(user1, "1", "diuk")
    # test_removeParticipant(user1, "2", "diuk")

    # test_addParticipant(user2, "1", "diuk")
    # test_removeParticipant(user2, "1", "diuk")

    # test_getConversations(user3) # questo non setta nullo come letto
    # test_getConversationInfo(user3, "1") # questo setta tutti i messaggi di quella conversazione come ricevuti
    # test_getMessageInfo(user3, "4") # questo setta il determinato messaggio come letto.
    

    # test_forwardMessage(user1, "1", "1") 
    # test_deleteMessage(user1, "1")

    # test_changePicture(user1, "https://i.imgur.com/testuno.jpg")
    # test_changeUsername(user1, "newTestUser")

    # test_changeGroupPhoto(user1, "1", "https://i.imgur.com/GROUPtest")
    # test_changeGroupName(user1, "1", "RRRrrrrepiteLoQueDicess")

    # test_addCommentToMessage(user1, "1", "comment to first message")
    # test_getCommentInfo(user1, "1")
    # test_deleteComment(user1, "1")