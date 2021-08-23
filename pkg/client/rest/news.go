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
	apiToken string
	hostname string
}

func (c *Client) NewRequest() *Request {
	return &Request{c, Params{}, "", ""}
}

func (r *Request) SetAPIToken(token string) *Request {
	r.apiToken = token
	return r
}

func (r *Request) AddChannels(channels ...string) *Request {
	r.params.Channels = append(r.params.Channels, channels...)
	return r
}

func (r *Request) SetChannels(channels ...string) *Request {
	r.params.Channels = channels
	return r
}

func (r *Request) AddTickers(tickers ...string) *Request {
	r.params.Tickers = append(r.params.Tickers, tickers...)
	return r
}

func (r *Request) SetTickers(tickers ...string) *Request {
	r.params.Tickers = tickers
	return r
}

func (r *Request) AddTopics(topics ...string) *Request {
	r.params.Topics = append(r.params.Topics, topics...)
	return r
}

func (r *Request) SetTopics(topics ...string) *Request {
	r.params.Topics = topics
	return r
}

func (r *Request) AddCUSIPs(cusips ...string) *Request {
	r.params.CUSIPs = append(r.params.Topics, cusips...)
	return r
}

func (r *Request) SetCUSIPs(cusips ...string) *Request {
	r.params.CUSIPs = cusips
	return r
}

func (r *Request) SetDate(date time.Time) *Request {
	r.params.Date = date
	return r
}

func (r *Request) SetDateFrom(dateFrom time.Time) *Request {
	r.params.DateFrom = dateFrom
	return r
}

func (r *Request) SetDateTo(dateTo time.Time) *Request {
	r.params.DateTo = dateTo
	return r
}

func (r *Request) SetPublishedSince(publishedSince time.Time) *Request {
	r.params.PublishedSince = publishedSince
	return r
}

func (r *Request) SetUpdatedSince(updatedSince time.Time) *Request {
	r.params.UpdatedSince = updatedSince
	return r
}

func (r *Request) SetDisplayOutput(outputOption OutputOption) *Request {
	r.params.DisplayOutput = &outputOption
	return r
}

func (r *Request) SetPageSize(n int) *Request {
	r.params.PageSize = n
	return r
}

func (r *Request) SetPage(n int) *Request {
	r.params.Page = n
	return r
}

func (r *Request) SetSortField(f SortField) *Request {
	if r.params.Sort == nil {
		r.params.Sort = &SortOption{}
	}

	r.params.Sort.Field = f

	return r
}

func (r *Request) SetSortDirection(d SortDirection) *Request {
	if r.params.Sort == nil {
		r.params.Sort = &SortOption{}
	}

	r.params.Sort.Direction = d

	return r
}

func (r *Request) URL() (*url.URL, error) {
	u, err := url.Parse(r.hostname + NewsAPIPath)
	if err != nil {
		return nil, fmt.Errorf("parse url error: %w", err)
	}

	q := u.Query()

	q.Set("token", r.apiToken)

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
