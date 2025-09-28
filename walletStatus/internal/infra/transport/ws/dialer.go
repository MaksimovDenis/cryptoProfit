package ws

import (
	"time"
	"walletStatus/internal/infra/logger"

	"github.com/gorilla/websocket"
)

const dialInterval = 5 * time.Second

func (hdl *Client) startDialer() {
	hdl.wg.Add(1)
	defer hdl.wg.Done()

	dialTicker := time.NewTicker(dialInterval)
	defer dialTicker.Stop()

	for {
		select {
		case <-hdl.ctx.Done():
			return
		case <-dialTicker.C:
			if hdl.connAlive.Load() {
				continue
			}

			logger.Debugf(hdl.ctx, "dial attempt")

			wSock, _, err := websocket.DefaultDialer.DialContext(hdl.ctx, hdl.addr.String(), nil)
			if err != nil {
				logger.Errorf(hdl.ctx, "dial attempt failed: %w", err)
				continue
			}

			if hdl.conn != nil {
				hdl.conn.Close()
			}

			hdl.conn = wSock
			hdl.connAlive.Store(true)

			logger.Debugf(hdl.ctx, "dial attempt successed")
		}
	}

}
