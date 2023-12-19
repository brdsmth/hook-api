package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	localConfig "hook-api/config"
	"hook-api/services"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gonanoid "github.com/matoous/go-nanoid"
)

// Job represents the structure of the job data
type Job struct {
	ID        string      `json:"id"`
	Payload   interface{} `json:"payload"` // Can be any JSON data
	URL       string      `json:"url"`
	ExecuteAt string      `json:"executeAt"`
	Status    string      `json:"Status"`
}

func AddJob(w http.ResponseWriter, r *http.Request) {
	var job Job
	var err error
	// Identify job
	id, err := gonanoid.Nanoid(6)
	if err != nil {
		log.Printf("error adding nanoid: %s", err)
		return
	}
	job.ID = id
	log.Printf("Add job:\t\t\t%s", id)

	// Decode the job from the request body
	err = json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		log.Printf("Cancel add job:\t\t%s", id)
		return
	}

	// Check if URL is provided
	if job.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		log.Printf("Cancel add job:\t\t%s", id)
		return
	}

	// Check if Payload is provided
	if job.Payload == nil {
		http.Error(w, "Payload is required", http.StatusBadRequest)
		return
	}

	// Marshal the Payload to a JSON string
	payloadBytes, err := json.Marshal(job.Payload)
	if err != nil {
		http.Error(w, "Error marshaling payload", http.StatusInternalServerError)
		log.Printf("Cancel add job:\t\t%s", id)
		return
	}
	payloadString := string(payloadBytes)

	// Set ExecutionTime to now if not provided
	if job.ExecuteAt == "" {
		job.ExecuteAt = time.Now().Format(time.RFC3339)
	}

	// Set job status
	job.Status = "QUEUED"

	// Update DynamoDB
	partitionKey := fmt.Sprintf("queue::%s::%s", job.ID, job.ExecuteAt)

	// Prepare the item to write to DynamoDB
	item := map[string]types.AttributeValue{
		"RowKey":    &types.AttributeValueMemberS{Value: partitionKey},
		"JobID":     &types.AttributeValueMemberS{Value: job.ID},
		"Payload":   &types.AttributeValueMemberS{Value: payloadString}, // Serialized JSON string
		"URL":       &types.AttributeValueMemberS{Value: job.URL},
		"ExecuteAt": &types.AttributeValueMemberS{Value: job.ExecuteAt},
		"Status":    &types.AttributeValueMemberS{Value: job.Status},
	}

	// Write to DynamoDB
	// The table name to store processed message will be stored in an environment variable
	dynamoDBQueueTable := localConfig.ReadEnv("DYNAMODB_QUEUE_TABLE")
	if dynamoDBQueueTable == "" {
		log.Fatal("DYNAMODB_QUEUE_TABLE environment variable not set")
	}
	tableName := dynamoDBQueueTable
	_, err = services.DynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("Failed to write to DynamoDB: %v", err)
		http.Error(w, "Failed to add job", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job added successfully"))
}
