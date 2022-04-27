package main

import (
	"github.com/BrandonICR/web_cl2_050422_8am/cmd/server/engine"
)

const JSON_STORE_FILENAME = "./transacciones.json"

// @title Transaction Management API
// @version 1.0
// @description This API Handle Transactions
// @contact.name Transactions Team
// @contact.url https://someurl.com/support
// @licence.name Apache 2.0
// @licence.url https://somelicence.com/licences/LICENCE-2.0.html
func main() {

	router := engine.GetEngine(JSON_STORE_FILENAME, "", "")

	if err := router.Run(); err != nil {
		panic("error: no se logro ejecutar el servidor")
	}
}
