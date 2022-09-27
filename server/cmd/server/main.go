package main

import (
	"os"
	"strconv"

	learnthephrase "github.com/CrazyHoneyBadger/Client_Server_TCP/server/Internal/learn_the_phrase"
	tcpserver "github.com/CrazyHoneyBadger/Client_Server_TCP/server/Internal/tcp_server"
	tcprotoko "github.com/CrazyHoneyBadger/TCP_Protocol"
	"github.com/CrazyHoneyBadger/pow"
)

func main() {
	lengthKeyStr := getEnv("LENGTH_KEY", "40")
	lengthKey, err := strconv.Atoi(lengthKeyStr)
	if err != nil {
		lengthKey = 40
	}
	complexityStr := getEnv("COMPLEXITY", "5")
	complexity, err := strconv.Atoi(complexityStr)
	if err != nil {
		complexity = 5
	}
	port := getEnv("PORT", ":8080")
	powServer := pow.NewPOWServer(lengthKey, complexity)

	proto := tcprotoko.NewProtocolServer(powServer)
	phrase := learnthephrase.NewPharse()
	tcpServer := tcpserver.NewTCPServer(proto, phrase)
	password := getEnv("PASSWORD", "")
	tcpServer.Start(port, password)
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
