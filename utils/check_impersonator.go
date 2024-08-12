package utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	openfgaClient "github.com/openfga/go-sdk/client"
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
