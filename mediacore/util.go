package mediacore

import (
	"math/rand"
	"net"
	"time"
)

func GetRandomUdpPort() int {
	min := 10000
	max := 65535
	rand.Seed(time.Now().UnixNano())
	port := 0
	for {
		p := rand.Intn(max-min) + min
		ln, err := net.ListenUDP("udp", &net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("127.0.0.1"),
		})
		if err != nil {
			continue
		}
		defer ln.Close()

		port = p

		break
	}

	return port
}
