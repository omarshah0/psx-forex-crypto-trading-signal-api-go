package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/models"
)

// NotificationSender is the interface for sending notifications
type NotificationSender interface {
	SendSignalNotification(signal *models.TradingSignal) error
}

// NotificationService manages multiple notification senders
type NotificationService struct {
	senders []NotificationSender
}

// NewNotificationService creates a new notification service with configured senders
func NewNotificationService(
	telegramEnabled bool, telegramBotToken, telegramChatID string,
	discordEnabled bool, discordWebhookURL string,
	expoEnabled bool,
) *NotificationService {
	var senders []NotificationSender

	if telegramEnabled && telegramBotToken != "" && telegramChatID != "" {
		senders = append(senders, NewTelegramNotificationService(telegramBotToken, telegramChatID))
	}

	if discordEnabled && discordWebhookURL != "" {
		senders = append(senders, NewDiscordNotificationService(discordWebhookURL))
	}

	if expoEnabled {
		senders = append(senders, NewExpoNotificationService())
	}

	return &NotificationService{
		senders: senders,
	}
}

// SendSignalNotification sends a signal notification to all configured senders
func (s *NotificationService) SendSignalNotification(signal *models.TradingSignal) error {
	if len(s.senders) == 0 {
		log.Println("No notification senders configured, skipping notification")
		return nil
	}

	var lastError error
	for _, sender := range s.senders {
		if err := sender.SendSignalNotification(signal); err != nil {
			log.Printf("Failed to send notification: %v", err)
			lastError = err
		}
	}

	return lastError
}

// TelegramNotificationService sends notifications to Telegram
type TelegramNotificationService struct {
	botToken   string
	chatID     string
	httpClient *http.Client
}

func NewTelegramNotificationService(botToken, chatID string) *TelegramNotificationService {
	return &TelegramNotificationService{
		botToken:   botToken,
		chatID:     chatID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *TelegramNotificationService) SendSignalNotification(signal *models.TradingSignal) error {
	message := fmt.Sprintf(
		"🚨 *New %s %s Signal!*\n\n"+
			"Asset: %s\n"+
			"Type: %s\n\n"+
			"Check the app for details! 📊",
		signal.AssetClass,
		signal.DurationType,
		signal.Symbol,
		signal.Type,
	)

	reqBody := map[string]interface{}{
		"chat_id":    s.chatID,
		"text":       message,
		"parse_mode": "Markdown",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram request: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create telegram request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send telegram request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	log.Printf("Telegram notification sent for signal ID %d", signal.ID)
	return nil
}

// DiscordNotificationService sends notifications to Discord
type DiscordNotificationService struct {
	webhookURL string
	httpClient *http.Client
}

func NewDiscordNotificationService(webhookURL string) *DiscordNotificationService {
	return &DiscordNotificationService{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *DiscordNotificationService) SendSignalNotification(signal *models.TradingSignal) error {
	embed := map[string]interface{}{
		"title":       fmt.Sprintf("🚨 New %s %s Signal!", signal.AssetClass, signal.DurationType),
		"description": fmt.Sprintf("A new trading signal has been posted for **%s**", signal.Symbol),
		"color":       0x00ff00, // Green color
		"fields": []map[string]interface{}{
			{
				"name":   "Asset",
				"value":  string(signal.Symbol),
				"inline": true,
			},
			{
				"name":   "Type",
				"value":  string(signal.Type),
				"inline": true,
			},
			{
				"name":   "Asset Class",
				"value":  string(signal.AssetClass),
				"inline": true,
			},
			{
				"name":   "Duration",
				"value":  string(signal.DurationType),
				"inline": true,
			},
		},
		"footer": map[string]string{
			"text": "Check the app for full details",
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	reqBody := map[string]interface{}{
		"content": "",
		"embeds":  []interface{}{embed},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal discord request: %w", err)
	}

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create discord request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("discord API returned status %d", resp.StatusCode)
	}

	log.Printf("Discord notification sent for signal ID %d", signal.ID)
	return nil
}

// ExpoNotificationService is a placeholder for Expo push notifications
type ExpoNotificationService struct{}

func NewExpoNotificationService() *ExpoNotificationService {
	return &ExpoNotificationService{}
}

func (s *ExpoNotificationService) SendSignalNotification(signal *models.TradingSignal) error {
	// TODO: Implement Expo push notification
	// This is a placeholder for future implementation
	log.Printf("[EXPO PLACEHOLDER] Would send notification for signal ID %d", signal.ID)
	log.Printf("[EXPO PLACEHOLDER] Signal: %s %s - %s", signal.AssetClass, signal.DurationType, signal.Symbol)
	
	// When implementing:
	// 1. Store user Expo push tokens in database
	// 2. Use Expo SDK or API to send push notifications
	// 3. Send to users who have active subscriptions for this signal type
	
	return nil
}

