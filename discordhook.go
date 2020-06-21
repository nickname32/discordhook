package discordhook

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/andersfylling/snowflake"
	jsoniter "github.com/json-iterator/go"
)

// WebhookAPI - for those, who wants to use single webhook url several times
type WebhookAPI struct {
	// URL of webhook
	URL *url.URL
	// Client which will be used for http requests
	Client *http.Client
	// Wait - if `true`, then the response will be expected and parsed in `msg`.
	Wait bool
}

// NewWebhookAPI creates WebhookExecuter (https://discord.com/api/webhooks/WEBHOOK_ID/WEBHOOK_TOKEN).
// if `wait` is `true`, then the response for method `Execute` will be expected and parsed in `msg`.
// client - *http.Client which will be used for http requests
func NewWebhookAPI(webhookID snowflake.Snowflake, webhookToken string, wait bool, client *http.Client) (*WebhookAPI, error) {
	u, err := url.Parse("https://discord.com/api/webhooks/" + strconv.FormatUint(uint64(webhookID), 10) + "/" + webhookToken + "?wait=" + strconv.FormatBool(wait))
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = http.DefaultClient
	}

	return &WebhookAPI{
		URL:    u,
		Client: client,
		Wait:   wait,
	}, nil
}

// WebhookExecuteParams represents webhook params payload structure
// https://discord.com/developers/docs/resources/webhook#execute-webhook-jsonform-params
type WebhookExecuteParams struct {
	// Content - the message contents (up to 2000 characters) [Required: one of content, file, embeds]
	Content string `json:"content,omitempty"`
	// Username - override the default username of the webhook [Required: false]
	Username string `json:"username,omitempty"`
	// AvatarURL - override the default avatar of the webhook [Required: false]
	AvatarURL string `json:"avatar_url,omitempty"`
	// TTS - true if this is a TTS message [Required: false]
	TTS bool `json:"tts,omitempty"`
	// Embeds - array of up to 10 embed objects [Required: one of content, file, embeds]
	Embeds []*Embed `json:"embeds,omitempty"`
	// AllowedMentions - allowed mentions for the message [Required: false]
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
}

// Execute - executes webhook. If `wait` was set to `true`, then the response will be expected and parsed in `msg`.
// Return message, http status code, error
func (wa *WebhookAPI) Execute(ctx context.Context, wep *WebhookExecuteParams, file io.Reader, filename string) (*Message, error) {
	bodyBuf := bytes.NewBuffer([]byte{})

	mw := multipart.NewWriter(bodyBuf)

	payloadPart, err := mw.CreateFormField("payload_json")
	if err != nil {
		return nil, err
	}

	err = jsoniter.NewEncoder(payloadPart).Encode(wep)
	if err != nil {
		return nil, err
	}

	if file != nil {
		filePart, err := mw.CreateFormFile("file", filename)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(filePart, file)
		if err != nil {
			return nil, err
		}
	}

	err = mw.Close()
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    wa.URL,
		Header: http.Header{
			"Content-Type": {mw.FormDataContentType()},
		},
		Body: ioutil.NopCloser(bodyBuf),
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := wa.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 && res.StatusCode != 201 && res.StatusCode != 204 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(b))
	}

	if wa.Wait {
		msg := new(Message)
		err = jsoniter.NewDecoder(res.Body).Decode(msg)
		if err != nil {
			return nil, err
		}

		return msg, nil
	}

	return nil, nil
}

// Webhook represents Webhook Structure
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-structure
type Webhook struct {
	// ID (snowflake) - the id of the webhook
	ID snowflake.Snowflake `json:"id"`
	// Type - the type of the webhook
	Type WebhookType `json:"type"`
	// GuildID (snowflake) - the guild id this webhook is for
	GuildID snowflake.Snowflake `json:"guild_id,omitempty"`
	// ChannelID (snowflake) - the channel id this webhook is for
	ChannelID snowflake.Snowflake `json:"channel_id"`
	// User - the user this webhook was created by (not returned when getting a webhook with its token)
	User *User `json:"user,omitempty"`
	// Name - the default name of the webhook
	Name string `json:"name,omitempty"`
	// Avatar - the default avatar of the webhook
	Avatar string `json:"avatar,omitempty"`
	// Token - the secure token of the webhook (returned for Incoming Webhooks)
	Token string `json:"token,omitempty"`
}

// WebhookType used for describing webhook types
// https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
type WebhookType int

const (
	// WebhookTypeIncoming - incoming Webhooks can post messages to channels with a generated token
	WebhookTypeIncoming WebhookType = 1
	// WebhookTypeChannelFollower - channel Follower Webhooks are internal webhooks used with Channel Following to post new messages into channels
	WebhookTypeChannelFollower WebhookType = 2
)

// Get - gets information about webhook.
// Return webhook, http status code, error
func (wa *WebhookAPI) Get(ctx context.Context) (*Webhook, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    wa.URL,
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := wa.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(b))
	}

	wh := new(Webhook)
	err = jsoniter.NewDecoder(res.Body).Decode(wh)
	if err != nil {
		return nil, err
	}

	return wh, nil
}

// WebhookModifyParams represents modify params for webhook api
// https://discord.com/developers/docs/resources/webhook#modify-webhook-json-params
type WebhookModifyParams struct {
	// Name - the default name of the webhook
	Name string `json:"name,omitemprty"`
	// Avatar - image for the default webhook avatar
	// Look https://discord.com/developers/docs/reference#image-data
	Avatar string `json:"avatar,omitemprty"`
	// ChannelID - the new channel id this webhook should be moved to
	ChannelID snowflake.Snowflake `json:"channel_id,omitemprty"`
}

// Modify - modifies a webhook.
// Return webhook, http status code, error
func (wa *WebhookAPI) Modify(ctx context.Context, wmp *WebhookModifyParams) (*Webhook, error) {
	bodyBuf := bytes.NewBuffer([]byte{})

	err := jsoniter.NewEncoder(bodyBuf).Encode(wmp)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPatch,
		URL:    wa.URL,
		Body:   ioutil.NopCloser(bodyBuf),
		Header: http.Header{
			"Content-Type": {"application/json"},
		},
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := wa.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 && res.StatusCode != 304 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(b))
	}

	wh := new(Webhook)
	err = jsoniter.NewDecoder(res.Body).Decode(wh)
	if err != nil {
		return nil, err
	}

	return wh, nil
}

// Delete - deletes a webhook.
// Returns http status code, error
func (wa *WebhookAPI) Delete(ctx context.Context) error {
	req := &http.Request{
		Method: http.MethodDelete,
		URL:    wa.URL,
	}

	res, err := wa.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 204 && res.StatusCode != 200 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(b))
	}

	return nil
}

// TODO: ExecuteSlack method
// TODO: ExecuteGitHub method
