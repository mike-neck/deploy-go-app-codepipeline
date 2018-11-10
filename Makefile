

.PHONY: deps clean build

dep:
	@echo dep
	@go get github.com/aws/aws-lambda-go/lambda

build:
	@echo build
	@GOOS=linux go build -o build/lambda-handler cmd/lambda-handler/main.go

clean:
	@echo clean
	@go clean
	@rm -rf build/
