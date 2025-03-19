# edusim-telegram-initdata

This repository contains a Go-based AWS Lambda function for validating Telegram app data. Below are the instructions for running, testing, and deploying the application.

---

## Prerequisites

1. **Go Installed**: Ensure you have Go installed on your system. You can download it from [https://go.dev/dl/](https://go.dev/dl/).
2. **AWS CLI Installed**: Install the AWS CLI and configure it with your credentials. Refer to [AWS CLI Installation Guide](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html).
3. **build-lambda-zip Tool**: Download the `build-lambda-zip` tool for Windows. Follow the instructions in the AWS Lambda Go documentation: [AWS Lambda Go for Developers on Windows](https://github.com/aws/aws-lambda-go?tab=readme-ov-file#for-developers-on-windows).

---

## Running the Application Locally

To run the application locally:

1. Navigate to the `src` directory:
   ```bash
   cd src
   ```

2. Run the Go application:
   ```bash
   go run main.go
   ```

---

## Testing the Application

To test the application:

1. Write unit tests in the `src` directory (e.g., `main_test.go`).
2. Run the tests using the following command:
   ```bash
   go test ./...
   ```

---

## Deploying to AWS Lambda

To deploy the application to AWS Lambda:

1. Ensure the `deploy.ps1` script is configured with the correct variables:
   - `$FunctionName`: Your Lambda function name.
   - `$AwsProfile`: Your AWS CLI profile.
   - `$Region`: The AWS region where your Lambda function is deployed.

2. Run the deployment script:
   ```powershell
   .\deploy.ps1
   ```

3. The script will:
   - Build the Go binary for Linux.
   - Package the binary into a ZIP file using the `build-lambda-zip` tool.
   - Deploy the ZIP file to AWS Lambda.

---

## Additional Resources

For more details on developing AWS Lambda functions in Go on Windows, refer to the official AWS Lambda Go documentation:  
[https://github.com/aws/aws-lambda-go?tab=readme-ov-file#for-developers-on-windows](https://github.com/aws/aws-lambda-go?tab=readme-ov-file#for-developers-on-windows)