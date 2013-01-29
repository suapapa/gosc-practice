package main

import (
	"./gosc"
	"fmt"
	"math"
	"net"
	"os"
	"time"
)

func main() {
	// go runOSCServer()

	/* laddr, resolverr := net.ResolveUDPAddr("udp", "192.168.200.21:9000") */
	fmt.Println(os.Args[1])
	laddr, resolverr := net.ResolveUDPAddr("udp", os.Args[1])
	if resolverr != nil {
		fmt.Println("ResolveUDPAddr:", resolverr)
		return
	}

	conn, dialerr := net.DialUDP("udp", nil, laddr)
	if dialerr != nil {
		fmt.Println("DialUDP:", dialerr)
		return
	}
	defer conn.Close()

	/* for i := 0; i < 1; i++ { */
		msgs := sineWaveOnMfader(4, 2, 1, 24, 0)
		b := osc.NewBundle(time.Now(), msgs...)
		b.WriteTo(conn)

		msleep(500)
	/* } */

	// var input string
	// for {
	// 	fmt.Scanln(&input)
	// 	if input == "quit" {
	// 		fmt.Println("quit")
	// 		break
	// 	}
	// }
}

func sineWaveOnMfader(panNo, mfaderNo int, sIdx, cnt int, offset int) []*osc.Message {

	offset %= cnt

	msgs := make([]*osc.Message, cnt)
	for i := sIdx - 1; i < cnt; i++ {
		m := new(osc.Message)
		m.Address = fmt.Sprintf("/%d/multifader%d/%d", panNo, mfaderNo, i+1)
		m.AppendArgs(float32(math.Sin((float64(i)*math.Pi*2)/float64(cnt))))
		msgs[i] = m
	}

	return msgs
}

func msleep(ms time.Duration) {
	time.Sleep(ms * time.Millisecond)
}
