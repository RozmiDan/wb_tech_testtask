package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type minimalOrder struct {
	OrderUID string `json:"order_uid"`
}

func main() {
	var (
		brokers  = flag.String("brokers", "kafka:9092", "comma-separated list of brokers")
		topic    = flag.String("topic", "orders", "kafka topic")
		dir      = flag.String("dir", "/data", "directory with *.json files")
		minDelay = flag.Duration("min", 150*time.Millisecond, "min delay between messages")
		maxDelay = flag.Duration("max", 800*time.Millisecond, "max delay between messages")
		timeout  = flag.Duration("timeout", 5*time.Second, "write timeout per message")
		retries  = flag.Int("retries", 3, "retries per message")
		shuffle  = flag.Bool("shuffle", true, "shuffle files order")
		repeat   = flag.Bool("repeat", false, "repeat endlessly over the directory")
	)
	flag.Parse()

	if *maxDelay < *minDelay {
		*maxDelay = *minDelay
	}
	files, err := collectJSON(*dir)
	must(err)
	if len(files) == 0 {
		must(errors.New("no *.json files found"))
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(*brokers, ",")...),
		Topic:                  *topic,
		Balancer:               &kafka.Hash{},
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: false,
		BatchTimeout:           50 * time.Millisecond,
	}
	defer w.Close()

	runOnce := func() {
		list := append([]string(nil), files...)
		if *shuffle {
			rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
		}
		for _, path := range list {
			payload, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
				continue
			}
			var o minimalOrder
			if err := json.Unmarshal(payload, &o); err != nil || o.OrderUID == "" {
				fmt.Fprintf(os.Stderr, "invalid json or empty order_uid in %s: %v\n", path, err)
				continue
			}
			msg := kafka.Message{Key: []byte(o.OrderUID), Value: payload}

			ok := false
			for i := 0; i < *retries; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), *timeout)
				err := w.WriteMessages(ctx, msg)
				cancel()
				if err == nil {
					ok = true
					break
				}
				time.Sleep(150 * time.Millisecond)
			}
			if ok {
				fmt.Printf("produced %-24s from %s\n", o.OrderUID, filepath.Base(path))
			} else {
				fmt.Fprintf(os.Stderr, "write failed for %s\n", filepath.Base(path))
			}

			if *maxDelay > 0 {
				d := *minDelay
				if *maxDelay > *minDelay {
					d += time.Duration(rand.Int63n(int64(*maxDelay - *minDelay)))
				}
				time.Sleep(d)
			}
		}
	}

	if *repeat {
		for {
			runOnce()
		}
	} else {
		runOnce()
	}
}

func collectJSON(dir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
