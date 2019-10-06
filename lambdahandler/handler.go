package lambdahandler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/kr/pretty"
)

type Response struct {
	Name    string
	Address string
}

//HandleRequest handles incoming request
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, _ = pretty.Println("parsed:", request.Body)
	return events.APIGatewayProxyResponse{Body: "response is working", StatusCode: 200}, nil
}