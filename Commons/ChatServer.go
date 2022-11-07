package commons

import (
	"fmt"
	"log"
	"time"
)

type ChatServer struct {
	Port         string
	MessageQueue map[string][]string
	Users        []string
	KeepMeAlive  chan bool
	AllMsgs      []string
}

// Register registers a client username with the chat server.
// It sends a message to all Users notifying a user has joined
func (c *ChatServer) Register(username string, reply *string) error {
	*reply = "Welcome to GoChat v1.0!\n"
	*reply += "List of Users online:\n"

	c.Users = append(c.Users, username)
	c.MessageQueue[username] = c.AllMsgs

	for _, value := range c.Users {
		*reply += value + "\n"
	}

	for k := range c.MessageQueue {
		msg := time.Now().Format("2006/01/02 15:04:05") + ": " + username + " has joined."
		c.MessageQueue[k] = append(c.MessageQueue[k], msg)
		c.AllMsgs = append(c.AllMsgs, msg)
	}

	log.Printf("%s has joined the chat.\n", username)

	return nil
}

func (c *ChatServer) CheckMessages(username string, reply *[]string) error {
	*reply = c.MessageQueue[username]
	c.MessageQueue[username] = nil
	return nil
}

func (c *ChatServer) List(none Nothing, reply *[]string) error {
	*reply = append(*reply, "Current online Users:")

	*reply = append(*reply, c.Users...)

	log.Println("Dumped list of Users to client output")

	return nil
}

func (c *ChatServer) DoesUserExsist(user *string, exist *bool) error {
	*exist = false
	for i := range c.Users {
		if c.Users[i] == *user {
			*exist = true
			break
		}
	}
	return nil
}
func (c *ChatServer) Tell(msg Message, reply *Nothing) error {

	if queue, ok := c.MessageQueue[msg.Target]; ok {
		m := msg.User + " tells you " + msg.Msg
		c.MessageQueue[msg.Target] = append(queue, m)
	} else {
		m := msg.Target + " does not exist"
		c.MessageQueue[msg.User] = append(queue, m)
	}

	*reply = false

	return nil
}

func (c *ChatServer) Say(msg Message, reply *Nothing) error {

	m := msg.MsgDate.Format("2006/01/02 15:04:05") + ": " + msg.User + " says " + msg.Msg
	for k, v := range c.MessageQueue {
		c.MessageQueue[k] = append(v, m)
	}
	c.AllMsgs = append(c.AllMsgs, m)

	*reply = true

	return nil
}

func (c *ChatServer) Logout(username string, reply *Nothing) error {

	delete(c.MessageQueue, username)

	for i := range c.Users {
		if c.Users[i] == username {
			c.Users = append(c.Users[:i], c.Users[i+1:]...)
			break
		}
	}
	msg := time.Now().Format("2006/01/02 15:04:05") + ": " + username + " has logged out."
	for k, v := range c.MessageQueue {
		c.MessageQueue[k] = append(v, msg)
	}
	c.AllMsgs = append(c.AllMsgs, msg)

	fmt.Println("User " + username + " has logged out.")

	*reply = false

	return nil
}
