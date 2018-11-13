package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

type Request struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Err     bool   `json:"err"`
}

type NormalResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Item    string `json:"item"`
	Message string `json:"message"`
}

func (res ErrorResponse) Error() string {
	return res.Message
}

func main() {
	lambda.Start(Handle)
}

func Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("normal", "request-context:", ctx)
	log.Println("normal", "request-body:", req.Body)
	log.Println("normal", "request-header:", req.Headers)

	request := Request{}
	e := json.Unmarshal([]byte(req.Body), &request)
	if e != nil {
		log.Println("error", "parse request:", "body:", req.Body, "error:", e)
		err := ErrorResponse{Item: "json", Message: "invalid request"}
		bytes, _ := json.Marshal(err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: string(bytes)}, &err
	}

	code := request.Status
	switch code {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusAccepted:
	case http.StatusBadRequest:
	case http.StatusNotFound:
	case http.StatusConflict:
	case http.StatusInternalServerError:
		if request.Err {
			response := ErrorResponse{Item: "user", Message: request.Message}
			bytes, _ := json.Marshal(response)
			return events.APIGatewayProxyResponse{StatusCode: code, Body: string(bytes)}, &response
		}
		response := NormalResponse{Message: request.Message}
		bytes, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{StatusCode: code, Body: string(bytes)}, nil
	}

	response := ErrorResponse{Item: "status-code", Message: "not in list"}
	bytes, _ := json.Marshal(response)
	if request.Err {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound, Body: string(bytes)}, &response
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound, Body: string(bytes)}, nil
}
