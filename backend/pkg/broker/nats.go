package broker

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/onlinejudge/backend/internal/models"
	"github.com/onlinejudge/backend/internal/services"
)

const (
	SubmissionSubject = "submission.evaluate"
)

type NATSClient struct {
	conn *nats.Conn
}

func NewNATSClient() (*NATSClient, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NATSClient{conn: nc}, nil
}

func (c *NATSClient) PublishSubmission(submission *models.Submission) error {
	data, err := json.Marshal(submission)
	if err != nil {
		return err
	}

	return c.conn.Publish(SubmissionSubject, data)
}

func (c *NATSClient) SubscribeToSubmissions(evaluator *services.Evaluator) error {
	_, err := c.conn.Subscribe(SubmissionSubject, func(msg *nats.Msg) {
		var submission models.Submission
		if err := json.Unmarshal(msg.Data, &submission); err != nil {
			log.Printf("Error unmarshaling submission: %v", err)
			return
		}

		if err := evaluator.Evaluate(&submission); err != nil {
			log.Printf("Error evaluating submission: %v", err)
		}
	})

	return err
}

func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
} 