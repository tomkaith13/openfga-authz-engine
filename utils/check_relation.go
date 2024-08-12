package utils

import (
	"context"
	"errors"
	"fmt"

	openfgaClient "github.com/openfga/go-sdk/client"
)

func Check(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, relation string, capabilityId string) error {
	options := openfgaClient.ClientCheckOptions{
		AuthorizationModelId: &modelId,
	}
	body := openfgaClient.ClientCheckRequest{
		User:     "user:" + impersonatorId,
		Relation: relation,
		Object:   "capability:" + capabilityId,
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
	fmt.Println("relation check results: allowed = ", *data.CheckResponse.Allowed)
	return nil
}
