package kvstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBKVStore implements the KVStore interface using DynamoDB as a storage backend.
type DynamoDBKVStore struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBKVStore creates a new instance of DynamoDBKVStore.
func NewDynamoDBKVStore(tableName string) (*DynamoDBKVStore, error) {
	// Load the AWS configuration and initialize a new DynamoDB client.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}
	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBKVStore{
		client:    client,
		tableName: tableName,
	}, nil
}

// Set inserts or updates an item in the DynamoDB table.
func (s *DynamoDBKVStore) Set(key, value string) error {
	_, err := s.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &s.tableName,
		Item: map[string]types.AttributeValue{
			"Key":   &types.AttributeValueMemberS{Value: key},
			"Value": &types.AttributeValueMemberS{Value: value},
		},
	})
	return err
}

// Get retrieves an item from the DynamoDB table.
func (s *DynamoDBKVStore) Get(key string) (string, error) {
	result, err := s.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &s.tableName,
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return "", err
	}

	if result.Item == nil {
		return "", errors.New("key not found")
	}

	value := result.Item["Value"].(*types.AttributeValueMemberS).Value
	return value, nil
}

// Delete removes an item from the DynamoDB table.
func (s *DynamoDBKVStore) Delete(key string) error {
	_, err := s.client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: &s.tableName,
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: key},
		},
	})
	return err
}

// Exist checks whether an item exists in the DynamoDB table.
func (s *DynamoDBKVStore) Exist(key string) (bool, error) {
	result, err := s.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &s.tableName,
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return false, err
	}
	return result.Item != nil, nil
}

// Ping checks if the DynamoDB service is reachable by calling ListTables.
func (s *DynamoDBKVStore) Ping() error {
	_, err := s.client.ListTables(context.Background(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(1),
	})
	return err
}
