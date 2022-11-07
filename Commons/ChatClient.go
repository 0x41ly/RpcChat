package commons

import (
	"bufio"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type Nothing bool

type ChatClient struct {
	Username string
	Address  string
	Client   *rpc.Client
}

func (c *ChatClient) getClientConnection() *rpc.Client {
	var err error

	if c.Client == nil {
		c.Client, err = rpc.Dial("tcp", c.Address)
		if err != nil {
			log.Panicf("Error establishing connection with host: %q", err)
		}
	}

	return c.Client
}

// Register takes a username and registers it with the server
func (c *ChatClient) Register() {
	var reply string
	c.Client = c.getClientConnection()
	for {
		if c.DoesUserExsist(&c.Username) || c.Username == "" {
			print("Username already exists. Please re-enter another username:")
			reader := bufio.NewReader(os.Stdin)
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error: %q\n", err)
			}
			c.Username = strings.TrimSuffix(line, "\n")
		} else {
			break
		}
	}
	err := c.Client.Call("ChatServer.Register", c.Username, &reply)
	if err != nil {
		log.Printf("Error registering user: %q", err)
	} else {
		log.Printf("Reply: %s", reply)
	}
}
func (c *ChatClient) DoesUserExsist(user *string) bool {
	var exist bool
	c.Client = c.getClientConnection()
	err := c.Client.Call("ChatServer.DoesUserExsist", user, &exist)
	if err != nil {
		log.Printf("Error Checking user existance : %q", err)
		return true
	} else {
		return exist
	}
}

// CheckMessages does a check every second for new messages for the user
func (c *ChatClient) CheckMessages() {
	var reply []string
	c.Client = c.getClientConnection()

	for {
		err := c.Client.Call("ChatServer.CheckMessages", c.Username, &reply)
		if err != nil {
			log.Fatalln("Chat has been shutdown. Goodbye.")
		}

		for i := range reply {
			println(reply[i])
		}

		time.Sleep(time.Second)
	}
}

// List lists all the Users in the chat currently
func (c *ChatClient) List() {
	var reply []string
	var none Nothing
	c.Client = c.getClientConnection()

	err := c.Client.Call("ChatServer.List", none, &reply)
	if err != nil {
		log.Printf("Error listing Users: %q\n", err)
	}

	for i := range reply {
		log.Println(reply[i])
	}
}

// Tell sends a message to a specific user
func (c *ChatClient) Tell(params []string) {
	var reply Nothing
	c.Client = c.getClientConnection()

	if len(params) > 2 {
		msg := strings.Join(params[2:], " ")
		message := Message{
			User:    c.Username,
			Target:  params[1],
			Msg:     msg,
			MsgDate: time.Now(),
		}

		err := c.Client.Call("ChatServer.Tell", message, &reply)
		if err != nil {
			log.Printf("Error telling Users something: %q", err)
		}
	} else {
		log.Println("Usage of tell: tell <user> <msg>")
	}
}

// Say sends a message to all Users
func (c *ChatClient) Say(params []string) {
	var reply Nothing
	c.Client = c.getClientConnection()

	msg := strings.Join(params, " ")
	message := Message{
		User:    c.Username,
		Target:  "",
		Msg:     msg,
		MsgDate: time.Now(),
	}

	err := c.Client.Call("ChatServer.Say", message, &reply)
	if err != nil {
		log.Printf("Error saying something: %q", err)
	}

}

// Logout logs out the current user and shuts down the client
func (c *ChatClient) Logout() {
	var reply Nothing
	c.Client = c.getClientConnection()

	err := c.Client.Call("ChatServer.Logout", c.Username, &reply)
	if err != nil {
		log.Printf("Error logging out: %q", err)
	}
}
