package benzinga

import "github.com/Benzinga/sdk-go/pkg/client/rest/news"

func (c *Client) News() *news.Request {
	return news.NewRequest(c.c)
}
