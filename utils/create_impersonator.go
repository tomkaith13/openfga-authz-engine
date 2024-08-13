package utils

import (
	"context"
	"fmt"
	"time"

	openfga "github.com/openfga/go-sdk"
	openfgaClient "github.com/openfga/go-sdk/client"
)

func CreateImpersonator(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userId string) error {
	options := openfgaClient.ClientWriteOptions{
		AuthorizationModelId: &modelId,
	}

	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02T15:04:05Z")
	fmt.Println("formated UTC time:", formattedTime)

	body := openfgaClient.ClientWriteRequest{
		Writes: []openfgaClient.ClientTupleKey{
			{
				User:     "user:" + userId,
				Relation: "impersonator",
				Object:   "user:" + impersonatorId,
				Condition: &openfga.RelationshipCondition{
					Name:    "check_expired",
					Context: &map[string]interface{}{"grant_time": formattedTime, "grant_duration": "1m"},
				},
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

	fmt.Printf("data from adding impersonator: %+v\n", data.Writes)
	return nil
}

func DeleteImpersonator(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userId string) error {
	options := openfgaClient.ClientWriteOptions{
		AuthorizationModelId: &modelId,
	}

	body := openfgaClient.ClientWriteRequest{
		Deletes: []openfgaClient.ClientTupleKeyWithoutCondition{
			{
				User:     "user:" + userId,
				Relation: "impersonator",
				Object:   "user:" + impersonatorId,
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

	fmt.Printf("data from deleting impersonator: %+v\n", data.Deletes)
	return nil
}

func GetImpersonator(fgaClient *openfgaClient.OpenFgaClient, storeId string, impersonatorId string, userId string) error {
	options := openfgaClient.ClientReadOptions{
		StoreId: &storeId,
	}

	userKey := "user:" + userId
	relation := "impersonator"
	objectKey := "user:" + impersonatorId
	body := openfgaClient.ClientReadRequest{
		User:     &userKey,
		Relation: &relation,
		Object:   &objectKey,
	}
	resp, err := fgaClient.Read(context.Background()).Options(options).Body(body).Execute()
	if err != nil {
		return err
	}

	tuples := resp.GetTuples()
	for _, t := range tuples {
		key := t.Key
		fmt.Println("user:", key.User, " relation:", key.Relation, " object:", key.Object)
		condition := key.Condition
		if condition != nil {
			context := condition.GetContext()
			fmt.Println("context:", context)
		}
	}
	return nil

}
