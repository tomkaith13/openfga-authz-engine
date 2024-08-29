package utils

import (
	"context"
	"errors"
	"fmt"

	openfgaClient "github.com/openfga/go-sdk/client"
)

func CheckManages(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userId string) error {
	options := openfgaClient.ClientCheckOptions{
		AuthorizationModelId: &modelId,
	}

	body := openfgaClient.ClientCheckRequest{
		User:     "user:" + impersonatorId,
		Relation: "manages",
		Object:   "user:" + userId,
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
