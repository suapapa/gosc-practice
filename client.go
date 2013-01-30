package main

import (
	"./gosc"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"time"
)

func connectToOSCServer(addr string) (net.Conn, error) {
	laddr, resolverr := net.ResolveUDPAddr("udp", addr)
	if resolverr != nil {
		return nil, resolverr
	}

	conn, dialerr := net.DialUDP("udp", nil, laddr)
	if dialerr != nil {
		return nil, dialerr
	}

	return conn, nil
}

func runFader(oscc chan *osc.Bundle, faderId int, sleepMs time.Duration) {
	colors := []string{"red", "green", "blue", "yellow", "pupple", "gray", "orange", "brown", "pink"}

	var i int
	for {
		i %= 24
		msgs := sineWaveOnMfader(4, faderId, 1, 24, i)
		if i == 0 {
			colorM := osc.NewMessage("/4/multifader1/12/color",
				colors[rand.Intn(len(colors))])
			msgs = append(msgs, colorM)
		}

		b := osc.NewBundle(time.Now(), msgs...)

		msleep(sleepMs)
		oscc <- b
		i += 1

	}
}

func main() {
	go runOSCServer()

	conn, err := connectToOSCServer(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	time.Sleep(100)

	oscc := make(chan *osc.Bundle)
	go runFader(oscc, 1, 50)
	/* go runFader(oscc, 2, 30) */
	for {
		b := <-oscc
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
