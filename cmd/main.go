package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/boypt/scraper"
	"github.com/jpillora/opts"
)

var VERSION = "0.0.0"

type config struct {
	ConfigFile string `type:"arg" help:"Path to JSON configuration file"`
	Host       string `help:"Listening interface"`
	Port       int    `help:"Listening port"`
	NoLog      bool   `help:"Disable access logs"`
}

func main() {

	h := &scraper.Handler{
		Log:   true,
		Debug: true,
		Headers: map[string]string{
			//we're a trusty browser :)
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
		},
	}

	c := config{
		Host: "0.0.0.0",
		Port: 3000,
	}

	opts.New(&c).
		Repo("github.com/boypt/scraper").
		Version(VERSION).
		Parse()

	h.Log = !c.NoLog

	go func() {
		for {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP)
			<-sig
			if err := h.LoadConfigFile(c.ConfigFile); err != nil {
				log.Printf("[scraper] Failed to load configuration: %s", err)
			} else {
				log.Printf("[scraper] Successfully loaded new configuration")
			}
		}
	}()

	if err := h.LoadConfigFile(c.ConfigFile); err != nil {
		log.Fatal(err)
	}

	log.Printf("[scraper] Listening on %d...", c.Port)
	log.Fatal(http.ListenAndServe(c.Host+":"+strconv.Itoa(c.Port), h))
}
