package services

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"backend/internal/models"
)

type NATSClient struct {
	conn *nats.Conn
}

func NewNATSClient(url string) (*NATSClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NATSClient{conn: nc}, nil
}

func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *NATSClient) PublishSubmission(submission *models.Submission) error {
	data, err := json.Marshal(submission)
	if err != nil {
		return err
	}

	return c.conn.Publish("submissions", data)
}

func (c *NATSClient) SubscribeToSubmissions(handler func(*models.Submission)) error {
	_, err := c.conn.Subscribe("submissions", func(msg *nats.Msg) {
		var submission models.Submission
		if err := json.Unmarshal(msg.Data, &submission); err != nil {
			log.Printf("Error unmarshaling submission: %v", err)
			return
		}
		handler(&submission)
	})

	return err
} 