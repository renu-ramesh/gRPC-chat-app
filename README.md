# gRPC-chat-app
The backend of a chat application, which makes it possible to connect to any room and
communicate directly with other users.

**Technologies Used**
- Golang
- Gorm
- gRPC
- Protobuf
- go gin

## Application Setup

To run the REST API endpoints :

    go run internal/main.go

To run the gRPC backend server :

    go run server/server.go


To run client server with group channels (one-to-many) :

    go run client/client.go -channel {channel_name} -username {sender_name} -server :5400

To run client server without group name specification (one-to-one) :

    go run client/client.go -username {sender_name} -server :5400

## Methods

Example methods as follows:

baseurl = 'http://localhost:8080/api/v1'

**Add new user**


    curl --request POST \
    --url {baseurl}/users \
    --header 'Content-Type: application/json' \
    --data '{"name":"renu"}'


**List all users**


    curl --request GET \
    --url {baseurl}/users

**Create a Group**


    curl --request POST \
    --url {baseurl}/channel \
    --header 'Content-Type: application/json' \
    --data '{"name":"Techversant",
    "type":"company"}'

**List all Groups**


    curl --request GET \
    --url {baseurl}/channel

**Delete a Group**


    curl --request DELETE \
    --url {baseurl}/channel/{:id} \
    --header 'Content-Type: application/json'

**Join a Group**


    curl --request POST \
    --url {baseurl}/users/{:id}/join \
    --header 'Content-Type: application/json' \
    --data '{"channel_name":"Techversant"}'

**List User's Group Details**


    curl --request GET \
    --url {baseurl}/users/channels

**Left from a Group**


    curl --request PUT \
    --url {baseurl}/users/{:id}/left \
    --header 'Content-Type: application/json' \
    --data '{"channel_name":"epixelsolutions"}'
