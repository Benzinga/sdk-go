package news

import "github.com/Benzinga/sdk-go/pkg/bztime"

type Body struct {
	APIVersion string `json:"api_version"`
	Kind       string `json:"kind"`
	Data       *Data  `json:"data"`
}

type Security struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
	Primary  bool   `json:"primary"`
}

type Content struct {
	ID         int    `json:"id"`
	RevisionID int    `json:"revision_id"`
	Type       string `json:"type"`

	Title      string     `json:"title"`
	Body       string     `json:"body"`
	Authors    []string   `json:"authors"`
	Teaser     string     `json:"teaser"`
	URL        string     `json:"url"`
	Tags       []string   `json:"tags"`
	Securities []Security `json:"securities"`
	Channels   []string   `json:"channels"`

	CreatedAt bztime.Time `json:"created_at"`
	UpdatedAt bztime.Time `json:"updated_at"`
}

type Data struct {
	Action    string      `json:"action"`
	ID        int64       `json:"id"`
	Content   *Content    `json:"content"`
	Timestamp bztime.Time `json:"timestamp"`
}
