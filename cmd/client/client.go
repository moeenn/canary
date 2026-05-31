package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
)

const pipeChannelBufferLength = 20

type config struct {
	AppName       string // e.g. my-app
	ServerAddress string // e.g. localhost:8000
}

func newConfig() (*config, error) {
	appName := os.Getenv("CANARY_APP")
	if appName == "" {
		return nil, errors.New("missing environment variable: CANARY_APP")
	}

	serverAdr := os.Getenv("CANARY_SERVER")
	if serverAdr == "" {
		return nil, errors.New("missin environment variable: CANARY_SERVER")
	}

	cfg := config{
		AppName:       strings.TrimSpace(appName),
		ServerAddress: strings.TrimSpace(serverAdr),
	}

	return &cfg, nil
}

type logPayload struct {
	App     string `json:"app"`
	Payload string `json:"payload"`
}

func readPipedInput() <-chan string {
	result := make(chan string, pipeChannelBufferLength)
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		var line string
		for scanner.Scan() {
			line = scanner.Text()
			result <- line
		}

		if err := scanner.Err(); err != nil {
			close(result)
		}
	}()

	return result
}

func connectToServer(serverAddress string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve server address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	return conn, nil
}

func sendLogMessage(logger *slog.Logger, conn *net.UDPConn, appName string, input string) {
	payload := logPayload{App: appName, Payload: input}
	encoded, err := json.Marshal(payload)
	if err != nil {
		logger.Error("failed to json encode log payload", "error", err.Error())
		return
	}

	if _, err := conn.Write(encoded); err != nil {
		logger.Error("failed to send packet to the canary server", "error", err.Error())
	}
}

func run() error {
	cfg, err := newConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	conn, err := connectToServer(cfg.ServerAddress)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer conn.Close()
	for input := range readPipedInput() {
		go sendLogMessage(logger, conn, cfg.AppName, input)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s.\n", err.Error())
		os.Exit(1)
	}
}
