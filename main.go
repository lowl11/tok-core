package main

import (
	"tok-core/src/api"
	"tok-core/src/definition"
)

func main() {
	definition.Init()
	api.StartServer()
}