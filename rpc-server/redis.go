package main

import (
	"context"
	"encoding/json"
	"log"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func (c *RedisClient) InitializeClient(ctx context.Context, address, password string) error {
	r := redis.NewClient(&redis.Options{
		Addr:     address,  // "redis:6379", 
		Password: password, // No password (empty string "")
		DB:       0,        // Use default DB --> Reference: https://redis.io/docs/clients/go/
	})

	// Test connection by pinging the Redis server to see if there's errors
	if err := r.Ping(ctx).Err(); err != nil {
		return err
	}

	c.client = r
	return nil
}

func (c *RedisClient) SaveMessage(ctx context.Context, roomID string, message *Message) error {
	// Store the message in JSON, if an error occurs, end the fn call & return the error
	// Reference: https://redis.io/docs/data-types/sorted-sets/
	text, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Using anonymous struct to store the message in the sorted set
	// Read up more on sorted sets here: https://www.youtube.com/watch?v=MUKlxdBQZ7g 
	// Sorted sets are sets that contain objects that are sorted by a key/score (in this case, the timestamp)
	// RoomID is used to identify which conversation the message belongs to (which room that the convo is in)
	member := &redis.Z{
		Score:  float64(message.Timestamp),
		Member: text,
	}
	log.Println("What is the member? ", member)

	/* 
		To summarize:
			roomID is the key or identifier for the Redis sorted set (identifies which conversation the message belongs to)
			The Score field in *redis.Z represents the score, which is the timestamp of the message. (The sorting key)
			The Member field in *redis.Z represents the message text or the member field in Redis. (The value)
			Once the member is constructed, it can be added to the sorted set in Redis using the ZAdd method.
			roomID is the key, *member will give us Score (which is the score) & Member (which is the message text/member field)
			Reference: 
				https://redis.io/commands/zadd/ (key, score, member)
				https://redis.io/commands/zadd/
				https://github.com/redis/go-redis/blob/master/commands_test.go (cltrl F 'redis.Z')
	*/
	_, err = c.client.ZAdd(ctx, roomID, *member).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisClient) GetMessagesByRoomID(ctx context.Context, roomID string, start, end int64, reverse bool) ([]*Message, error) {
	var (
		rawMessages []string
		messages    []*Message
		err         error
	)

	if reverse {
		// Descending order with time -> first message is the latest message
		// https://redis.io/commands/zrevrange/
		rawMessages, err = c.client.ZRevRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	} else {
		// Ascending order with time -> first message is the earliest message
		// https://redis.io/commands/zrange/
		rawMessages, err = c.client.ZRange(ctx, roomID, start, end).Result()
		if err != nil {
			return nil, err
		}
	}

	// Reference:
	// https://gobyexample.com/json
	// https://pkg.go.dev/encoding/json#Unmarshal
	for _, msg := range rawMessages {
		temp := &Message{}
		err := json.Unmarshal([]byte(msg), temp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, temp)
	}

	return messages, nil
}