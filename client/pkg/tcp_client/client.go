package tcpclient

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type (
	protocol interface {
		Request(data []byte) (map[string]string, error)
		Response(data map[string]string, key string, complexity int) ([]byte, error)
	}
	TCPClient struct {
		prtc   protocol
		adress string
	}
)

func NewTCPClient(adress string, prtc protocol) *TCPClient {
	return &TCPClient{adress: adress, prtc: prtc}
}
func (c TCPClient) SendResponse(dataRequest map[string]string) (map[string]string, error) {
	conn, err := net.Dial("tcp", c.adress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	buf, err := read(conn)
	if err != nil {
		return nil, err
	}
	data, err := c.prtc.Request(buf)
	if err != nil {
		return nil, err
	}
	complexityString := data["complexity"]
	complexity, _ := strconv.Atoi(complexityString)
	for key, value := range dataRequest {
		data[key] = value
	}
	buf, err = c.prtc.Response(data, data["POW_KEY"], complexity)
	if err != nil {
		return nil, err
	}
	conn.Write(buf)
	buf, err = read(conn)
	if err != nil {
		return nil, err
	}
	data, err = c.prtc.Request(buf)
	if err != nil {
		return nil, err
	}
	if er, ok := data["error"]; ok {
		return nil, fmt.Errorf(er)
	}
	conn.Close()
	return data, nil
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
