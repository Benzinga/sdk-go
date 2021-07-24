package news

import (
	"context"
	"errors"
	"fmt"

	"github.com/Benzinga/sdk-go/pkg/models/websocket/news"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var WebsocketReadLimit int64 = 1000 << 16

var ErrNilWebsocketMessageHandler = errors.New("websocket message handler must not be nil")

type WebsocketMessageHandler interface {
	Handle(b *news.Body) error
}

func RunWebsocket(ctx context.Context, url string, mh WebsocketMessageHandler) error {
	if mh == nil {
		return ErrNilWebsocketMessageHandler
	}

	wsc, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("websocket dial error: %w", err)
	}

	wsc.SetReadLimit(WebsocketReadLimit)

	for {
		var m news.Body

		if err := wsjson.Read(ctx, wsc, &m); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("websocket read message error: %w", err)
		}

		if err := mh.Handle(&m); err != nil {
			return fmt.Errorf("handle message error: %w", err)
		}
	}
}
