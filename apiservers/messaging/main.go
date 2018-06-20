package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jadiego/react-slack-clone/apiservers/messaging/handlers"
	"github.com/jadiego/react-slack-clone/apiservers/messaging/models/messages"
)

const (
	defaultPort      = "80"
	defaultMongoPort = "27017"
)

const (
	apiRoot             = "/v1/"
	apiChannels         = apiRoot + "channels"
	apiMessages         = apiRoot + "messages"
	apiSpecificChannels = apiChannels + "/"
	apiSpecificMessages = apiMessages + "/"
)

var (
	port   = os.Getenv("PORT")
	host   = os.Getenv("HOST")
	dbAddr = os.Getenv("DBADDR")
)

func init() {
	if len(host) == 0 {
		fmt.Println("Host not set. Defaulting to empty host")
		host = ""
	}
	if len(port) == 0 {
		fmt.Println("Port number not set. Defaulting to port: " + defaultPort)
		port = defaultPort
	}
	if len(dbAddr) == 0 {
		fmt.Println("DB address not set. Defaulting to port: " + defaultMongoPort)
		dbAddr = defaultMongoPort
	}
}

func main() {
	addr := fmt.Sprintf("%s:%s", host, port)

	// set Messages DB
	mdbstore, err := messages.NewMongoStore(dbAddr, "", "", "")
	if err != nil {
		log.Fatalf("error starting messages DB: %v", err.Error())
	}
	defer mdbstore.Session.Close()

	// setup Handler context
	ctx := handlers.NewHandlerContext(mdbstore)

	// setup messages Routes
	mux := http.NewServeMux()
	mux.HandleFunc(apiChannels, ctx.ChannelsHandler)
	mux.HandleFunc(apiSpecificChannels, ctx.SpecificChannelHandler)
	mux.HandleFunc(apiMessages, ctx.MessageHandler)
	mux.HandleFunc(apiSpecificMessages, ctx.SpecificMessageHandler)

	fmt.Printf("Messages Route is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
