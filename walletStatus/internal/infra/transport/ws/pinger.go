package ws

import (
	"time"
	"walletStatus/internal/infra/logger"

	"github.com/gorilla/websocket"
)

const pingInterval = 5 * time.Second

func (hdl *Client) startPinger() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hdl.ctx.Done():
			return
		case <-ticker.C:
			if hdl.stopped.Load() {
				return
			}

			for !hdl.connAlive.Load() {
				if hdl.stopped.Load() {
					return
				}

				time.Sleep(500 * time.Millisecond)
			}

			if err := hdl.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Errorf(hdl.ctx, "failed to ping websocket %w", err)
				hdl.connAlive.Store(false)
			}
		}
	}
}
