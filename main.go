package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/madjake/server"
	"github.com/madjake/webserver"
)

// Config encompasses top level application variables like interfaces & ports
type Config struct {
	webAddress  string
	gameAddress string
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	shutdown := make(chan bool)

	config := Config{
		getEnvironmentValue("IFF_WEB_ADDRESS", "localhost:8080"),
		getEnvironmentValue("IFF_GAME_ADDRESS", "localhost:8888"),
	}

	webServerAddress := flag.String("webServerAddress", config.webAddress, "http web server address")
	gameServerAddress := flag.String("gameServerAddress", config.gameAddress, "tcp game server address")

	go func() {

		ch := make(chan string)

		go func(ch chan string) {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("IF Game Terminal\nType HELP for supported commands\n---------------------------")

			for {
				fmt.Print("-> ")
				text, err := reader.ReadString('\n')

				if err != nil {
					close(ch)
					break
				}

				ch <- text
			}
		}(ch)

		for {
			select {
			case stdin, ok := <-ch:
				if !ok {
					break
				} else {
					text := strings.Replace(stdin, "\n", "", -1)

					//CommandParser + CommandHandler concepts needed to abstract this away
					if strings.Compare("who", text) == 0 {
						fmt.Println("Users connected: ", len(server.Users))
					}

					if strings.Compare("mem", text) == 0 {
						printMemUsage()
					}

					if strings.Compare("quit", text) == 0 {
						close(ch)
						shutdown <- true
						break
					}
				}
			}
		}
	}()

	go func() {
		doEvery(50*time.Millisecond, gameTick)
		shutdown <- true
	}()

	go func() {
		doEvery(10*time.Second, gameInfo)
		shutdown <- true
	}()

	go func() {
		fmt.Printf("Starting game server on %s", config.webAddress)
		server.NewGameServer(gameServerAddress)
		shutdown <- true
	}()

	go func() {
		fmt.Printf("Starting web server on %s", config.webAddress)
		webserver.NewWebServer(webServerAddress)
		shutdown <- true
	}()

	<-shutdown
	fmt.Println("Shutting down...")

}

func getEnvironmentValue(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func gameInfo(t time.Time) {
	//	fmt.Printf("%+v", server.Users)

}

func gameTick(t time.Time) {
	//fmt.Println("Tick")
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
