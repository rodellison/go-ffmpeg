package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CustomResponseOutput struct {
	Message string `json:"message"`
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
//type Response events.APIGatewayProxyResponse

type Response CustomResponseOutput

// Handler is our lambda handler invoked by the `lambda.Start` function call
// Uncomment the version needed based on what service is sending data to this lambda function..
// See: https://github.com/aws/aws-lambda-go/tree/master/events

//Lambda function gets called as a result of an APIGateway request
//func Handler(ctx context.Context, inEventData events.APIGatewayProxyRequest) (Response, error) {

//Lambda function gets called as a result of an S3 Event (PutObject completed, etc.)
func Handler(ctx context.Context, event events.S3Event) (Response, error) {

	//Do something here!
	fmt.Printf("Received: S3 Event = %v", event)
	bucket := event.Records[0].S3.Bucket.Name
	key := event.Records[0].S3.Object.Key

	fmt.Printf("A new thumbnail file was generated at 'https://s3.amazonaws.com/%s/%s'", bucket, key)

	return responseHandler(true, nil)
}

func responseHandler(success bool, errToRespond error) (Response, error) {

	//var buf bytes.Buffer
	//body, err := json.Marshal(map[string]interface{}{
	//	"message": "Okay so your other function also executed successfully!",
	//})
	//if err != nil {
	//	return Response{StatusCode: 404}, err
	//}
	//json.HTMLEscape(&buf, body)
	//
	//resp := Response{
	//	StatusCode:      200,
	//	IsBase64Encoded: false,
	//	Body:            buf.String(),
	//	Headers: map[string]string{
	//		"Content-Type":           "application/json",
	//		"X-MyCompany-Func-Reply": "world-handler",
	//	},
	//}

	if !success {
		return Response{
			Message: "NOT Successful!",
		}, errToRespond
	} else
	{
		return Response{
			Message: "Successful!",
		}, errToRespond
	}

}

func main() {
	lambda.Start(Handler)
}
