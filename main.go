package main

import (
	"flag"
	"github.com/fxfitz/journald-forwarder/journald"
	"github.com/fxfitz/journald-forwarder/loggly"
	"github.com/fxfitz/journald-forwarder/kafka"
	"log"
	"os"
	"runtime"
	"time"
)

var token = flag.String("token", "", "Loggly Token")
var logFile = flag.String("logfile", "/var/log/journald-forwarder.log", "Path to log file to write")
var tag = flag.String("tag", "", "What tag to use on Loggly")
var broker = flag.String("brokers", "", "Kafka Brokers")
var topic = flag.String("topic", "", "Kafka Topic")

func main() {
	flag.Parse()

	if *token == "" {
		log.Fatalf("-token is required.")
	}

	if *brokers == "" {
		log.Fatalf("-brokers is required.")
	}

	if *topic == "" {
		log.Fatalf("-topic is required.")
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	c := make(chan journald.JournalEntry, 2)
	uri := loggly.GenerateUri(*token, *tag)
	go journald.CollectJournal(c)
	go loggly.ProcessJournal(c, uri)
	go kafka.ProcessJournal(c, brokers, topic)

	for {
		time.Sleep(1 * time.Second)
	}
}
