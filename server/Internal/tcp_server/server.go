package tcpserver

import (
	"context"
	"net"
	"sync"
)

type (
	Protocol interface {
		Request(data []byte, key string) (map[string]string, error)
		Response(data map[string]string) ([]byte, string)
		ResponseError(err error) ([]byte, string)
	}
	Work interface {
		DoWork(payload map[string]string) (string, error)
	}
	TCPServer struct {
		protocol Protocol
		iwork    Work
	}
)

func NewTCPServer(protocol Protocol, iwork Work) *TCPServer {
	return &TCPServer{
		protocol: protocol,
		iwork:    iwork,
	}
}
func (s TCPServer) Start(listenPort string, password string) error {
	server, err := net.Listen("tcp", listenPort)
	if err != nil {
		return err
	}
	defer server.Close()
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := server.Accept()
			if err != nil {
				return err
			}
			handler := newHandler(conn, s.protocol, s.iwork, password)
			wg.Add(1)
			go handler.handle(wg, ctx, cancel)
		}

	}
}
