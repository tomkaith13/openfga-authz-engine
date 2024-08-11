package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	openfga "github.com/openfga/go-sdk"
	openfgaClient "github.com/openfga/go-sdk/client"
)

func main() {
	r := chi.NewRouter()

	// Middleware for logging and recovery
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	_, err := openfga.NewConfiguration(openfga.Configuration{
		ApiUrl: os.Getenv("FGA_API_URL"), // required, e.g. https://api.fga.example
	})

	if err != nil {
		fmt.Println("error::", err)
	}
	fgaClient, err := openfgaClient.NewSdkClient(&openfgaClient.ClientConfiguration{
		ApiUrl: os.Getenv("FGA_API_URL"), // required, e.g. https://api.fga.example
	})

	if err != nil {
		// .. Handle error
	}

	resp, err := fgaClient.CreateStore(context.Background()).
		Body(openfgaClient.ClientCreateStoreRequest{Name: "FGA Authz Engine"}).Execute()
	if err != nil {
		// .. Handle error
	}

	fmt.Println("store details:", resp.Name, resp.Id)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("World!!"))
	})
	http.ListenAndServe(":8888", r)
}
