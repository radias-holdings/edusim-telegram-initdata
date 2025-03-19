package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func main() {
	lambda.Start(handler)
}

func handleParseInitData(initData string) events.LambdaFunctionURLResponse {
	parsedData, err := initdata.Parse(initData)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error parsing init data: %v", err),
		}
	}

	// Convert parsed data to JSON for the response
	parsedDataJSON, err := json.Marshal(parsedData)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error serializing parsed data: %v", err),
		}
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(parsedDataJSON),
	}
}

func handleValidateInitData(initData string, botID int64, expIn time.Duration) events.LambdaFunctionURLResponse {
	err := initdata.ValidateThirdParty(initData, botID, expIn)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Validation failed: %v", err),
		}
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "Validation successful",
	}
}

func handler(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	path := request.RequestContext.HTTP.Path

	var response events.LambdaFunctionURLResponse

	switch path {
	case "/api/parseinitdata":
		initData := request.QueryStringParameters["initData"]
		response = handleParseInitData(initData)
	case "/api/validateinitdata":
		initData := request.QueryStringParameters["initData"]
		botIDStr := request.QueryStringParameters["botID"]
		botID, err := strconv.ParseInt(botIDStr, 10, 64)
		if err != nil {
			return events.LambdaFunctionURLResponse{StatusCode: 400, Body: fmt.Sprintf("Invalid botID: %v", err)}, nil
		}
		expIn := 24 * time.Hour // Expiration time
		response = handleValidateInitData(initData, botID, expIn)

	default:
		response = events.LambdaFunctionURLResponse{StatusCode: 404, Body: "Not Found"}
	}

	return response, nil
}
