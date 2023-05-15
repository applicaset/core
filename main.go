package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nasermirzaei89/env"
	"net/http"
)

func main() {
	svc := NewStore()

	h := NewHandler(svc)

	apiAddress := env.GetString("API_ADDRESS", ":8080")

	if err := http.ListenAndServe(apiAddress, h); err != nil {
		panic(fmt.Errorf("error on listen and serve http"))
	}
}
