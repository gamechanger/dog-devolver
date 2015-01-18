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

func parseAddress(host string, portString string, lookup bool) (net.IP, int, error) {
	ip := net.ParseIP(host)
	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, 0, err
	}

	if ip == nil && lookup == true {
		ips, err := net.LookupIP(host)
		ip = ips[0]
		if err != nil {
			return nil, 0, err
		}
	}

	return ip, port, nil
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

// Take the StatsD hosts and ports from config and merge them
// into a slice of targetSpecs
func zipTargets(hosts, ports []string) []targetSpec {
	var result []targetSpec
	for idx, host := range hosts {
		if len(ports) <= idx {
			return result
		}

		ip, port, err := parseAddress(host, ports[idx], true)
		if err != nil {
			panic(fmt.Sprintf("Invalid IP (%s) or port (%s) spec, got error: %s", host, ports[idx], err))
		}

		if ip == nil {
			result = append(result, targetSpec{Host: host, Port: port})
		} else {
			result = append(result, targetSpec{IP: ip, Port: port})
		}
	}
	return result
}

func ProxyToDogStatsD(msg string) error {
	log.Debug("Proxying unaltered data to DogStatsD: %s", msg)
	ip, port, err := parseAddress(config.DOGSTATSD_HOST, config.DOGSTATSD_PORT, true)
	if err != nil {
		log.Warning("Error in resolving DogStatsD address: %s", err)
		return err
	}

	udpAddr := net.UDPAddr{IP: ip, Port: port}
	sock, err := net.DialUDP("udp4", nil, &udpAddr)
	if err != nil {
		log.Warning("Error opening UDP socket for address %+v: %s", udpAddr, err)
		return err
	}
	defer sock.Close()

	bytesWritten, err := sock.Write([]byte(msg))
	if err != nil {
		log.Warning("Error writing to DogStatsD socket: %s", err)
		return err
	}

	log.Debug("Successfully wrote %d bytes to DogStatsD", bytesWritten)
	return nil
}

func ProxyToStatsD(msg string) error {
	devolved, err := devolve.Devolve(msg)
	if err != nil {
		log.Warning("Error in devolving message %s. Error was: %s", msg, err)
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
	log.Debug("Successfully wrote %d bytes to %+v", bytesWritten, udpAddr)
	return nil
}
