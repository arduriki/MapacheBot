package youtube

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Client is a wrapper around the YouTube API client
type Client struct {
	service    *youtube.Service
	liveChatID string
}

// NewClient creates a new YouTube API client
func NewClient(ctx context.Context, apiKey, liveChatID string) (*Client, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube service: %w", err)
	}

	return &Client{
		service:    service,
		liveChatID: liveChatID,
	}, nil
}

// GetLiveChatID retrieves the live chat ID for a live stream video
func (c *Client) GetLiveChatID(ctx context.Context, videoID string) (string, error) {
	call := c.service.Videos.List([]string{"liveStreamingDetails"}).Id(videoID)
	resp, err := call.Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get video details: %w", err)
	}

	if len(resp.Items) == 0 {
		return "", fmt.Errorf("video not found: %s", videoID)
	}

	video := resp.Items[0]
	if video.LiveStreamingDetails == nil || video.LiveStreamingDetails.ActiveLiveChatId == "" {
		return "", fmt.Errorf("video is not a live stream or live chat is not available")
	}

	return video.LiveStreamingDetails.ActiveLiveChatId, nil
}

// GetMessages gets messages from the live chat
func (c *Client) GetMessages(ctx context.Context, pageToken string) (*youtube.LiveChatMessageListResponse, error) {
	call := c.service.LiveChatMessages.List(c.liveChatID, []string{"snippet", "authorDetails"})

	// Set page token if provided
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	return call.Context(ctx).Do()
}

// DeleteMessage deletes a message from the live chat
// Requires moderator privileges
func (c *Client) DeleteMessage(ctx context.Context, messageID string) error {
	err := c.service.LiveChatMessages.Delete(messageID).Context(ctx).Do()
	return err
}

// SendMessage sends a message to the live chat
// This will be used to send warning messages
func (c *Client) SendMessage(ctx context.Context, message string) error {
	_, err := c.service.LiveChatMessages.Insert(
		[]string{"snippet"},
		&youtube.LiveChatMessage{
			Snippet: &youtube.LiveChatMessageSnippet{
				LiveChatId: c.liveChatID,
				Type:       "textMessageEvent",
				TextMessageDetails: &youtube.LiveChatTextMessageDetails{
					MessageText: message,
				},
			},
		},
	).Context(ctx).Do()
	return err
}

// SetLiveChatID updates the live chat ID
func (c *Client) SetLiveChatID(liveChatID string) {
	c.liveChatID = liveChatID
}
