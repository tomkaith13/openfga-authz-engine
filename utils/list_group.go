package utils

import (
	"context"
	"fmt"

	openfgaClient "github.com/openfga/go-sdk/client"
)

func ListUserGroup(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string) error {
	options := openfgaClient.ClientListObjectsOptions{
		AuthorizationModelId: &modelId,
		// This consistency level is the default if you look at the source code
		// Consistency: openfga.CONSISTENCYPREFERENCE_MINIMIZE_LATENCY.Ptr(),
	}

	fmt.Println("entered list user group")

	body := openfgaClient.ClientListObjectsRequest{
		User:     "user:" + impersonatorId,
		Relation: "member",
		Type:     "group",
	}

	data, err := fgaClient.ListObjects(context.Background()).
		Body(body).
		Options(options).
		Execute()

	if err != nil {
		fmt.Println("error listing: %+v", err)
		return err
	}

	fmt.Printf("list objects results: %+v\n", data.GetObjects())
	return nil
}
