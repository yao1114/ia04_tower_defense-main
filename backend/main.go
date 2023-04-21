package main

import (
	"fmt"

	"gitlab.utc.fr/michenwe/ia04_tower_defense/backend/agt/server"
)

func main() {
	rsa := server.NewRestServerAgent(":8000")
	rsa.Start()
	fmt.Scanln()
}
