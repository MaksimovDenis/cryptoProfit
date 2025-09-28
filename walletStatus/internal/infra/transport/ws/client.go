package ws

import (
	"context"
	"net/url"
	"sync"
	"walletStatus/internal/domain"
	"walletStatus/internal/infra/logger"

	"go.uber.org/atomic"

	"github.com/gorilla/websocket"
)

type Client struct {
	addr      *url.URL
	ctx       context.Context
	ctxCancel context.CancelFunc

	wg        *sync.WaitGroup
	connAlive *atomic.Bool
	stopped   *atomic.Bool
	conn      *websocket.Conn

	In  chan []byte
	Out chan<- []byte
}

type Opts struct {
	Ctx  context.Context
	Addr string
	In   chan []byte
	Out  chan<- []byte
}

func NewClient(opts Opts) (*Client, error) {
	addr, err := url.Parse(opts.Addr)
	if err != nil {
		return nil, domain.ErrFailedToParseAddr
	}

	client := &Client{
		addr:      addr,
		connAlive: atomic.NewBool(false),
		stopped:   atomic.NewBool(false),
		wg:        &sync.WaitGroup{},
		In:        opts.In,
		Out:       opts.Out,
	}

	client.ctx, client.ctxCancel = context.WithCancel(opts.Ctx)

	return client, nil
}

func (hdl *Client) Start() {
	go hdl.startDialer()
	go hdl.startWriter()
	go hdl.startListener()
	go hdl.startPinger()
}

func (hdl *Client) Stop() {
	logger.Debugf(hdl.ctx, "call stop")

	hdl.ctxCancel()
	hdl.stopped.Store(true)
	hdl.wg.Wait()
}

func (hdl *Client) Alive() bool { return hdl.connAlive.Load() }
