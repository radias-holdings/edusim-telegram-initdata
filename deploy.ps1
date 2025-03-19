# Variables
$FunctionName = "dev-telegram-app-validator" # Replace with your Lambda function name
$ZipFile = "telegram-app-validation.zip"
$GoFile = "src/main.go" # Updated to include the relative path to main.go
$AwsProfile = "237893020496_MS-DEV" # AWS CLI profile for dev environment
$BuildDir = "build"
$Region = "ap-southeast-2"
$BuildLambdaZipPath = "$env:USERPROFILE\Go\bin\build-lambda-zip.exe" # Path to build-lambda-zip tool

# Step 1: Build the Go binary
Write-Host "Building Go binary..."
if (!(Test-Path $BuildDir)) {
    New-Item -ItemType Directory -Path $BuildDir | Out-Null
}
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o "$BuildDir/bootstrap" $GoFile
if ($LASTEXITCODE -ne 0) {
    Write-Host "Go build failed. Exiting."
    Exit 1
}

# Step 2: Create a ZIP file using build-lambda-zip
Write-Host "Creating ZIP file using build-lambda-zip..."
& $BuildLambdaZipPath -o "$BuildDir/$ZipFile" "$BuildDir/bootstrap"
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to create ZIP file. Exiting."
    Exit 1
}

# Step 3: Deploy to AWS Lambda
Write-Host "Deploying to AWS Lambda..."
aws lambda update-function-code --function-name $FunctionName --zip-file fileb://"$BuildDir/$ZipFile" --profile $AwsProfile --region $Region
if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to deploy to AWS Lambda. Exiting."
    Exit 1
}

# Cleanup
Write-Host "Cleaning up..."
Remove-Item -Recurse -Force $BuildDir

Write-Host "Deployment successful!"