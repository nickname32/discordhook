package discordhook

import (
	"time"

	"github.com/andersfylling/snowflake"
)

// Embed - embedded `rich` content
// For the webhook embed objects, you can set every
// field except `type` (it will be rich regardless of if you try to set it),
// `provider`, `video`, and any `height`, `width`, or `proxy_url` values for images.
// https://discord.com/developers/docs/resources/channel#embed-object-embed-structure
type Embed struct {
	// Title - title of embed
	Title string `json:"title,omitempty"`
	// Type - type of embed (always "rich" for webhook embeds)
	Type EmbedType `json:"type,omitempty"`
	// Description - description of embed
	Description string `json:"description,omitempty"`
	// URL - url of embed
	URL string `json:"url,omitempty"`
	// Timestamp - timestamp of embed content
	Timestamp *time.Time `json:"timestamp,omitempty"`
	// Color - color code of the embed
	Color int `json:"color,omitempty"`
	// Footer - footer information
	Footer *EmbedFooter `json:"footer,omitempty"`
	// Image - image information
	Image *EmbedImage `json:"image,omitempty"`
	// Thumbnail - thumbnail information
	Thumbnail *EmbedThumbnail `json:"thumbnail,omitempty"`
	// Video - video information
	Video *EmbedVideo `json:"video,omitempty"`
	// Provider - provider information
	Provider *EmbedProvider `json:"provider,omitempty"`
	// Author - author information
	Author *EmbedAuthor `json:"author,omitempty"`
	// Fields - fields information
	Fields []*EmbedField `json:"fields,omitempty"`
}

// EmbedType - embed types are "loosely defined" and, for the most part,
// are not used by Discord's clients for rendering. Embed attributes power what is rendered.
// Embed types should be considered deprecated and might be removed in a future API version.
// https://discord.com/developers/docs/resources/channel#embed-object-embed-types
type EmbedType string

const (
	// EmbedTypeRich - generic embed rendered from embed attributes
	EmbedTypeRich EmbedType = "rich"
	// EmbedTypeImage - image embed
	EmbedTypeImage EmbedType = "image"
	// EmbedTypeVideo - video embed
	EmbedTypeVideo EmbedType = "video"
	// EmbedTypeGifv - animated gif image embed rendered as a video embed
	EmbedTypeGifv EmbedType = "gifv"
	// EmbedTypeArticle - article embed
	EmbedTypeArticle EmbedType = "article"
	// EmbedTypeLink - link embed
	EmbedTypeLink EmbedType = "link"
)

// EmbedFooter represents Embed Footer Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-footer-structure
type EmbedFooter struct {
	// Text - footer text (required)
	Text string `json:"text"`
	// IconURL - url of footer icon (only supports http(s) and attachments)
	IconURL string `json:"icon_url,omitempty"`
	// ProxyIconURL - a proxied url of footer icon
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// EmbedImage represents Embed Image Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-image-structure
type EmbedImage struct {
	// URl source url of image (only supports http(s) and attachments)
	URL string `json:"url,omitempty"`
	// ProxyURL - a proxied url of the image
	ProxyURL string `json:"proxy_url,omitempty"`
	// Height - height of image
	Height int `json:"height,omitempty"`
	// Width - width of image
	Width int `json:"width,omitempty"`
}

// EmbedThumbnail represents Embed Thumbnail Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-thumbnail-structure
type EmbedThumbnail struct {
	// URL - source url of thumbnail (only supports http(s) and attachments)
	URL string `json:"url,omitempty"`
	// ProxyURL - a proxied url of the thumbnail
	ProxyURL string `json:"proxy_url,omitempty"`
	// Height - height of thumbnail
	Height int `json:"height,omitempty"`
	// Wifth - width of thumbnail
	Width int `json:"width,omitempty"`
}

// EmbedVideo represents Embed Video Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-video-structure
type EmbedVideo struct {
	// URL - source url of video
	URL string `json:"url,omitempty"`
	// Height - height of video
	Height int `json:"height,omitempty"`
	// Width - width of video
	Width int `json:"width,omitempty"`
}

// EmbedProvider represents Embed Provider Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-provider-structure
type EmbedProvider struct {
	// Name - name of provider
	Name string `json:"name,omitempty"`
	// URL - url of provider
	URL string `json:"url,omitempty"`
}

// EmbedAuthor - represenst Embed Author Structure
// https://discord.com/developers/docs/resources/channel#embed-object-embed-author-structure
type EmbedAuthor struct {
	// Name - name of author
	Name string `json:"name,omitempty"`
	// URL - url of author
	URL string `json:"url,omitempty"`
	// IconURL - url of author icon (only supports http(s) and attachments)
	IconURL string `json:"icon_url,omitempty"`
	// ProxyIconURL - a proxied url of author icon
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// EmbedField - used to represent an embed field object
// https://discord.com/developers/docs/resources/channel#embed-object-embed-field-structure
type EmbedField struct {
	// Name - name of the field (required)
	Name string `json:"name"`
	// Value - value of the field (required)
	Value string `json:"value"`
	// Inline - whether or not this field should display inline
	Inline bool `json:"inline,omitempty"`
}

// AllowedMentions allows for more granular control over
// mentions without various hacks to the message content.
// This will always validate against message content
// to avoid phantom pings (e.g. to ping everyone,
// you must still have @everyone in the message content),
// and check against user/bot permissions.
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object
type AllowedMentions struct {
	// Parse - an array of allowed mention types to parse from the content.
	Parse []MentionType `json:"parse,omitempty"`
	// Roles - an array of role_ids to mention (Max size of 100)
	Roles []snowflake.Snowflake `json:"roles,omitempty"`
	// Users - Array of user_ids to mention (Max size of 100)
	Users []snowflake.Snowflake `json:"users,omitempty"`
}

// MentionType - discord mention type
// https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mention-types
type MentionType string

const (
	// MentionTypeRole can be used in AllowedMentions to allow mention roles
	MentionTypeRole MentionType = "roles"
	// MentionTypeUser can be used in AllowedMentions to allow mention users
	MentionTypeUser MentionType = "users"
	// MentionTypeEveryone can be used in AllowedMentions to allow mention @everyone
	MentionTypeEveryone MentionType = "everyone"
)
