package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type LambdaResponse struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func main() {
	lambda.Start(handler)
}

func handleParseInitData(initData string) LambdaResponse {
	parsedData, err := initdata.Parse(initData)
	if err != nil {
		return LambdaResponse{
			StatusCode: 400,
			Error:      fmt.Sprintf("Error parsing init data: %v", err),
		}
	}

	// Convert parsed data to JSON for the response
	if err != nil {
		return LambdaResponse{
			StatusCode: 500,
			Error:      fmt.Sprintf("Error serializing parsed data: %v", err),
		}
	}

	return LambdaResponse{
		StatusCode: 200,
		Body:       parsedData,
	}
}

func handleValidateInitData(initData string, botID int64, expIn time.Duration) LambdaResponse {
	err := initdata.ValidateThirdParty(initData, botID, expIn)
	if err != nil {
		return LambdaResponse{
			StatusCode: 400,
			Error:      fmt.Sprintf("Validation failed: %v", err),
		}
	}

	return LambdaResponse{
		StatusCode: 200,
		Body:       "Validation successful",
	}
}

func handler(request events.LambdaFunctionURLRequest) (LambdaResponse, error) {

	requestJSON, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to serialize request: %v", err)
	} else {
		log.Printf("Incoming request: %s", string(requestJSON))
	}

	path := request.RequestContext.HTTP.Path

	var response LambdaResponse

	switch path {
	case "/api/parseinitdata":
		initData := request.QueryStringParameters["initData"]
		response = handleParseInitData(initData)
	case "/api/validateinitdata":
		initData := request.QueryStringParameters["initData"]
		botIDStr := request.QueryStringParameters["botID"]
		expInStr := request.QueryStringParameters["expIn"] // Get expiration from query

		botID, err := strconv.ParseInt(botIDStr, 10, 64)
		if err != nil {
			return LambdaResponse{
				StatusCode: 400,
				Error:      fmt.Sprintf("Failed botID: %v", err),
			}, nil
		}
		expIn := 24 * time.Hour
		if expInStr != "" {
			expInSeconds, err := strconv.ParseInt(expInStr, 10, 64)
			if err != nil || expInSeconds <= 0 {
				return LambdaResponse{
					StatusCode: 400,
					Error:      "Invalid expIn: must be a positive integer representing seconds",
				}, nil
			}
			expIn = time.Duration(expInSeconds) * time.Second
		}

		response = handleValidateInitData(initData, botID, expIn)

		// Log the outgoing response
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to serialize response: %v", err)
		} else {
			log.Printf("Outgoing response: %s", string(responseJSON))
		}
	default:
		response = LambdaResponse{
			StatusCode: 404,
			Error:      "Not Found",
		}
	}

	return response, nil
}
