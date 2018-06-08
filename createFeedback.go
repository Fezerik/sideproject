package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	uuid "github.com/satori/go.uuid"
)

var r = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

type feedback struct {
	FeedbackID string `json:"feedbackid"`
	Title      string `json:"title"`
	Content    string `json:"content"`
}

func createFeedback(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}
	sessionID := strings.TrimPrefix(req.Path, "/feedback/")
	// id := req.PathParameters["id"]
	// fmt.Print("this is my UUID? ", id)
	// fmt.Printf(id)
	// if !r.MatchString(id) {
	// 	return clientError(http.StatusBadRequest)
	// }

	// sid, err := getSessionItem(id)
	// if err != nil {
	// 	return serverError(err)
	// }
	// if sid == nil {
	// 	return clientError(http.StatusNotFound)
	// }

	fbid := uuid.NewV4().String()

	feedback := &feedback{
		// ID:         id,
		FeedbackID: fbid,
		Title:      "",
		Content:    "",
	}

	err := json.Unmarshal([]byte(req.Body), feedback)
	if err != nil {
		// return clientError(http.StatusUnprocessableEntity)
		log.Print(err)
	}

	if feedback.Title == "" || feedback.Content == "" {
		return clientError(http.StatusBadRequest)
	}

	err = createFeedbackItem(feedback, sessionID)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		// Headers:    map[string]string{"Location": fmt.Sprintf("/feedback/{id=%s}/review", fb.ID)},
	}, nil
}

func main() {
	lambda.Start(createFeedback)
}
