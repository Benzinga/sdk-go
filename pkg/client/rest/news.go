package news

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Benzinga/sdk-go/pkg/models/rest/news"
)

const (
	NewsAPIPath = "/api/v2/news"
	DateFormat  = "2006-01-02"
)

type Client struct {
	c *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		c: httpClient,
	}
}

type Params struct {
	Channels       []string      `form:"channels"`
	Date           time.Time     `form:"date" time_format:"2006-01-02" time_location:"Etc/GMT+4"`
	DateFrom       time.Time     `form:"dateFrom" time_format:"2006-01-02" time_location:"Etc/GMT+4"`
	DateTo         time.Time     `form:"dateTo" time_format:"2006-01-02" time_location:"Etc/GMT+4"`
	PublishedSince time.Time     `form:"publishedSince" time_format:"unix" time_utc:"1"`
	UpdatedSince   time.Time     `form:"updatedSince" time_format:"unix" time_utc:"1"`
	DisplayOutput  *OutputOption `form:"displayOutput" validate:"alpha"`
	Page           int           `form:"page" validate:"lte=10000"`
	PageSize       int           `form:"pageSize"`
	Sort           *SortOption   `form:"sort"`
	Tickers        []string      `form:"tickers"`
	Topics         []string      `form:"topics"`
	CUSIPs         []string      `form:"cusips"`
	APIToken       string        //TODO(darwin)
}

type OutputOption int

const (
	HeadlineOutput OutputOption = iota
	AbstractOutput
	FullOutput
)

func (o OutputOption) String() string {
	return [...]string{"headline", "abstract", "full"}[o]
}

type SortField int

const (
	UpdatedField SortField = iota
	CreatedField
	IdField
)

func (s SortField) String() string {
	return [...]string{"updated", "created", "id"}[s]
}

type SortDirection int

const (
	Ascending = iota
	Descending
)

func (s SortDirection) String() string {
	return [...]string{"asc", "desc"}[s]
}

type SortOption struct {
	Field     SortField
	Direction SortDirection
}

func (s SortOption) String() string {
	return s.Field.String() + ":" + s.Direction.String()
}

type Request struct {
	*Client
	params   Params
	hostname string
}

func (c *Client) NewRequest() *Request {
	return &Request{c, Params{}, ""}
}

func (r *Request) AddChannels(channels ...string) {
	r.params.Channels = append(r.params.Channels, channels...)
}

func (r *Request) SetChannels(channels ...string) {
	r.params.Channels = channels
}

func (r *Request) AddTickers(tickers ...string) {
	r.params.Tickers = append(r.params.Tickers, tickers...)
}

func (r *Request) SetTickers(tickers ...string) {
	r.params.Tickers = tickers
}

func (r *Request) AddTopics(topics ...string) {
	r.params.Topics = append(r.params.Topics, topics...)
}

func (r *Request) SetTopics(topics ...string) {
	r.params.Topics = topics
}

func (r *Request) AddCUSIPs(cusips ...string) {
	r.params.CUSIPs = append(r.params.Topics, cusips...)
}

func (r *Request) SetCUSIPs(cusips ...string) {
	r.params.CUSIPs = cusips
}

func (r *Request) SetDate(date time.Time) {
	r.params.Date = date
}

func (r *Request) SetDateFrom(dateFrom time.Time) {
	r.params.DateFrom = dateFrom
}

func (r *Request) SetDateTo(dateTo time.Time) {
	r.params.DateTo = dateTo
}

func (r *Request) SetPublishedSince(publishedSince time.Time) {
	r.params.PublishedSince = publishedSince
}

func (r *Request) SetUpdatedSince(updatedSince time.Time) {
	r.params.UpdatedSince = updatedSince
}

func (r *Request) SetDisplayOutput(outputOption OutputOption) {
	r.params.DisplayOutput = &outputOption
}

func (r *Request) SetPageSize(n int) {
	r.params.PageSize = n
}

func (r *Request) SetPage(n int) {
	r.params.Page = n
}

func (r *Request) SetSortField(f SortField) {
	if r.params.Sort == nil {
		r.params.Sort = &SortOption{}
	}

	r.params.Sort.Field = f
}

func (r *Request) SetSortDirection(d SortDirection) {
	if r.params.Sort == nil {
		r.params.Sort = &SortOption{}
	}

	r.params.Sort.Direction = d
}

func (r *Request) URL() (*url.URL, error) {
	u, err := url.Parse(r.hostname + NewsAPIPath)
	if err != nil {
		return nil, fmt.Errorf("parse url error: %w", err)
	}

	q := u.Query()

	if len(r.params.Channels) > 0 {
		q.Set("channels", strings.Join(r.params.Channels, ","))
	}

	if len(r.params.Tickers) > 0 {
		q.Set("tickers", strings.Join(r.params.Tickers, ","))
	}

	if len(r.params.Topics) > 0 {
		q.Set("topics", strings.Join(r.params.Topics, ","))
	}

	if len(r.params.CUSIPs) > 0 {
		q.Set("cusips", strings.Join(r.params.CUSIPs, ","))
	}

	if !r.params.Date.IsZero() {
		q.Set("date", r.params.Date.Format(DateFormat))
	}

	if !r.params.DateFrom.IsZero() {
		q.Set("dateFrom", r.params.DateFrom.Format(DateFormat))
	}

	if !r.params.DateTo.IsZero() {
		q.Set("dateTo", r.params.DateTo.Format(DateFormat))
	}

	if !r.params.PublishedSince.IsZero() {
		q.Set("publishedSince", strconv.FormatInt(r.params.PublishedSince.Unix(), 10))
	}

	if !r.params.UpdatedSince.IsZero() {
		q.Set("updatedSince", strconv.FormatInt(r.params.UpdatedSince.Unix(), 10))
	}

	if r.params.DisplayOutput != nil {
		q.Set("displayOutput", r.params.DisplayOutput.String())
	}

	if r.params.Page > 0 {
		q.Set("page", strconv.Itoa(r.params.Page))
	}

	if r.params.PageSize > 0 {
		q.Set("pageSize", strconv.Itoa(r.params.PageSize))
	}

	if r.params.Sort != nil {
		q.Set("sort", r.params.Sort.String())
	}

	u.RawQuery = q.Encode()

	return u, nil
}

type ErrUnexpectedResponse struct {
	StatusCode int
	Message    string
}

func (e ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("unexpected api response with code: %d, message: %s", e.StatusCode, e.Message)
}

func (r *Request) Exec(ctx context.Context) (news.Stories, error) {
	u, err := r.URL()
	if err != nil {
		return nil, fmt.Errorf("error building request url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	// Gzip?

	resp, err := r.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, ErrUnexpectedResponse{
			StatusCode: resp.StatusCode,
			Message:    string(b),
		}
	}

	var stories news.Stories

	if err := json.Unmarshal(b, &stories); err != nil {
		return nil, fmt.Errorf("error parsing json response: %w", err)
	}

	return stories, nil
}
