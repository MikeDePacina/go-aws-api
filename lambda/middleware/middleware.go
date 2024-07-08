package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error))(func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error)){
	return func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
		//get token from header
		token := extractTokenFromHeaders(request.Headers)

		if token == "" {
			return events.APIGatewayProxyResponse{
				Body: "Unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, nil
		
		}

		claims, err := parseToken(token)

		if err != nil {
			return events.APIGatewayProxyResponse{
				Body: "Unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body: "Unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		return next(request)
	}
}

func extractTokenFromHeaders(headers map[string]string) string{
	authHeader, ok := headers["Authorization"]

	if !ok {
		return ""
	}	

	token := strings.Split(authHeader, "Bearer ")

	if len(token) != 2 {
		return ""
	}

	return token[1]
	
}

func parseToken(tokenString string) (jwt.MapClaims, error){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		secret := "team secret"
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("token is invalid")
	}

	return claims, nil
}