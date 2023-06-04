// go test handler_test.go -v (need -v flag to log more details about the tests e.g. your logs, output etc)

package main

import (
	"fmt"
	"context"
	"testing"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// MockRedisClient is a mock implementation of the Redis client.
type MockRedisClient struct{}

func (m *MockRedisClient) InitializeClient(ctx context.Context, address, password string) *redis.IntCmd {
	return nil
}

func (m *MockRedisClient) SaveMessage(ctx context.Context, roomID string, message *SendRequest) *redis.StringSliceCmd {
	return nil
}

type SendRequest struct{
	Chat string
	Text string
	Sender string
	SendTime int64
}

type SendResponse struct {
	Message string
	ErrorCode int
}

type Message struct {
	Chat string
	Text string
	Sender string
	SendTime int64
}

type PullRequest struct {
	Chat string
	Cursor int64 
	Limit int32
	Reverse bool
}

type PullResponse struct {
	Code int32      
	Msg string     
	Messages []*Message
	HasMore *bool
	NextCursor *int64
}

type IMServiceImpl struct {
	redisClient *MockRedisClient
}

func (s *IMServiceImpl) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	// Mock implementation of the Send method
	message := "Message saved!"
	errorCode := 0

	// Return the SendResponse instance with the message and error code
	return &SendResponse{
		Message:  message,
		ErrorCode: errorCode,
	}, nil
}

// Unit Test One
func TestIMServiceImplementation_Send(t *testing.T) {
	type args struct {
		ctx context.Context
		req *SendRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &SendRequest{
					Chat: `john:sam"`,
					Text: `Hello World"`,
					Sender:`john`,
					SendTime: time.Now().Unix(),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an instance of the mock Redis client
			mockRedisClient := &MockRedisClient{}

			// Create an instance of the IMServiceImpl, passing the mock Redis client
			s := &IMServiceImpl{
				redisClient: mockRedisClient,
			}

			// Call the Send method
			returnedValue, err := s.Send(tt.args.ctx, tt.args.req)

			// Assert that the method executed successfully without errors
			fmt.Println("What is error code:", err)
			fmt.Println("What is returned value: ", returnedValue)
			assert.NoError(t, err)
			assert.NotNil(t, returnedValue)
		})
	}
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *PullRequest) (*PullResponse, error) {
	// Mock implementation of the Pull method - mocking 5 messages in redis DB
	code := int32(200)
	msg := "success"
	messages := []*Message{
		{
			Chat: `john:sam`,
			Text: `Hello World`,
			Sender:`john`,
			SendTime: time.Now().Unix(),
		},
		{
			Chat: `sam:john`,
			Text: `Hi too`,
			Sender:`sam`,
			SendTime: time.Now().Unix(),
		},
		{
			Chat: `sam:john`,
			Text: `I'm new to GoLang`,
			Sender:`sam`,
			SendTime: time.Now().Unix(),
		},
		{
			Chat: `john:sam`,
			Text: `Me too! Let's learn together!`,
			Sender:`john`,
			SendTime: time.Now().Unix(),
		},
		{
			Chat: `sam:john`,
			Text: `For sure, hooray!`,
			Sender:`sam`,
			SendTime: time.Now().Unix(),
		},
	}
	hasMore := true
	nextCursor := int64(3)

	// Return the PullResponse instance with the message and error code
	return &PullResponse{
		Code: code,
		Msg: msg,
		Messages: messages,
		HasMore: &hasMore,
		NextCursor: &nextCursor,
	}, nil
}

// Unit Test Two
func TestIMServiceImplementation_Pull(t *testing.T) {
	// Define the expected values for the test - there should be 5 messages returned based on the mock implementation above
	expectedMessagesLength := 5

	// Create an instance of the mock Redis client
	mockRedisClient := &MockRedisClient{}

	// Create an instance of the IMServiceImpl, passing the mock Redis client
	s := &IMServiceImpl{
		redisClient: mockRedisClient,
	}

	// Define the payload for the Pull method
	payload := &PullRequest{
    Chat: "jack:tom",
    Cursor: 0,
    Limit: 10,
    Reverse: true,
	}

	// Call the Pull method
	messages, err := s.Pull(context.Background(), payload)

	// Assert that the method executed successfully without errors
	assert.NoError(t, err)

	// Assert that the returned messages' length (5) match the expected messages' length (5)
	fmt.Println("What is the messages array: ", messages.Messages)
	fmt.Println("What is the length of the messages array: ", len(messages.Messages))
	assert.Equal(t, expectedMessagesLength, len(messages.Messages))
}

