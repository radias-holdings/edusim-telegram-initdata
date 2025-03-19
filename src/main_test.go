package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler_NotFound(t *testing.T) {
	request := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Path: "/api/unknown",
			},
		},
		Headers: map[string]string{},
	}

	response, err := handler(request)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.StatusCode)
	assert.Equal(t, "Not Found", response.Body)
}

func TestHandler_ParseInitData(t *testing.T) {
	initData := "user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%20%2B%20-%20%3F%20%5C%2F%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg%22%7D&chat_instance=8134722200314281151&chat_type=private&auth_date=1733509682&signature=TYJxVcisqbWjtodPepiJ6ghziUL94-KNpG8Pau-X7oNNLNBM72APCpi_RKiUlBvcqo5L-LAxIc3dnTzcZX_PDg&hash=a433d8f9847bd6addcc563bff7cc82c89e97ea0d90c11fe5729cae6796a36d73"

	request := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Path: "/api/parseinitdata",
			},
		},
		QueryStringParameters: map[string]string{
			"initData": initData,
		},
	}

	response, err := handler(request)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, response.Body, `"id":279058397`)                   // Check if parsed data contains expected fields
	assert.Contains(t, response.Body, `"first_name":"Vladislav + - ? /"`) // Updated to match the actual response
}

func TestHandler_ValidateInitData_InvalidInitData(t *testing.T) {
	invalidInitData := "invalid_data"

	request := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Path: "/api/validateinitdata",
			},
		},
		QueryStringParameters: map[string]string{
			"initData": invalidInitData,
			"botID":    "7342037359",
			"secret":   "your-secret-key",
		},
	}

	response, err := handler(request)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Contains(t, response.Body, "Validation failed")
}

func TestHandler_ValidateInitData_InvalidBotID(t *testing.T) {
	validInitData := "user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%20%2B%20-%20%3F%20%5C%2F%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg%22%7D&chat_instance=8134722200314281151&chat_type=private&auth_date=1733584787&hash=2174df5b000556d044f3f020384e879c8efcab55ddea2ced4eb752e93e7080d6&signature=zL-ucjNyREiHDE8aihFwpfR9aggP2xiAo3NSpfe-p7IbCisNlDKlo7Kb6G4D0Ao2mBrSgEk4maLSdv6MLIlADQ"

	request := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Path: "/api/validateinitdata",
			},
		},
		QueryStringParameters: map[string]string{
			"initData": validInitData,
			"botID":    "invalid_bot_id",
			"secret":   "your-secret-key",
		},
	}

	response, err := handler(request)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Contains(t, response.Body, "Invalid botID: strconv.ParseInt: parsing \"invalid_bot_id\": invalid syntax")
}
