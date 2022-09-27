package tcpserver

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type handler struct {
	cn       net.Conn
	protocol Protocol
	iwork    Work
	key      string
	password string
}

func newHandler(cn net.Conn, protocol Protocol, iwork Work, password string) *handler {
	data, key := protocol.Response(map[string]string{})
	cn.Write(data)
	return &handler{
		cn:       cn,
		protocol: protocol,
		iwork:    iwork,
		key:      key,
		password: password,
	}
}

func (h *handler) handle(wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc) {
	defer wg.Done()
	defer h.cn.Close()
	isWork := true
	for isWork {
		select {
		case <-ctx.Done():
			isWork = false
		default:
			netdata, err := read(h.cn)
			if err != nil {
				fmt.Println(err)
				return
			}
			data, err := h.protocol.Request(netdata, h.key)
			if err != nil {
				data, key := h.protocol.ResponseError(err)
				h.key = key
				h.cn.Write(data)
			}
			switch data["command"] {
			case "get_phrase":
				phrase, _ := h.iwork.DoWork(data)
				data, key := h.protocol.Response(map[string]string{
					"phrase": phrase,
				})
				h.key = key
				h.cn.Write(data)
			case "shutdown":
				if data["password"] == h.password {
					cancel()
					isWork = false
					data, key := h.protocol.Response(map[string]string{"Message": "Server shutdown"})
					h.key = key
					h.cn.Write(data)

				} else {
					data, key := h.protocol.ResponseError(fmt.Errorf("wrong password"))
					h.key = key
					h.cn.Write(data)
				}

			}
		}
	}
}

func read(conn net.Conn) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}
