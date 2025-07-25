$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -tags lambda.norpc -o bootstrap
~\Go\Bin\build-lambda-zip.exe -o sms-relay-forwarder.zip bootstrap