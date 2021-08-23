package news

import (
	"encoding/json"
	"time"
)

type Image struct {
	Size string `json:"size"`
	URL  string `json:"url"`
}

// Stock ...
type Stock struct {
	Name  string `json:"name"`
	CUSIP string `json:"cusip,omitempty"`
}

// ChannelTag is a shared format for Channels and Tags
type ChannelTag struct {
	Name string `json:"name"`
}

type Story struct {
	ID       int          `json:"id"`
	Author   string       `json:"author"`
	Created  time.Time    `json:"created"`
	Updated  time.Time    `json:"updated"`
	Title    string       `json:"title"`
	Teaser   string       `json:"teaser"`
	Body     string       `json:"body"`
	URL      string       `json:"url"`
	Image    []Image      `json:"image"`
	Channels []ChannelTag `json:"channels"`
	Stocks   []Stock      `json:"stocks"`
	Tags     []ChannelTag `json:"tags"`
}

type Stories []Story

// Implement unmarshal for Created, Updated.
func (s *Story) UnmarshalJSON(data []byte) error {
	type TempStory Story

	tmp := struct {
		Created int `json:"created,string"`
		Updated int `json:"updated,string"`
		*TempStory
	}{
		TempStory: (*TempStory)(s),
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	s.Created = time.Unix(int64(tmp.Created), 0)
	s.Updated = time.Unix(int64(tmp.Updated), 0)

	return nil
}
