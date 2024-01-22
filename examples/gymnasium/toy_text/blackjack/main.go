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

	jsoniter "github.com/json-iterator/go"
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

	jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary
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

type ObsT [3]int
type Update struct {
	Obs        ObsT
	Action     int
	Reward     float64
	Terminated bool
	NextObs    ObsT
}

func (a *agentT) envHandler(msg *nats.Msg) {
	a.msgCount++
	if *debug {
		printMsg(msg, a.msgCount)
	}
	switch msg.Subject {
	case "gym.env.reset": // e.g. [gym.env.reset]: '[[11, 4, 0], {}]'
	case "gym.env.obs": // e.g. [gym.env.obs]: '[11, 4, 0]'
		var obs ObsT
		if err := jsoncomp.Unmarshal(msg.Data, &obs); err != nil {
			log.Fatal(err)
		}
		a.nc.Publish(agentSubject, []byte("1"))
	case "gym.env.update": // e.g. [gym.env.update]: '[[11, 4, 0], 1, 0.0, false, [21, 4, 0]]'
		update := parseUpdate(msg.Data)
		log.Printf("update: %+v", update)
	case "gym.env.decay_epsilon": // e.g. [gym.env.decay_epsilon]: ''
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

func parseObs(data any) ObsT {
	var obs ObsT
	tmpObs, ok := data.([]any)
	if !ok {
		log.Fatalf("parseObs: data: %T", data)
	}
	if len(tmpObs) != len(obs) {
		log.Fatalf("parseObs: tmpObs=%+v, want %v items", tmpObs, len(obs))
	}
	for i := 0; i < len(obs); i++ {
		obs[i] = parseInt(tmpObs[i])
	}
	return obs
}

func parseBool(data any) bool {
	tmpBool, ok := data.(bool)
	if !ok {
		log.Fatalf("parseBool: data: %T", data)
	}
	return tmpBool
}

func parseInt(data any) int {
	tmpInt, ok := data.(float64)
	if !ok {
		log.Fatalf("parseInt: data: %T", data)
	}
	return int(tmpInt)
}

func parseFloat64(data any) float64 {
	tmpFloat64, ok := data.(float64)
	if !ok {
		log.Fatalf("parseFloat64: data: %T", data)
	}
	return tmpFloat64
}

func parseUpdate(data []byte) Update {
	var update []any
	if err := jsoncomp.Unmarshal(data, &update); err != nil {
		log.Fatal(err)
	}
	return Update{
		Obs:        parseObs(update[0]),
		Action:     parseInt(update[1]),
		Reward:     parseFloat64(update[2]),
		Terminated: parseBool(update[3]),
		NextObs:    parseObs(update[4]),
	}
}
