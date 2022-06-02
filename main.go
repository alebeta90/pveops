package main

import (
	"fmt"
	"net/http"
	"os"

	"git.gonkar.com/gonkar/infra-cmd/routers"
)

func main() {

	routers.Main()

	port := os.Getenv("PORT") //Get port from .env file, we did not specify any port so this should return an empty string when tested locally
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Printf("Starting server on port: %s\n", port)

	handler := routers.C.Handler(routers.MainRouter)
	err := http.ListenAndServe(":"+port, handler) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
