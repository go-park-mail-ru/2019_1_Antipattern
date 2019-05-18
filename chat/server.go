package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"../providers/auth"
	"../providers/user_data"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	ID   string `json:"_id,omitempty"`
	UID  string `json:"uid"`
	Text string `json:"text"`
}

type Client struct {
	isConnected bool
	uid         string
	conn        *websocket.Conn
}
type MessageJSON struct {
	Status  string    `json:"status"`
	Payload []Message `json:"payload"`
}

var authProvider auth.JWTProvider = auth.JWTProvider{
	ServerAddress: "identity_service:8081",
	Secure:        false,
	AuthDomain:    ".kpacubo.xyz",
}

var apiProvider user_data.GrpcProvider = user_data.GrpcProvider{
	ServerAddress: "api:8081",
}

func ParseAuth(r *http.Request) (string, error) {
	return authProvider.GetUUID(r)
}

var _client *mongo.Client

func dbConnect() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if _client == nil {
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://chat_db:27017"))
		if err != nil {
			return nil, err
		}
		_client = client
		err = _client.Connect(ctx)
		if err != nil {
			return nil, err
		}
	}
	err := _client.Ping(ctx, nil)
	return _client, err
}

func (client *Client) ReceiveMessage(messageChan chan *Message) {
	defer func() {
		client.isConnected = false
		client.conn.Close()
	}()
	for {
		message := Message{}
		err := client.conn.ReadJSON(&message)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
		dbClient, err := dbConnect()
		if err != nil {
			fmt.Println("Failed to connect DB")
			return
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		collection := dbClient.Database("kpacubo").Collection("messages")

		message.UID = client.uid
		result, err := collection.InsertOne(ctx, message)
		if err != nil {
			fmt.Println("Failed to create message")
		}
		message.ID = result.InsertedID.(primitive.ObjectID).Hex()
		messageChan <- &message
	}
}

func (client *Client) SendMessage(message *Message) {
	client.conn.WriteJSON(*message)
}

func ChatRoom(clientChan chan *Client, messageChan chan *Message) {
	var mtx sync.Mutex
	clients := make(map[string]*Client)
	for {
		select {
		case newClient := <-clientChan:
			clients[newClient.uid] = newClient
			newClient.conn.SetCloseHandler(func(code int, text string) error {
				mtx.Lock()
				defer mtx.Unlock()
				delete(clients, newClient.uid)
				return nil
			})
			fmt.Println("Client joined")
			userData, err := apiProvider.GetUsers([]string{newClient.uid})
			if err != nil && len(userData) != 0 {
				for _, client := range clients {
					if client.isConnected {
						message := Message{UID: "", Text: userData[0].Login + " joined"}
						go client.SendMessage(&message)
					}
				}
			}
		case message := <-messageChan:
			mtx.Lock()

			for _, client := range clients {
				if client.isConnected {
					go client.SendMessage(message)
				}
			}
			mtx.Unlock()
		}
	}
}

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	dbClient, err := dbConnect()
	if err != nil {
		fmt.Println("Failed to connect DB")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := dbClient.Database("kpacubo").Collection("messages")
	options := options.Find()
	options.SetLimit(int64(50)).SetSort(bson.M{"_id": -1})

	cursor, err := collection.Find(ctx, bson.D{}, options)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var messages []Message
	for cursor.Next(ctx) {
		m := Message{}
		err = cursor.Decode(&m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		messages = append(messages, m)
	}
	messageJSON := MessageJSON{
		Status:  "success",
		Payload: messages,
	}
	json, _ := json.Marshal(messageJSON)
	w.Write(json)

}
func upgraderHandler(w http.ResponseWriter, r *http.Request, clientChan chan *Client, messageChan chan *Message) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// TODO: Implement
			return true
		},
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := Client{
		true,
		"",
		connection,
	}
	uid, err := ParseAuth(r)
	if err == nil {
		client.uid = uid
	}

	go client.ReceiveMessage(messageChan)
	clientChan <- &client
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func InitDB() {
	dbClient, err := dbConnect()
	if err != nil {
		fmt.Println("Failed to connect DB")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := dbClient.Database("kpacubo").Collection("messages")
	m := Message{Text: "Server started!", UID: ""}
	collection.InsertOne(ctx, m)
}

func main() {
	messageChan := make(chan *Message)
	clientChan := make(chan *Client)
	go ChatRoom(clientChan, messageChan)
	InitDB()

	http.HandleFunc("/messages", HandleGetMessages)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgraderHandler(w, r, clientChan, messageChan)
	})

	http.ListenAndServeTLS(":2000", "/cert/live/kpacubo.xyz/fullchain.pem", "/cert/live/kpacubo.xyz/privkey.pem", nil)

}
