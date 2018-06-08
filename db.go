package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("ap-southeast-2"))

func createFeedbackItem(feedback *feedback, sessionID string) error {
	newfeedback, _ := dynamodbattribute.Marshal(feedback)

	input := &dynamodb.UpdateItemInput{

		TableName: aws.String("sidekick-dev-dfezer"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(sessionID),
			},
		},
		// UpdateExpression: aws.String("set Title = :title, Content = :content, FeedbackID = :feedbackid"),
		UpdateExpression: aws.String("SET feedbacks = list_append(if_not_exists(feedbacks, :empty_list), :feedback)"),

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":feedback": {L: []*dynamodb.AttributeValue{newfeedback}},

			":empty_list": {L: []*dynamodb.AttributeValue{}},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err := db.UpdateItem(input)
	fmt.Print(input)
	return err

}
