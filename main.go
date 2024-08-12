package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	openfga "github.com/openfga/go-sdk"
	openfgaClient "github.com/openfga/go-sdk/client"
	"github.com/tomkaith13/openfga-authz-engine/utils"
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
		fmt.Println("Unable to create sdk client ", err)
		return
	}

	resp, err := fgaClient.CreateStore(context.Background()).
		Body(openfgaClient.ClientCreateStoreRequest{Name: os.Getenv("FGA_STORE_NAME")}).Execute()
	if err != nil {
		fmt.Println("Unable to create FGA Store:", err)
		return
	}

	fmt.Println("store details:", resp.Name, resp.Id)

	// We need to reinit the client to use the new StoreId we got from CreateStore
	fgaClient, err = openfgaClient.NewSdkClient(&openfgaClient.ClientConfiguration{
		ApiUrl:  os.Getenv("FGA_API_URL"), // required, e.g. https://api.fga.example
		StoreId: resp.Id,
		// AuthorizationModelId: "v1",
	})
	if err != nil {
		fmt.Println("Unable to create openfgaclient with new storeId", err)
		return
	}
	configFile := os.Getenv("FGA_CONFIG_FILE")
	bytes, err := utils.LoadConfig(configFile)
	if err != nil {
		fmt.Println("Unable to read config file", err)
		return
	}
	fmt.Println("config:", string(bytes))
	var body openfgaClient.ClientWriteAuthorizationModelRequest
	if err := json.Unmarshal(bytes, &body); err != nil {
		fmt.Println("Unable to load config to Authz Model Request", err)
		return
	}

	// setup config
	data, err := fgaClient.WriteAuthorizationModel(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Println("error writing authz model to openFGA server", err)
		return
	}

	fmt.Println("data model id:", data.AuthorizationModelId)

	// Now loading the data
	err = utils.TupleLoader(fgaClient, data.AuthorizationModelId)
	if err != nil {
		fmt.Println("Unable to load tuples. err:", err)
		return
	}

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("World!!"))
	})
	http.ListenAndServe(":8888", r)
}
