package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	initdata "github.com/telegram-mini-apps/init-data-golang"
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
	assert.Nil(t, response.Body)                 // Body should be nil for errors
	assert.Equal(t, "Not Found", response.Error) // Check the error message
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

	// Assert that the response body is a map and contains the expected fields
	body, ok := response.Body.(initdata.InitData)
	assert.True(t, ok, "Response body should be a map")

	// Check specific fields in the parsed data
	assert.Equal(t, int64(279058397), body.User.ID) // User ID
	assert.Equal(t, "Vladislav + - ? /", body.User.FirstName)
	assert.Equal(t, "ru", body.User.LanguageCode)
	assert.Equal(t, true, body.User.IsPremium)
}

func TestHandler_ValidateInitData_InvalidInitData(t *testing.T) {
	invalidInitData := "invalid_data" // Simulate invalid initData

	request := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Path: "/api/validateinitdata",
			},
		},
		QueryStringParameters: map[string]string{
			"initData": invalidInitData,
			"botID":    "7342037359",
		},
	}

	response, err := handler(request)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)

	// Ensure the Body is nil for errors
	assert.Nil(t, response.Body)

	// Check the exact error message returned
	assert.Equal(t, "Validation failed: sign is missing", response.Error)
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
	assert.Nil(t, response.Body)
	assert.Contains(t, response.Error, "Invalid botID: strconv.ParseInt: parsing \"invalid_bot_id\": invalid syntax")
}

func TestHandler_ValidateInitData_WithInvalidExpiration(t *testing.T) {
	validInitData := "user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%7D&auth_date=1662771648&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2"

	// Test with invalid expiration (negative value)
	request := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Path: "/api/validateinitdata",
			},
		},
		QueryStringParameters: map[string]string{
			"initData": validInitData,
			"botID":    "7342037359",
			"expIn":    "-3600", // Invalid expiration (negative value)
		},
	}

	response, err := handler(request)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Nil(t, response.Body) // Body should remain nil for errors
	assert.Equal(t, "Invalid expIn: must be a positive integer representing seconds", response.Error)
}

func TestHandler_ValidateInitData_WithMissingExpiration(t *testing.T) {
	validInitData := "user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%7D&auth_date=1662771648&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2"

	// Test with missing expiration (should default to 24 hours)
	request := events.LambdaFunctionURLRequest{
		RequestContext: events.LambdaFunctionURLRequestContext{
			HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
				Path: "/api/validateinitdata",
			},
		},
		QueryStringParameters: map[string]string{
			"initData": validInitData,
			"botID":    "7342037359",
		},
	}

	response, err := handler(request)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)               // Validation will fail due to invalid initData
	assert.Nil(t, response.Body)                            // Body should remain nil for errors
	assert.Contains(t, response.Error, "Validation failed") // Generic validation failure
}
