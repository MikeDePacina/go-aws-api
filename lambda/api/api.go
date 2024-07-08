package api

import (
	"encoding/json"
	"fmt"
	"lambda/func/database"
	"lambda/func/types"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct{
	dbStore database.UserStore
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)
	
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:"Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body: "Request has empty parameters",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("request has empty parameters")
	}

	userExists, err := api.dbStore.DoesUserExist(registerUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Error checking if user exists",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error checking if user exists: %v", err)
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body: "User already exists",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("user already exists")
	}
	
	user, err := types.NewUser(registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Error creating user",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error creating user: %v", err)
	}

	err = api.dbStore.RegisterUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Error inserting user in db",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user in db: %v", err)
	
	}

	return events.APIGatewayProxyResponse{
		Body: "User registered",
		StatusCode: http.StatusCreated,
	}, nil
}

func (api ApiHandler) LoginUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	var loginDTO types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &loginDTO)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if loginDTO.Username == "" || loginDTO.Password == "" {
		return events.APIGatewayProxyResponse{
			Body: "Request has empty parameters",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("request has empty parameters")
	}

	user, err := api.dbStore.GetUser(loginDTO.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Error getting user",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error getting user: %v", err)
	}

	if passwordMatch, _ := types.ValidatePassword(user.PasswordHash, loginDTO.Password); !passwordMatch {
		return events.APIGatewayProxyResponse{
			Body: "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
		}, fmt.Errorf("invalid credentials")

		
	}

	accessToken := types.CreateToken(user)
	successMsg := fmt.Sprintf("User logged in, token: %s", accessToken)
	

	return events.APIGatewayProxyResponse{
		Body: successMsg,
		StatusCode: http.StatusOK,
	}, nil

}

func NewApiHandler(dbStore database.UserStore) ApiHandler{
	return ApiHandler{
		dbStore: dbStore,
	}
}

