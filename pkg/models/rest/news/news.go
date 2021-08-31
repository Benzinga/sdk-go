package news

import (
	"encoding/json"
	"time"
)

const NewsTimeFormat = "Mon, 02 Jan 2006 15:04:05 -0700"

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
		Created string `json:"created"`
		Updated string `json:"updated"`
		*TempStory
	}{
		TempStory: (*TempStory)(s),
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	s.Created, err = time.Parse(NewsTimeFormat, tmp.Created)
	if err != nil {
		return err
	}

	s.Updated, err = time.Parse(NewsTimeFormat, tmp.Updated)
	if err != nil {
		return err
	}

	return nil
}
