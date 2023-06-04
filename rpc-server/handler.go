package main

import (
	"context"
	"fmt"
	"strings"
	"time"
	"log"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func validateSendRequest(req *rpc.SendRequest) error {
	senders := strings.Split(req.Message.Chat, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("Invalid Chat ID '%s', should be in the format of user1:user2", req.Message.GetChat())
		return err
	}
	sender1, sender2 := senders[0], senders[1]

	if req.Message.GetSender() != sender1 && req.Message.GetSender() != sender2 {
		err := fmt.Errorf("Sender '%s' is not in the chat room", req.Message.GetSender())
		return err
	}

	return nil
}

func getRoomID(chat string) (string, error) {
	var roomID string

	lowercaseString := strings.ToLower(chat)
	senders := strings.Split(lowercaseString, ":")
	if len(senders) != 2 {
		err := fmt.Errorf("Invalid Chat ID '%s', should be in the format of user1:user2", chat)
		return "", err
	}

	// Compare the sender and receiver alphabetically, and sort it ascending to form the roomID
	sender1, sender2 := senders[0], senders[1]
	comparison := strings.Compare(sender1, sender2); 
	if comparison == 1 {
		roomID = fmt.Sprintf("%s:%s", sender2, sender1)
	} else {
		roomID = fmt.Sprintf("%s:%s", sender1, sender2)
	}

	return roomID, nil
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	err := validateSendRequest(req); 
	if err != nil {
		return nil, err
	}
  
	// Refer to line 14 & 36 of rpc-server\redis.go
	// The *Message in the SaveMessage() method refers to the Message struct defined in rpc-server\redis.go
	timestamp := time.Now().Unix()
	message := &Message{
		Sender:    req.Message.GetSender(),
		Message:   req.Message.GetText(),
		Timestamp: timestamp,
	}
	log.Printf("Received message: %v", message)

	roomID, err := getRoomID(req.Message.GetChat())
	if err != nil {
		return nil, err
	}

	err = rdb.SaveMessage(ctx, roomID, message)
	if err != nil {
		return nil, err
	}

	resp := rpc.NewSendResponse()
	resp.Code, resp.Msg = 0, "success" // '0' means success in GoLang
	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	roomID, err := getRoomID(req.GetChat())
	if err != nil {
		return nil, err
	}

	limit := int64(req.GetLimit())
	if limit == 0 {
		limit = 20 // Default limit is 20
	}
	start := req.GetCursor()
	end := start + limit // 0 + 20 = 20 --> 21 items --> Did not minus 1 due to hasMore check later

	messages, err := rdb.GetMessagesByRoomID(ctx, roomID, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	// The Message struct used here is in idl_rpc.go (refer to line 3 - rpc package, and line 12 - Message struct)
	respMessages := make([]*rpc.Message, 0)
	var counter int64 = 0
	var nextCursor int64 = 0
	hasMore := false
	for _, msg := range messages {
		if counter + 1 > limit {
			// If counter + 1 > limit, it means that there are more messages to be pulled
			hasMore = true
			nextCursor = end
			break // Do not return the last message (i.e. the 21st message; index=20)
		}

		// Line 12 of rpc-server/idl_rpc.go
		temp := &rpc.Message{
			Chat:     req.GetChat(),
			Text:     msg.Message,
			Sender:   msg.Sender,
			SendTime: msg.Timestamp,
		}
		respMessages = append(respMessages, temp)
		counter += 1
	}

	// Line 1150 of rpc-server/idl_rpc.go
	// The * in line 1154 & 1155 means that we will deference the pointers here to get the value later
	resp := rpc.NewPullResponse()
	resp.Code = 0
	resp.Msg = "success"
	resp.Messages = respMessages
	resp.HasMore = &hasMore
	resp.NextCursor = &nextCursor

	return resp, nil
}
