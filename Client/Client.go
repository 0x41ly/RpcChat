package main

import (
	commons "RpcChat/Commons"
	"bufio"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// Globals/Constants
var (
	DEFAULT_PORT = 3410
	DEFAULT_HOST = "localhost"
)

func createClientFromFlags() (*commons.ChatClient, error) {
	var c *commons.ChatClient = &commons.ChatClient{}
	var host string

	flag.StringVar(&c.Username, "user", "", "Your username")
	flag.StringVar(&host, "host", "localhost", "The host you want to connect to")

	flag.Parse()

	if !flag.Parsed() {
		return c, errors.New("unable to create user from commandline flags. Please try again")
	}

	// Check for the structure of the flag to see if we can make any educated guesses for them
	if len(host) != 0 {

		if strings.HasPrefix(host, ":") { // Begins with a colon means :3410 (just Port)
			c.Address = DEFAULT_HOST + host
		} else if strings.Contains(host, ":") { // Contains a colon means host:Port
			c.Address = host
		} else { // Otherwise, it's just a host
			c.Address = net.JoinHostPort(host, strconv.Itoa(DEFAULT_PORT))
		}

	} else {
		c.Address = net.JoinHostPort(DEFAULT_HOST, strconv.Itoa(DEFAULT_PORT)) // Default to our default Port and host
	}

	return c, nil
}

func mainLoop(c *commons.ChatClient) {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error: %q\n", err)
		}

		line = strings.TrimSpace(line)
		params := strings.Fields(line)

		if strings.HasPrefix(line, "list") {
			c.List()
		} else if strings.HasPrefix(line, "tell") {
			c.Tell(params)
		} else if strings.HasPrefix(line, "logout") {
			c.Logout()
			break
		} else {
			c.Say(params)

		}
	}
}

func main() {
	// Set MAX PROCS
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Start by parsing any flags given to the program
	client, err := createClientFromFlags()
	if err != nil {
		log.Panicf("Error creating client from flags: %q", err)
	}

	client.Register()

	// Listen for messages
	go client.CheckMessages()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.Logout()
		os.Exit(1)
	}()

	mainLoop(client)
}
