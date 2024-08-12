package utils

import (
	"context"
	"fmt"

	openfgaClient "github.com/openfga/go-sdk/client"
)

func TupleLoader(fgaClient *openfgaClient.OpenFgaClient, modelId string) error {
	options := openfgaClient.ClientWriteOptions{
		AuthorizationModelId: &modelId,
	}
	body := openfgaClient.ClientWriteRequest{
		Writes: []openfgaClient.ClientTupleKey{
			{
				User:     "user:beth",
				Relation: "member",
				Object:   "group:claims_reader",
			},
			{
				User:     "group:claims_reader#member",
				Relation: "reader",
				Object:   "capability:claims",
			},
		},
	}
	data, err := fgaClient.Write(context.Background()).
		Body(body).
		Options(options).
		Execute()

	if err != nil {
		return err
	}
	fmt.Printf("writes status: %+v\n", data.Writes)

	return nil
}
