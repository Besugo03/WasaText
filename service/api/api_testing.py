import requests
import json

# URL del server
# BASE_URL = "http://localhost:6969"
BASE_URL = "http://localhost:3000"

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
        "conversationId": conversationID,
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
        "conversationId": conversationID,
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
        "conversationId": conversationID,
        "username" : participantID,
    }
    headers = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.put(url, json=data, headers=headers)
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
    url = f"{BASE_URL}/conversations/message/delete/{messageID}"
    headers = {
        "Authorization" : "Bearer " + userID,
    }
    response = requests.delete(url, headers=headers)
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
    if response.status_code == 200 or response.status_code == 201 or response.status_code == 204:
        print("Change Username and Picture Test Passed")
        print("Response:", response.text)
    else:
        print("Change Username and Picture Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_changeUsername(userID, newUsername):
    url = f"{BASE_URL}/user/username"
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
        "conversationId": conversationID,
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
    url = f"{BASE_URL}/conversations/message/info/{messageID}"
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    # body = {
    #     "message_id": messageID,
    # }

    # response = requests.get(url, headers=headers, json=body)
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        print("Get Message Info Test Passed")
        print("Response:", response.text)
    else:
        print("Get Message Info Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_getAllUsers(userID):
    url = f"{BASE_URL}/users"
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        print("Get All Users Test Passed")
        print("Response:", response.text)
    else:
        print("Get All Users Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)

def test_getUserExistence(userID, username):
    url = f"{BASE_URL}/users/exists/{username}"
    headers = {
        "Authorization" : "Bearer " + userID,
    }

    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        print("Get User Existence Test Passed")
        print("Response:", response.text)
    else:
        print("Get User Existence Test Failed")
        print("Status Code:", response.status_code)
        print("Response:", response.text)


if __name__ == "__main__":
    # Test "Hello World" endpoint
    test_hello_world()

    # user1 = test_login("bobbolino")
    # user2 = test_login("besugo")
    user3 = test_login("diuk")
    # user4 = test_login("vi")
    # user5 = test_login("johnny")
    # user6 = test_login("panam")
    # user7 = test_login("riccardo")

    # test_getAllUsers(user1)
    # test_getUserExistence(user1, "diuk")

    # test_getConversations(user1)
    # test_getConversations(user2)
    # test_getConversations(user3)

    # test_createConversation_withDiuk(user7)

    # test_getMessageInfo(user3, "24")

    # for i in range (1, 10):
        # test_sendMessage(user7, "1", f"messaggio di test numero {i}", False)
    
    # test_sendMessage(user3, "3", "https://i.imgflip.com/1m36cf.jpg", False)
    # test_sendMessage(user3, "3", "EH! https://i.imgflip.com/1m36cf.jpg VOLEVI! https://i.imgflip.com/1m36cf.jpg MUZUNNA RAHAZZI! https://i.imgflip.com/1m36cf.jpg https://i.imgflip.com/1m36cf.jpg https://i.imgflip.com/1m36cf.jpg", False)
    # test_getConversationInfo(user3, "5")

    # test_addParticipant(user3, "1", "bobbolino")

    # test_removeParticipant(user1, "1", "diuk")
    # test_removeParticipant(user3, "2", "newTestUser")

    # test_addParticipant(user3, "3", "panam")
    # test_removeParticipant(user2, "1", "diuk")

    # test_getConversations(user3) # questo non setta nullo come letto
    # test_getConversationInfo(user3, "1") # questo setta tutti i messaggi di quella conversazione come ricevuti
    # test_getMessageInfo(user3, "4") # questo setta il determinato messaggio come letto.

    # test_forwardMessage(user1, "1", "1")
    # test_deleteMessage(user3, "5")

    # test_changePicture(user3, "https://media.tenor.com/lll_dEVH2s4AAAAi/nagatoro-nekotoro.gif")
    # test_changePicture(user7, "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fstatic.wikia.nocookie.net%2Fendaven%2Fimages%2Fe%2Fe3%2FAstolfo_pfp.jpg%2Frevision%2Flatest%2Fscale-to-width-down%2F1200%3Fcb%3D20230429063130&f=1&nofb=1&ipt=e5769024e356248b7d6a8924486200bad278ece3959877e74e07ab8799a38f8d&ipo=images")
    # test_changeUsername(user1, "bobbolino")

    # test_changeGroupPhoto(user1, "1", "https://i.imgur.com/GROUPtest")
    # test_changeGroupName(user1, "3", "RRRrrrrepiteLoQueDicess")

    # test_addCommentToMessage(user1, "1", "comment to first message")
    # test_getCommentInfo(user1, "1")
    # test_deleteComment(user1, "1")
    
    # test_getConversationInfo(user3, "3")