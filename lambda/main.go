package main

import (
	"fmt"
	"lambda/func/app"
	"lambda/func/middleware"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username is empty")
	}

	return fmt.Sprintf("sucessfully called by %s", event.Username), nil
}

func protectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	return events.APIGatewayProxyResponse{
		Body: "Protected",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path{
		case "/register":
			return myApp.ApiHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.ApiHandler.LoginUserHandler(request)
		case "/protected":
			return middleware.ValidateJWTMiddleware(protectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body: "Not Found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}