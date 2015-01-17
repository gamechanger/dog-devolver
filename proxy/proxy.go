package proxy

import (
	"fmt"
	"net"
	"strconv"

	"github.com/gamechanger/dog-devolver/config"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("dog-devolver")

type targetSpec struct {
	IP   net.IP
	Host string
	Port int
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

func ProxyToStatsD(incomingData []byte, statsdTarget targetSpec) error {
	log.Debug("Proxying devolved data to StatsD at %s: %s", statsdTarget, incomingData)
	return nil
}

var STATSD_TARGETS = zipTargets(config.STATSD_HOSTS, config.STATSD_PORTS)
