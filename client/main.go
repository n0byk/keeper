package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/n0byk/keeper/client/internal"
)

var done chan interface{}
var interrupt chan os.Signal
var wsconn *websocket.Conn

var isConnected bool

var addr = flag.String("addr", "localhost:3001", "http service address")
var u = url.URL{}

func main() {
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Keeper Добро пожаловать!")
	printMenu()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	done = make(chan interface{})
	interrupt = make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	defer conn.Close()
	if err != nil {
		log.Fatal("dial:", err)
	}
	wsconn = conn
	if err != nil {
		isConnected = false
	} else {
		isConnected = true
	}

	go recieveHandler()

	scanner.Scan()

	for scanner.Text() != "/quit" {

		cmd, msg := sanitizeInput(scanner.Text())
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		switch cmd {
		case "err":
		case internal.Register:
			registerAction(msg)
		case internal.Auth:
			authAction(msg)
		case internal.LobbyCreate:
			CreateLobbyAction(msg)
		case internal.GetPublicToken:
			GetPublicTokenAction(msg)
		case internal.GetLobbyId:
			GetLobbyIdAction(msg)
		case internal.LobbyInvite:
			LobbyInviteAction(msg)

		}
		scanner.Scan()
	}
}

func quit() {
	log.Printf("Bey friend")
	os.Exit(0)
}

func recieveHandler() {
	defer close(done)

	for {
		// err = wsconn.ReadJSON(&msg)
		_, message, err := wsconn.ReadMessage()
		if err != nil {
			log.Println("Error reading json: ", err)
		}
		log.Println(string(message))

	}
}

func sanitizeInput(userIn string) (string, string) {
	var command string
	var message string
	if userIn != "" && userIn[0] == '/' {
		result := strings.SplitN(userIn, " ", 2)
		command = result[0]
		if len(result) == 2 {
			message = result[1]
		}

		command = command[1:]
	} else {
		command = ""
		message = userIn
	}
	return command, message
}

func registerAction(message string) {
	words := strings.Fields(message)
	msg := internal.TRegistrationRequest{
		Action:   internal.Register,
		Login:    words[0],
		Password: words[1],
	}
	err := wsconn.WriteJSON(msg)
	if err == nil {
		return
	}
	fmt.Println(`Вы зарегистрировались , Ваш логин: ` + words[0])
}

func authAction(message string) {
	words := strings.Fields(message)
	msg := internal.TAuthRequest{
		Action:   internal.Auth,
		Login:    words[0],
		Password: words[1],
	}
	err := wsconn.WriteJSON(msg)
	if err == nil {
		return
	}
}

func CreateLobbyAction(message string) {
	words := strings.Fields(message)
	msg := internal.TCreateLobbyRequest{
		Action:    internal.LobbyCreate,
		Token:     words[0],
		LobbyName: words[1],
	}
	err := wsconn.WriteJSON(msg)
	if err == nil {
		return
	}
}

func GetPublicTokenAction(message string) {
	msg := internal.TGetPublicTokenRequest{
		Action: internal.GetPublicToken,
		Login:  message,
	}
	err := wsconn.WriteJSON(msg)
	if err == nil {
		return
	}
}

func GetLobbyIdAction(message string) {
	msg := internal.TLobbyIdRequest{
		Action:    internal.GetLobbyId,
		LobbyName: message,
	}
	err := wsconn.WriteJSON(msg)
	if err == nil {
		return
	}
}

func LobbyInviteAction(message string) {
	words := strings.Fields(message)

	msg := internal.TAddToLobbyRequest{
		Action:      internal.LobbyInvite,
		LobbyName:   words[0],
		Token:       words[1],
		PublicToken: words[2],
	}

	err := wsconn.WriteJSON(msg)
	if err == nil {
		return
	}
}

func SetDataAction(message string) {
	words := strings.Fields(message)

	msg := internal.TSetDataRequest{
		Action:    internal.SetData,
		Token:     words[0],
		LobbyName: words[1],
		LobbyId:   words[2],
		RowData:   words[3],
	}

	err := wsconn.WriteJSON(msg)
	if err == nil {
		return
	}
}

func printMenu() {
	fmt.Println(">1. Введите логин и пароль через пробел для регистрации пример: \"/registration login password\"")
	fmt.Println(">2. Авторизация в системе пример: \"/auth login password\"")
	fmt.Println(">3. Создание комнаты, для шеринга приватных данных пример: \"/create_lobby token lobby_name\"")
	fmt.Println(">4. Получение публичного ключа пользователя пример: \"/get_public_token login\"")
	fmt.Println(">5. Получение Id Лобби пример : \"/lobby_id lobby_name\"")
	fmt.Println(">6. Пригласить пользователя в лобби, дать другому пользователю права читать сообщения : \"/invite_lobby lobby_name token PublicToken\"")
	fmt.Println(">7. Добавление данных в систему : \"/set_data token lobby_name lobby_id data\"")
	fmt.Println(">8. Выход из приложения пример: \"/quit\"")
}
