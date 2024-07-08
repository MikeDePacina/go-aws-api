package database

import (
	"fmt"
	"lambda/func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	USERS_TABLE = "Users"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	RegisterUser(user *types.User) error
	GetUser(username string) (types.User, error)
}

type DynamoDBClient struct{
	databaseStore *dynamodb.DynamoDB
}


func (d DynamoDBClient) DoesUserExist(username string) (bool, error){
	res, err := d.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USERS_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return true, err
	}

	if res.Item == nil {
		return false, nil
	}

	return true, nil
}

func (d DynamoDBClient) RegisterUser(user *types.User) error{
	item := &dynamodb.PutItemInput{
		TableName: aws.String(USERS_TABLE),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}
	
	_, err := d.databaseStore.PutItem(item)

	if err != nil {
		return err
	}

	return nil
}

func (d DynamoDBClient) GetUser(username string) (types.User, error){
	var user types.User
	res, err := d.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USERS_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		return user, err
	}

	if res.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(res.Item, &user)

	if err != nil {
		return user, err
	}

	return user, nil

}

func NewDynamoDBClient() DynamoDBClient{
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)
	
	return DynamoDBClient{
		databaseStore: db,
	}
}