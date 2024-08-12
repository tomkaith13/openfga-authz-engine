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
			{ // beth is in claims_reader
				User:     "user:beth",
				Relation: "member",
				Object:   "group:claims_reader",
			},
			{
				User:     "group:claims_reader#member",
				Relation: "reader",
				Object:   "capability:claims",
			},
			{ // jerry is in wallet deleter
				User:     "user:jerry",
				Relation: "member",
				Object:   "group:wallet_deleter",
			},
			{
				User:     "group:wallet_deleter#member",
				Relation: "deleter",
				Object:   "capability:wallet",
			},
			{ // jerry is in wallet reader
				User:     "user:jerry",
				Relation: "member",
				Object:   "group:wallet_reader",
			},
			{
				User:     "group:wallet_reader#member",
				Relation: "reader",
				Object:   "capability:wallet",
			},
			{ // morty is in claims updater
				User:     "user:morty",
				Relation: "member",
				Object:   "group:claims_updater",
			},
			{
				User:     "group:claims_updater#member",
				Relation: "updater",
				Object:   "capability:claims",
			},
			{ // morty is in claims reader
				User:     "user:morty",
				Relation: "member",
				Object:   "group:claims_reader",
			},
			{ // rick is in journey reader
				User:     "user:rick",
				Relation: "member",
				Object:   "group:journey_reader",
			},
			{
				User:     "group:journey_reader#member",
				Relation: "reader",
				Object:   "capability:journey",
			},
			{ // summer is in journey creator
				User:     "user:summer",
				Relation: "member",
				Object:   "group:journey_creator",
			},
			{
				User:     "group:journey_creator#member",
				Relation: "creator",
				Object:   "capability:journey",
			},
			{ // birdman is in claims admin
				User:     "user:birdman",
				Relation: "member",
				Object:   "group:claims_admin",
			},
			{
				User:     "group:claims_admin#member",
				Relation: "admin",
				Object:   "capability:claims",
			},
			{ // squanchy is in claims admin
				User:     "user:squanchy",
				Relation: "member",
				Object:   "group:journey_admin",
			},
			{
				User:     "group:journey_admin#member",
				Relation: "admin",
				Object:   "capability:journey",
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

func LoadAssertions(fgaClient *openfgaClient.OpenFgaClient, modelId string) error {
	options := openfgaClient.ClientWriteAssertionsOptions{
		AuthorizationModelId: &modelId,
	}
	requestBody := openfgaClient.ClientWriteAssertionsRequest{
		openfgaClient.ClientAssertion{
			User:        "user:birdman",
			Relation:    "can_delete",
			Object:      "capability:claims",
			Expectation: true,
		},
		openfgaClient.ClientAssertion{
			User:        "user:birdman",
			Relation:    "can_all",
			Object:      "capability:claims",
			Expectation: true,
		},
		openfgaClient.ClientAssertion{
			User:        "user:birdman",
			Relation:    "can_delete",
			Object:      "capability:journey",
			Expectation: false,
		},
		openfgaClient.ClientAssertion{
			User:        "user:squanchy",
			Relation:    "can_delete",
			Object:      "capability:journey",
			Expectation: true,
		},
		openfgaClient.ClientAssertion{
			User:        "user:beth",
			Relation:    "can_read",
			Object:      "capability:claims",
			Expectation: true,
		},
		openfgaClient.ClientAssertion{
			User:        "user:beth",
			Relation:    "can_create",
			Object:      "capability:claims",
			Expectation: false,
		},
		openfgaClient.ClientAssertion{
			User:        "user:morty",
			Relation:    "can_read",
			Object:      "capability:claims",
			Expectation: true,
		},
		openfgaClient.ClientAssertion{
			User:        "user:morty",
			Relation:    "can_update",
			Object:      "capability:claims",
			Expectation: true,
		},
		openfgaClient.ClientAssertion{
			User:        "user:beth",
			Relation:    "can_delete",
			Object:      "capability:claims",
			Expectation: false,
		},
	}
	_, err := fgaClient.WriteAssertions(context.Background()).
		Body(requestBody).
		Options(options).
		Execute()

	if err != nil {
		return err
	}
	return nil
}
