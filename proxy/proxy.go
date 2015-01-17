package proxy

import (
	"fmt"
	"net"
	"strconv"

	"github.com/gamechanger/dog-devolver/config"
	"github.com/gamechanger/dog-devolver/devolve"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("dog-devolver")

var statsdTargets = zipTargets(config.STATSD_HOSTS, config.STATSD_PORTS)

type targetSpec struct {
	IP   net.IP
	Host string
	Port int
}

func (spec targetSpec) getIP() (net.IP, error) {
	if spec.IP != nil {
		return spec.IP, nil
	}
	addrs, err := net.LookupIP(spec.Host)
	if err != nil {
		return nil, err
	}
	log.Debug("Resolved IP %s for name %s", addrs[0], spec.Host)
	return addrs[0], nil
}

func zipTargets(hosts, ports []string) []targetSpec {
	var result []targetSpec
	for idx, host := range hosts {
		if len(ports) <= idx {
			return result
		}

		ip := net.ParseIP(host)
		port, err := strconv.Atoi(ports[idx])
		if err != nil {
			panic(fmt.Sprintf("Invalid port spec %s", ports[idx]))
		}

		if ip == nil {
			result = append(result, targetSpec{Host: host, Port: port})
		} else {
			result = append(result, targetSpec{IP: ip, Port: port})
		}
	}
	return result
}

func ProxyToDogStatsD(incomingData []byte) error {
	log.Debug("Proxying to DogStatsD: %s", incomingData)
	return nil
}

func ProxyToStatsD(incomingData []byte) error {
	devolved, err := devolve.Devolve(string(incomingData))
	if err != nil {
		return err
	}

	devolvedBytes := []byte(devolved)
	for _, target := range statsdTargets {
		go proxyToSingleStatsD(devolvedBytes, target)
	}
	return nil
}

func proxyToSingleStatsD(devolved []byte, statsdTarget targetSpec) error {
	log.Debug("Proxying devolved data to StatsD at %+v: %s", statsdTarget, devolved)
	ip, err := statsdTarget.getIP()
	if err != nil {
		log.Warning("Error getting IP from spec %+v: %s", statsdTarget, err)
		return err
	}

	udpAddr := net.UDPAddr{IP: ip, Port: statsdTarget.Port}
	sock, err := net.DialUDP("udp4", nil, &udpAddr)
	if err != nil {
		log.Warning("Error opening UDP socket for address %+v: %s", udpAddr, err)
		return err
	}
	defer sock.Close()

	bytesWritten, err := sock.Write(devolved)
	if err != nil {
		log.Warning("Error writing to socket for address %s: %s", udpAddr, err)
		return err
	}
	log.Debug("Wrote %d bytes to %+v", bytesWritten, udpAddr)
	return nil
}
