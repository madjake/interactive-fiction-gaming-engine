package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

//Users is a list of users
var Users map[string]*GameUser

// GameUser describes a User connected to the game
type GameUser struct {
	IP         string
	Name       string
	Connection *websocket.Conn
}

// NewGameServer creates a TCP server that handles data sent from the FE game client
func NewGameServer(addr *string) {
	flag.Parse()
	log.SetFlags(0)

	Users = make(map[string]*GameUser)

	http.HandleFunc("/game", handleUser)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("Failed to upgrade web socket with client:", err)
		return
	}

	//	log.Printf("User connected from %s", c.RemoteAddr().String())

	user := addUser(c)

	defer c.Close()

	c.WriteMessage(1, []byte("Welcome to the server"))

	for {
		mt, messageBytes, err := user.Connection.ReadMessage()

		if err != nil {
			userQuitWithError(c, err)
			break
		}

		//		log.Printf("recv: %s", message)

		msg := string(messageBytes)
		selfMsg := fmt.Sprintf("You say: %s", msg)
		otherMsg := fmt.Sprintf("Somebody else says: %s", msg)

		err = sendMessageToUser(user, mt, selfMsg)
		sendMessageToAllUsersExcept([]*GameUser{user}, mt, otherMsg)

		if err != nil {
			userQuitWithError(c, err)
			break
		}
	}

	//log.Println("Disconnected")
}

func sendMessageToUser(user *GameUser, messageType int, message string) error {
	return user.Connection.WriteMessage(messageType, []byte(message))
}

func sendMessageToAllUsers(messageType int, message string) {
	for _, user := range Users {
		sendMessageToUser(user, messageType, message)
	}
}

func sendMessageToAllUsersExcept(excludedUsers []*GameUser, messageType int, message string) {
	for _, user := range Users {
		for _, excludedUser := range excludedUsers {
			if user != excludedUser {
				sendMessageToUser(user, messageType, message)
			}
		}
	}
}

func addUser(con *websocket.Conn) *GameUser {
	user := &GameUser{
		con.RemoteAddr().String(),
		"User",
		con,
	}

	Users[con.RemoteAddr().String()] = user

	return user
}

func userQuitWithError(con *websocket.Conn, err error) {
	log.Println("Socket Error:", err)
	// log or do something with error
	userQuit(con)
}

func userQuit(con *websocket.Conn) {
	// do stuff because user quit
	delete(Users, con.RemoteAddr().String())
}
