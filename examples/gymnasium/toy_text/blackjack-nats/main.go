// -*- compile-command: "go run main.go"; -*-

// blackjack solves the gymnasium Blackjack game with
// Reinforcement Learning over NATS.
package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	agentSubject = "gym.agent.action"
)

func usage() {
	log.Printf("Usage: blackjack")
	flag.PrintDefaults()
}

func showUsageAndExit(exitCode int) {
	usage()
	os.Exit(exitCode)
}

var (
	debug    = flag.Bool("d", false, "Show debug information")
	showHelp = flag.Bool("h", false, "Show help message")
	showTime = flag.Bool("t", false, "Display timestamps")
	subj     = flag.String("sub", "gym.env.*", "NATS subject to listen to for env events")
	urls     = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by commas)")
)

func main() {
	log.SetFlags(0) // disable output of timestamp from log.* functions.
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	opts := setupConnOptions("NATS GEP Blackjack")
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	agent := &agentT{nc: nc}
	nc.Subscribe(*subj, agent.envHandler)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%v]", *subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()
}

type agentT struct {
	nc       *nats.Conn
	msgCount int
}

func (a *agentT) envHandler(msg *nats.Msg) {
	a.msgCount++
	if *debug {
		printMsg(msg, a.msgCount)
	}
	switch msg.Subject {
	case "gym.env.reset":
	case "gym.env.obs":
		a.nc.Publish(agentSubject, []byte("1"))
	case "gym.env.update":
	case "gym.env.decay_epsilon":
	default:
		log.Fatalf("unknown NATS subject [%v]", msg.Subject)
	}
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%v] Received on [%v]: '%v'", i, m.Subject, string(m.Data))
}

func setupConnOptions(name string) []nats.Option {
	totalWait := 1 * time.Minute
	reconnectDelay := time.Second

	return []nats.Option{
		nats.Name(name),
		nats.ReconnectWait(reconnectDelay),
		nats.MaxReconnects(int(totalWait / reconnectDelay)),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("Disconnected due to:%v, will attempt reconnects for %.0fm", err, totalWait.Minutes())
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Reconnected [%v]", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Fatalf("Exiting: %v", nc.LastError())
		}),
	}
}
