package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/gamechanger/dog-devolver/config"
	"github.com/gamechanger/dog-devolver/proxy"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("dog-devolver")
var logLevel = logging.DEBUG

func initLogger() error {
	var backend *logging.LogBackend
	format := logging.MustStringFormatter(
		"%{color}%{time:01/02 15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
	)

	if config.LOG_TO == "" {
		backend = logging.NewLogBackend(os.Stderr, "", 0)
	} else {
		logFile, err := os.OpenFile(config.LOG_TO, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			return err
		}
		backend = logging.NewLogBackend(logFile, "", 0)
	}

	formatter := logging.NewBackendFormatter(backend, format)

	backendLeveled := logging.AddModuleLevel(formatter)
	backendLeveled.SetLevel(logLevel, "")

	logging.SetBackend(backendLeveled)
	return nil
}

func handleData(incomingData []byte) {
	log.Debug("Received: %s", incomingData)
	go proxy.ProxyToDogStatsD(incomingData)
	for _, target := range proxy.STATSD_TARGETS {
		go proxy.ProxyToStatsD(incomingData, target)
	}
}

func initSocket() (*net.UDPConn, error) {
	port, err := strconv.Atoi(config.LISTEN_PORT)
	if err != nil {
		panic(fmt.Sprintf("Port config value %s cannot be cast to int", config.LISTEN_PORT))
	}
	ip := net.ParseIP(config.LISTEN_IP)
	bindSpec := net.UDPAddr{
		IP:   ip,
		Port: port,
	}
	log.Info("Opening UDP socket at %s:%d", ip, port)
	sock, err := net.ListenUDP("udp4", &bindSpec)
	return sock, err
}

func main() {
	err := initLogger()
	if err != nil {
		panic(fmt.Sprintf("Error initializing logger: %s", err))
	}

	sock, err := initSocket()
	if err != nil {
		panic(fmt.Sprintf("Error opening UDP socket: %s", err))
	}
	defer sock.Close()

	incomingData := make([]byte, 65507)

	for {
		bytesRead, returnAddr, err := sock.ReadFromUDP(incomingData)
		if err != nil {
			log.Warning("Error receiving packet from addr %s, received bytes %s, error was: %s", returnAddr, incomingData, err)
			continue
		}
		log.Info(fmt.Sprintf("Read %d bytes from address %s", bytesRead, returnAddr))
		go handleData(incomingData)
	}
}
