package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	openfgaClient "github.com/openfga/go-sdk/client"
)

var LRUCache = expirable.NewLRU[string, bool](1, nil, time.Second*600)

const (
	externalResolverKey string = "ext"
)

func CheckImpersonator(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userId string) error {
	options := openfgaClient.ClientCheckOptions{
		AuthorizationModelId: &modelId,
	}

	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02T15:04:05Z")
	fmt.Println("formated UTC time:", formattedTime)

	body := openfgaClient.ClientCheckRequest{
		User:     "user:" + userId,
		Relation: "impersonator",
		Object:   "user:" + impersonatorId,
		Context:  &map[string]interface{}{"current_time": formattedTime},
	}

	data, err := fgaClient.Check(context.Background()).
		Body(body).
		Options(options).
		Execute()

	if err != nil {
		return err
	}

	if data.CheckResponse.Allowed != nil && !*data.CheckResponse.Allowed {
		return errors.New("Not Allowed Error")
	}
	fmt.Println("check results: allowed = ", *data.CheckResponse.Allowed)
	return nil
}

// CheckImpersonatorWithExternalResolver passes the current time in Context and also reaches out to external resolvers.
// This is the equivalent of http_send in OPA
func CheckImpersonatorWithExternalResolver(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userId string) error {
	options := openfgaClient.ClientCheckOptions{
		AuthorizationModelId: &modelId,
	}

	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02T15:04:05Z")
	fmt.Println("formated UTC time:", formattedTime)

	// Enable caching of external resolvers
	// we cant cache these using context since context need to accept real time info which is time.
	extValue, ok := LRUCache.Get(externalResolverKey)
	if !ok {
		LRUCache.Add(externalResolverKey, externalResolver())
	}
	extValue, _ = LRUCache.Get(externalResolverKey)

	body := openfgaClient.ClientCheckRequest{
		User:     "user:" + userId,
		Relation: "impersonator",
		Object:   "user:" + impersonatorId,
		Context:  &map[string]interface{}{"current_time": formattedTime, "external_check": extValue}, // we are passing the results from externalResolver
	}

	data, err := fgaClient.Check(context.Background()).
		Body(body).
		Options(options).
		Execute()

	if err != nil {
		return err
	}

	if data.CheckResponse.Allowed != nil && !*data.CheckResponse.Allowed {
		return errors.New("Not Allowed Error")
	}
	fmt.Println("check results: allowed = ", *data.CheckResponse.Allowed)
	return nil
}

func externalResolver() bool {
	// Simulate work
	time.Sleep(5 * time.Second)
	return true
}
