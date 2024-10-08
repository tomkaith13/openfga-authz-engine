package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	openfga "github.com/openfga/go-sdk"
	openfgaClient "github.com/openfga/go-sdk/client"
)

var ErrorAlreadyExists error = errors.New("Tuple already exists")

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

func CreateImpersonatorWithExt(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userId string) error {
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
					Name:    "check_expired_with_ext",
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

// DeleteAndAddImpersonator results in an ERROR if the tuple already exists. It does not work!!
// Just added it here as a proof! Do not use!
func DeleteAndAddImpersonator(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userId string) error {
	options := openfgaClient.ClientWriteOptions{
		AuthorizationModelId: &modelId,
	}

	currentTime := time.Now().UTC()
	formattedTime := currentTime.Format("2006-01-02T15:04:05Z")
	fmt.Println("formated UTC time:", formattedTime)

	body := openfgaClient.ClientWriteRequest{
		Deletes: []openfgaClient.ClientTupleKeyWithoutCondition{
			{
				User:     "user:" + userId,
				Relation: "impersonator",
				Object:   "user:" + impersonatorId,
			},
		},
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

func CreateUserManagementTuples(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userIds []string) error {
	options := openfgaClient.ClientWriteOptions{
		AuthorizationModelId: &modelId,
	}
	writeTuples := []openfgaClient.ClientTupleKey{}
	for _, userId := range userIds {

		writeTuples = append(writeTuples, openfga.TupleKey{
			User:     "user:" + impersonatorId,
			Relation: "manages",
			Object:   "user:" + userId,
		})
	}
	body := openfgaClient.ClientWriteRequest{
		Writes: writeTuples,
	}
	data, err := fgaClient.Write(context.Background()).
		Body(body).
		Options(options).
		Execute()

	if err != nil {
		if strings.Contains(err.Error(), "tuple which already exists") {
			return ErrorAlreadyExists
		}
	}

	fmt.Printf("data from adding impersonator: %+v\n", data.Writes)
	return nil

}

// func DeleteUserManagementTuple(fgaClient *openfgaClient.OpenFgaClient, modelId string, impersonatorId string, userIds []string) error {
// 	options := openfgaClient.ClientWriteOptions{
// 		AuthorizationModelId: &modelId,
// 	}
// 	writeTuples := []openfgaClient.ClientTupleKey{}
// 	for _, userId := range userIds {

// 		writeTuples = append(writeTuples, openfga.TupleKey{
// 			User:     "user:" + impersonatorId,
// 			Relation: "manages",
// 			Object:   "user:" + userId,
// 		})
// 	}
// 	body := openfgaClient.ClientWriteRequest{
// 		Writes: writeTuples,
// 	}
// 	data, err := fgaClient.Write(context.Background()).
// 		Body(body).
// 		Options(options).
// 		Execute()

// 	if err != nil {
// 		if strings.Contains(err.Error(), "tuple which already exists") {
// 			return ErrorAlreadyExists
// 		}
// 	}

// 	fmt.Printf("data from adding impersonator: %+v\n", data.Writes)
// 	return nil

// }
