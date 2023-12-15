package main

import (
	"fmt"

	"github.com/vvthai10/transaction-mongodb/bootstrap"
)

func main() {
	app := bootstrap.App()
	fmt.Println(app.Env.DBUri)
}