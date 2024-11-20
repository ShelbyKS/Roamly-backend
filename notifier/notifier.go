package notifier

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/ShelbyKS/Roamly-backend/notifier/config"
)

type Notifier struct {
	config  *config.Config
	clients map[string]*Client
}

type Client struct {
	SessionToken string
	MessageChan  chan []byte
}

func New(cfg *config.Config) *Notifier {
	return &Notifier{
		config: cfg,
	}
}

func (app *Notifier) Run() {
	broadcast := make(chan []byte)
	app.clients = make(map[string]*Client)

	go app.startKafkaConsumer(broadcast)
	go app.broadcastMessages(broadcast)

	r := app.newRouter()

	if err := r.Run(":" + app.config.ServerPort); err != nil {
		log.Fatalf("Failed to start notifier: %v", err)
	}
}

func (app *Notifier) startKafkaConsumer(broadcast chan []byte) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", app.config.KafkaConfig.Host, app.config.KafkaConfig.Port),
		"group.id":          app.config.KafkaConfig.Group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to init kafka consumer: %s\n", err)
	}
	defer consumer.Close()

	err = consumer.Subscribe(app.config.KafkaConfig.Topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe kafka topic: %s\n", err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			broadcast <- msg.Value
		} else {
			log.Printf("Failed to read message: %v\n", err)
		}
	}
}

func (app *Notifier) newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/notifications", app.websocketHandler)

	return router
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *Notifier) websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for websocket: %v", err)
		return
	}
	defer conn.Close()

	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session token"})
		c.Abort()
		return
	}

	client := &Client{
		SessionToken: sessionToken,
		MessageChan:  make(chan []byte),
	}
	app.clients[sessionToken] = client

	defer func() {
		delete(app.clients, sessionToken)
		close(client.MessageChan)
	}()

	for msg := range client.MessageChan {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			break
		}

		log.Printf("Sent message: %v", string(msg))
	}
}

type EventPayload struct {
	Action  string `json:"action"`
	Author  string `json:"author"`
	TripID  string `json:"trip_id"`
	Message string `json:"message"`
}

type KafkaMessage struct {
	Payload EventPayload `json:"payload"`
	Clients []string     `json:"clients"`
}

func (app *Notifier) broadcastMessages(broadcast chan []byte) {
	for {
		message := <-broadcast

		var kafkaMessage KafkaMessage
		if err := json.Unmarshal(message, &kafkaMessage); err != nil {
			log.Printf("Failed to decode message: %v", err)
			continue
		}

		payloadBytes, err := json.Marshal(kafkaMessage.Payload)
		if err != nil {
			log.Printf("Failed to encode payload: %v", err)
			continue
		}

		for _, actionParticipant := range kafkaMessage.Clients {
			client, exists := app.clients[actionParticipant]
			if !exists {
				continue
			}
			client.MessageChan <- payloadBytes
		}
	}
}
