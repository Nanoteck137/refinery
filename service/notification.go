package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nanoteck137/refinery/config"
)

type Priority int

const (
	PriorityMin    Priority = 1
	PriorityLow    Priority = 2
	PriorityNormal Priority = 3
	PriorityHigh   Priority = 4
	PriorityUrgent Priority = 5
)

func (p Priority) String() string {
	switch p {
	case PriorityMin:
		return "min"
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityUrgent:
		return "urgent"
	default:
		return "normal"
	}
}

const (
	NOTIFY_TAG_WARNING = "warning"
)

type NotificationService struct {
	logger *slog.Logger

	BaseUrl string
	Topic   string
}

func NewNotificationService(logger *slog.Logger, config *config.Config) *NotificationService {
	return &NotificationService{
		logger:  logger,
		BaseUrl: "https://ntfy.nanoteck137.net",
		Topic:   "test",
	}
}

type notificationMessageBody struct {
	Topic    string   `json:"topic"`
	Message  string   `json:"message,omitempty"`
	Title    string   `json:"title"`
	Tags     []string `json:"tags,omitempty"`
	Priority int      `json:"priority,omitempty"`

	Markdown bool `json:"markdown"`
}

type SimpleNotificationOptions struct {
	Tags     []string
	Priority Priority
}

func (s *NotificationService) sendNotification(body notificationMessageBody) error {
	d, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest("POST", s.BaseUrl, bytes.NewReader(d))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("returned non 200 status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *NotificationService) SendSimple(title, message string, opts SimpleNotificationOptions) error {
	body := notificationMessageBody{
		Topic:    s.Topic,
		Message:  message,
		Title:    title,
		Tags:     opts.Tags,
		Priority: int(opts.Priority),
		Markdown: true,
	}

	err := s.sendNotification(body)
	if err != nil {
		s.logger.Error("failed to send simple notification",
			slog.String("title", title),
			slog.String("message", message),
			slog.String("priority", opts.Priority.String()),
			slog.Any("tags", opts.Tags),
			slog.Any("err", err),
		)

		return fmt.Errorf("failed to send notification: %w", err)
	}

	s.logger.Info("sent simple notification",
		slog.String("title", title),
		slog.String("message", message),
		slog.String("priority", opts.Priority.String()),
		slog.Any("tags", opts.Tags),
	)

	return nil
}
