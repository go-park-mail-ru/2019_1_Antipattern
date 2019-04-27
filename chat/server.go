package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func JWTParse(w http.ResponseWriter, r *http.Request) (string, error) {
	secret := []byte("secret")
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", nil
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, ok := claims["uid"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return "", err
		}
		if uid == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return "", err
		}
		return uid, nil
	}
	return "", nil
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
	var clients []*Client
	for {
		select {
		case newClient := <-clientChan:
			clients = append(clients, newClient)
			fmt.Println("Client joined")
		case message := <-messageChan:
			for _, client := range clients {
				if client.isConnected {
					go client.SendMessage(message)
				} else {
					// Delete client
					//clients[index] = clients[len(clients)-1]
					//clients = clients[:len(clients)-1]
				}
			}
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
	uid, err := JWTParse(w, r)
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
