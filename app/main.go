package main

import (
	"goserve/admin"
	"goserve/server"
)

func main() {
	admin.AdminPanel()
	server.CreateServer()
}
