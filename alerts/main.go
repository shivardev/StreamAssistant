package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"myproject/utils"
	"strings"
	"sync"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type MessagePayload struct {
	Messages []utils.ChatMessage `json:"messages"`
}
type GifsPayload struct {
	GifUrl string `json:"url"`
}

var messageQueue chan utils.ChatMessage
var speakQueue chan utils.ChatMessage
var userDatabaseQueue chan utils.ChatMessage
var gifQueue chan GifsPayload
var connections = make(map[*websocket.Conn]bool)
var mu sync.Mutex // To ensure thread-safety
var Users = hashset.New()
var relangiData = utils.GetRelangiJSON()
var messageHandlers = map[string]func(){}
var videoLink string

func init() {
	messageQueue = make(chan utils.ChatMessage, 100)
	speakQueue = make(chan utils.ChatMessage, 100)
	gifQueue = make(chan GifsPayload, 100)
	userDatabaseQueue = make(chan utils.ChatMessage, 100)
}

func processMsgQueue() {
	for {
		msg := <-messageQueue
		println(msg.MessageContent, msg.AuthorName)
		if utils.IgnoredUsers.Contains(msg.AuthorName) {
			continue
		}
		if strings.Contains(msg.MessageContent, "frog") {
			utils.DoAction(utils.GetAction(string(utils.Frog)))
		} else if strings.Contains(msg.MessageContent, "iron") {
			utils.DoAction(utils.GetAction(string(utils.Ironman)))
		} else if strings.Contains(msg.MessageContent, "bat") {
			fmt.Println("Batman")
			utils.DoAction(utils.GetAction(string(utils.Batman)))
		} else if strings.Contains(msg.MessageContent, "joke") || strings.Contains(msg.MessageContent, "clown") {
			utils.DoAction(utils.GetAction(string(utils.Clown)))
		} else if strings.Contains(msg.MessageContent, "eye") {
			utils.DoAction(utils.GetAction(string(utils.Eyes)))
		} else if strings.Contains(msg.MessageContent, "thug") {
			utils.DoAction(utils.GetAction(string(utils.Thug)))
		}
		if len(msg.MessageContent) > 6 && (strings.HasPrefix(msg.MessageContent, "!speak") || strings.HasPrefix(msg.MessageContent, "! speak")) {
			speakQueue <- msg
		}

		lowerMessage := strings.ToLower(msg.MessageContent)
		for command, response := range utils.CommandHandler {
			if strings.HasPrefix(lowerMessage, command) {
				fmt.Println(response)
				utils.SendMsgToYoutube(response)
				continue
			}
		}
		userDatabaseQueue <- msg
	}
}

// func initMessageHandlers() {
// 	val := reflect.ValueOf(relangiData)
// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Type().Field(i)
// 		fieldName := field.Tag.Get("json")
// 		messageHandlers[fieldName] = func(f reflect.Value) func() {
// 			return func() {
// 				if f.Len() > 0 {
// 					utils.SendMsgToYoutube(f.Index(rand.Intn(f.Len())).Interface().(string))
// 				}
// 			}
// 		}(val.Field(i))
// 	}
// }

// Define structs to match the JSON response structure
type Action struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Group           string `json:"group"`
	Enabled         bool   `json:"enabled"`
	SubactionsCount int    `json:"subactions_count"`
}

type Response struct {
	Count   int      `json:"count"`
	Actions []Action `json:"actions"`
}

func processSpeakQueue() {
	for {
		msg := <-speakQueue

		if len(msg.MessageContent) > 6 && strings.HasPrefix(msg.MessageContent, "!speak") {
			fmt.Println("Processing !speak message:", msg.MessageContent)
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshalling message to JSON:", err)
				continue
			}
			mu.Lock()
			for conn := range connections {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonMsg)); err != nil {
					log.Println("WebSocket write error:", err)
					conn.Close()
					delete(connections, conn)
				}
			}
			mu.Unlock()
		}

	}
}

// For each User youtube comment msg, will be handling the databases and if the points are a multiple of 10, 20 then will trigger a congrats msg to youtube.
// also processes to store the information in Sqlite db.
func processUserPoints() {
	for {

		msg := <-userDatabaseQueue
		fmt.Println("Processing User", msg.AuthorName)
		var user utils.User
		exists, user, err := utils.CheckUserExists(msg.AuthorId)
		if err != nil {
			fmt.Println("Error checking user existence: ", err)
			continue
		}
		if !exists {
			fmt.Println("User doesn't exist, inserting")
			user, err = utils.InsertUser(msg)
			if err != nil {
				fmt.Println("Error Inserting user: ", err)
				continue
			}
		}
		if exists {
			fmt.Println("User exists, updating")
			user, err = utils.UpdateUser(msg, user)
			if err != nil {
				fmt.Println("Error fetching user: ", err)
				continue
			}
		}
		fmt.Println("User: ", user.UserName, " -> ", user.Points)

		if user.Points == 1 {
			utils.SendMsgToYoutube(fmt.Sprintf(utils.WelcomeMsgs[rand.Intn(len(utils.WelcomeMsgs))], msg.AuthorName))
			Users.Add(msg.AuthorId)
		} else if user.Points%10 == 0 {
			utils.SendMsgToYoutube(fmt.Sprintf(utils.CongratsMessages[rand.Intn(len(utils.CongratsMessages))], msg.AuthorName, user.Points))
		} else if strings.HasPrefix(msg.MessageContent, "!points") {
			utils.SendMsgToYoutube(fmt.Sprintf(utils.PointsMsgs[rand.Intn(len(utils.PointsMsgs))], msg.AuthorName, user.Points))
		} else {
			if !utils.IgnoredUsers.Contains(msg.AuthorName) && !Users.Contains(msg.AuthorId) {
				Users.Add(msg.AuthorId)
				utils.SendMsgToYoutube(fmt.Sprintf(utils.SubscriberMsgs[rand.Intn(len(utils.SubscriberMsgs))], msg.AuthorName, user.Points))
			}
		}
		go func() {
			utils.SetLeaderBoard(user)
		}()
		fmt.Println("User: ", user.UserName, "Processed")
	}
}
func main() {
	utils.GetActionList()
	utils.DataBaseConnection()
	// initMessageHandlers()
	go processMsgQueue()
	go processSpeakQueue()
	go processUserPoints()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type,Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}))
	app.Static("/", "./static")
	app.Get(("/"), func(c *fiber.Ctx) error {
		return c.SendString("Alive")
	})
	app.Get("/speak", func(c *fiber.Ctx) error {
		return c.SendFile("./static/speak.html")
	})
	app.Get("/streamReset", func(c *fiber.Ctx) error {
		Users.Clear()
		fmt.Println("Stream reset")
		return c.SendString("Stream reset successfully")
	})
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Post("/streamurl", func(c *fiber.Ctx) error {
		var streamUrl utils.StreamUrl

		if err := c.BodyParser(&streamUrl); err != nil {
			fmt.Println("Error parsing body:", err)
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
		}
		if videoLink != streamUrl.Url {
			fmt.Println("Stream url changed, resettings data!")
			// stream url changed
			reset()
		}
		videoLink = streamUrl.Url
		return c.Status(fiber.StatusOK).SendString("Stream url received successfully")
	})
	app.Get("/streamurl", func(c *fiber.Ctx) error {
		return c.SendString(videoLink)
	})
	app.Post("/takestats", func(c *fiber.Ctx) error {
		// will takes current stats from youtube live stream contains live likes and viewer count of the stream.
		var statPayload utils.StatsPayload

		if err := c.BodyParser(&statPayload); err != nil {
			fmt.Println("Error parsing body:", err)
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
		}

		likeCompare := statPayload.Stats.Likes - statPayload.Stats.PreviousLikes
		if statPayload.Stats.ShouldCongratulate && statPayload.Stats.MaxLikes > 0 {
			utils.SendMsgToYoutube(fmt.Sprintf("Congrats! Stream has reached %d likes", statPayload.Stats.MaxLikes))
		} else if likeCompare > 0 {
			utils.DoAction(utils.GetAction(string(utils.Likes)))
			fmt.Println("Thanks for liking the stream!")
			// utils.SendMsgToYoutube(fmt.Sprint("Thanks for liking the stream!"))
		} else if likeCompare < 0 {
			fmt.Println("Someone disliked the stream!")
			utils.DoAction(utils.GetAction(string(utils.Cry)))
			// utils.SendMsgToYoutube(fmt.Sprint("Someone removed a like!"))
		}
		return c.Status(fiber.StatusOK).SendString("Stats received successfully")
	})
	app.Post("/takemsgs", func(c *fiber.Ctx) error {
		fmt.Println("Processing messages")
		var chatMessages MessagePayload

		if err := c.BodyParser(&chatMessages); err != nil {
			fmt.Println("Error parsing body:", err)
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
		}
		go func() {
			for _, msg := range chatMessages.Messages {
				messageQueue <- msg
			}
		}()

		return c.Status(fiber.StatusOK).SendString("Messages received successfully")
	})
	app.Post("/takegifs", func(c *fiber.Ctx) error {
		var eachGif GifsPayload

		if err := c.BodyParser(&eachGif); err != nil {
			fmt.Println("Error parsing body:", err)
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
		}
		fmt.Println(eachGif)
		go func() {
			gifQueue <- eachGif
		}()

		return c.Status(fiber.StatusOK).SendString("Gif received successfully")
	})
	type InstagramNotification struct {
		Platform string `json:"platform"`
		UserName string `json:"userName"`
		Action   string `json:"action"`
	}
	app.Get("/leaderboard", func(c *fiber.Ctx) error {

		isVertical := c.Query("isVertical", "false") == "true"
		topTwo := utils.GetTopLeaderBoard(isVertical)
		fmt.Println(topTwo)
		if len(topTwo) == 0 {
			return c.JSON([]map[string]interface{}{})
		}
		leaderboard := []map[string]interface{}{}
		for i := 0; i < len(topTwo); i++ {
			leaderboard = append(leaderboard, map[string]interface{}{
				"userName": topTwo[i].UserName,
				"points":   topTwo[i].Points,
			})
		}

		return c.JSON(leaderboard)
	})
	app.Post("/insta", func(c *fiber.Ctx) error {
		var notification InstagramNotification

		if err := json.Unmarshal(c.Body(), &notification); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
		}
		fmt.Println("Platform:", notification.Platform, "UserName:", notification.UserName, "Action:", notification.Action)
		if notification.Platform == "instagram" {

		}
		return c.Status(fiber.StatusOK).SendString("Notification received successfully")
	})
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		mu.Lock()
		connections[c] = true // Track the connection
		mu.Unlock()

		defer func() {
			mu.Lock()
			delete(connections, c) // Clean up when the connection closes
			mu.Unlock()
			c.Close()
		}()

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			// Echo the message back to the client
			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))
	log.Fatal(app.Listen(":3000"))
}
func reset() {
	Users.Clear()
	utils.ResetLeaderBoard()
	fmt.Println("Stream data reset")
}
