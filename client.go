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

	bc := make(chan *osc.Bundle)

	go func() {
		var i int
		for {
			i %= 24
			msgs := sineWaveOnMfader(4, 1, 1, 24, i)
			b := osc.NewBundle(time.Now(), msgs...)
			/* b.WriteTo(conn) */

			msleep(50)
			bc <- b
			i += 1
		}
	}()

	go func() {
		var i int
		for {
			i %= 24
			msgs := sineWaveOnMfader(4, 2, 1, 24, i)
			b := osc.NewBundle(time.Now(), msgs...)
			/* b.WriteTo(conn) */

			msleep(30)
			bc <- b
			i += 1
		}
	}()

	for {
		b := <-bc
		b.WriteTo(conn)
	}

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
		sinV := float32(math.Sin((float64(i+offset) * math.Pi * 2) / float64(cnt)))
		m.AppendArgs((sinV + 1) / 2)
		msgs[i] = m
	}

	return msgs
}

func msleep(ms time.Duration) {
	time.Sleep(ms * time.Millisecond)
}
