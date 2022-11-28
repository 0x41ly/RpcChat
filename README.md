# RpcChat

## Usage
- Server

    `go run Server.go [-Port <port>]`
- Client

    `go run Client.go -user <user> [-host <address:port>]`

- Inside the Client you can:
    - Send Private msg to a specific user:
        `/!tell <username> <msg>`
    - List Users:
        `/!list`
    - Send msg to all:
        `<msg>`
    - Logout:
        `/!logout`