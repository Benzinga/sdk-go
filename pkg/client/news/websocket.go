package news

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Benzinga/sdk-go/pkg/models/websocket/news"

	bolt "go.etcd.io/bbolt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Sets Reader Limit for Websocket Messages.
var WebsocketReadLimit int64 = 1000 << 16

const websocketBufferBucketName = "news_websocket"

var (
	ErrNilWebsocketMessageHandler = errors.New("websocket message handler must not be nil")
)

// ErrServerError indicates some issue with the server and should provide hints on how to best handle.
type ErrServerError struct {
	Code    int
	Message string
}

func (e ErrServerError) Error() string {
	return fmt.Sprintf("server returned unexpected error: %d - %s", e.Code, e.Message)
}

type WebsocketMessageHandler interface {
	Handle(b *news.Body) error
}

// Websocket contains configuration values for News Websocket connections.
type Websocket struct {
	url  string
	mh   WebsocketMessageHandler
	opts *WebsocketOptions
}

type WebsocketOptions struct {
	AutoClearBuffer bool
	UseDiskBuffer   bool
	DiskBufferPath  string
	UseMemoryBuffer bool
	buffer          []news.Body
	diskBuffer      *bolt.DB
}

func NewWebsocket(url string, mh WebsocketMessageHandler, opts *WebsocketOptions) (*Websocket, error) {
	if mh == nil {
		return nil, ErrNilWebsocketMessageHandler
	}

	if opts == nil {
		opts = &WebsocketOptions{}
	}

	if err := opts.bufferSetup(); err != nil {
		return nil, fmt.Errorf("buffer configuration failed with error: %w", err)
	}

	return &Websocket{url, mh, opts}, nil
}

func (w *WebsocketOptions) bufferSetup() error {
	if w.UseMemoryBuffer && w.buffer == nil {
		w.buffer = make([]news.Body, 0)
	}

	if w.UseDiskBuffer {
		if w.DiskBufferPath == "" {
			w.DiskBufferPath = path.Join(os.TempDir(), "bz_news_websocket_buffer.db")
		}

		bopts := bolt.DefaultOptions
		bopts.Timeout = 10 * time.Second

		db, err := bolt.Open(w.DiskBufferPath, 0600, bopts)
		if err != nil {
			return err
		}

		w.diskBuffer = db
	}

	return nil
}

func (w *Websocket) Run(ctx context.Context) error {
	wsc, r, err := websocket.Dial(ctx, w.url, nil)
	if err != nil {
		return fmt.Errorf("websocket dial error: %w", err)
	}

	switch r.StatusCode {
	case http.StatusTooManyRequests:
		return ErrServerError{Code: r.StatusCode, Message: "server reports too many connections, wait before reconnecting or disconnect other sessions"}
	case http.StatusInternalServerError, http.StatusServiceUnavailable:
		return ErrServerError{Code: r.StatusCode, Message: "server unavaiable, delay before attempting reconnect"}
	}

	wsc.SetReadLimit(WebsocketReadLimit)

	for {
		var b *news.Body

		if err := wsjson.Read(ctx, wsc, &b); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("websocket read message error: %w", err)
		}

		if err := w.SaveToBuffer(b); err != nil {
			return err
		}

		if err := w.mh.Handle(b); err != nil {
			return fmt.Errorf("handle message error: %w", err)
		}
	}
}

func (w *Websocket) RemoveFromBuffer(b *news.Body) error {
	if w.opts.UseMemoryBuffer {
		for i := 0; i < len(w.opts.buffer); i++ {
			if w.opts.buffer[i].Data.ID == b.Data.ID {
				w.opts.buffer = append(w.opts.buffer[:i], w.opts.buffer[i+1:]...) // TODO: Fix this. Handle slice length issues.
			}
		}
	}

	return nil
}

func (w *Websocket) RetrieveFromBuffer() {

}

func (w *Websocket) SaveToBuffer(b *news.Body) error {
	if w.opts.UseMemoryBuffer {
		w.opts.buffer = append(w.opts.buffer, *b)
	}

	if w.opts.UseDiskBuffer {
		err := w.opts.diskBuffer.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(websocketBufferBucketName))
			if err != nil {
				return err
			}

			val, err := json.Marshal(b)
			if err != nil {
				return err
			}

			if err := bucket.Put([]byte(strconv.FormatInt(b.Data.ID, 10)), val); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
