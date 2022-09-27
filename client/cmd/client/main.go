package main

import (
	"fmt"
	"os"

	tcpclient "github.com/CrazyHoneyBadger/Client_Server_TCP/client/pkg/tcp_client"
	tcpprotocol "github.com/CrazyHoneyBadger/TCP_Protocol"
	"github.com/CrazyHoneyBadger/pow"
)

func main() {
	netAddr := getEnv("NET_ADDR", "localhost:8080")
	isShutdown := getEnv("SHUTDOWN", "false")
	password := getEnv("PASSWORD", "12345618")
	powClient := pow.NewPOWClient(1 << 23)
	proto := tcpprotocol.NewProtocolClient(powClient)
	tcp := tcpclient.NewTCPClient(netAddr, proto)
	phrase, err := getPhrase(tcp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(phrase)
	if isShutdown == "true" {
		err = shutdownServer(tcp, password)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func getPhrase(tcp *tcpclient.TCPClient) (string, error) {
	data, err := tcp.SendResponse(map[string]string{"command": "get_phrase"})
	if err != nil {
		return "", err
	}
	return data["phrase"], nil
}
func shutdownServer(tcp *tcpclient.TCPClient, password string) error {
	data, err := tcp.SendResponse(map[string]string{"command": "shutdown", "password": password})
	if err != nil {
		return err
	}
	fmt.Println(data["Message"])
	return nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
