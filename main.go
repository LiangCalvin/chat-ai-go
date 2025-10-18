package main

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pusher/pusher-http-go/v5"
)

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	IsUser    bool      `json:"isUser"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	messages []Message
	mu       sync.Mutex
)
func main() {
    // Initialize a new Fiber app
    app := fiber.New()

	app.Use(cors.New())

	pusherClient := pusher.Client{
		AppID: "2065516",
		Key: "82571582980be44b96ec",
		Secret: "0449218ee44cbf94531a",
		Cluster: "ap1",
		Secure: true,
	}

    // Define a route for the GET method on the root path '/'
    app.Get("/messages", func(c *fiber.Ctx) error {
        // Send a string response to the client
        // return c.SendString("Hello, World 👋!")
		mu.Lock()
		defer mu.Unlock()
		return c.JSON(messages)
    })

	app.Post("/messages", func(c *fiber.Ctx) error {
		// var data map[string]string

		// if err := c.BodyParser(&data); err != nil{
		// 	return err
		// }

		// pusherClient.Trigger("chat", "message", data)
        
		// return c.JSON([]string{})
		
		var msg Message
		if err := c.BodyParser(&msg); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if msg.Timestamp.IsZero() {
			msg.Timestamp = time.Now()
		}

		mu.Lock()
		messages = append(messages, msg)
		mu.Unlock()

		aiResponse := Message{
		ID:        "ai-" + msg.ID,
		Content:   "Hi!" + msg.Content + "This is a very long AI response. " + 
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
			"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
			"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. " +
			"…(repeat to exceed 500 chars)…",
		IsUser:    false,
		Timestamp: time.Now(),
	}
// Save AI response
	mu.Lock()
	messages = append(messages, aiResponse)
	mu.Unlock()

	// Trigger pusher events for both messages
	pusherClient.Trigger("chat", "message", msg)        // user message
	pusherClient.Trigger("chat", "message", aiResponse) // AI response

	// Return the AI response in the POST response
	return c.JSON(aiResponse)
	
		// pusherClient.Trigger("chat", "message", msg)

		// return c.JSON(msg)
    })

    // Start the server on port 3000
    log.Fatal(app.Listen(":3000"))
}