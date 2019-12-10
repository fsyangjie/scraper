package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/boypt/scraper"
	"github.com/jpillora/opts"
)

var VERSION = "0.0.0"

type config struct {
	ConfigFile string `opts:"help=Path to JSON configuration file"`
	EntryPoint string `opts:"help=entrypoint"`
	Query      string `opts:"help=query"`
	Item       string `opts:"help=item"`
	TestAll    bool   `opts:"help=all"`
}

func testAll(h *scraper.Handler, c *config) {

	for k, endpoint := range h.Config {
		if strings.Contains(k, "item") {
			log.Println("skiped ", k)
			continue
		}

		log.Println("^^^^^^^^^^^^^  testing for ", k)

		param := map[string]string{"query": c.Query}
		result, err := endpoint.Execute(param)
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}

		log.Println("^^^^^^^^^^^^^  result size: ", len(result))
	}
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

	c := config{}
	opts.New(&c).
		Repo("github.com/boypt/scraper").
		Version(VERSION).
		Parse()

	if err := h.LoadConfigFile(c.ConfigFile); err != nil {
		log.Fatal(err)
	}

	if c.TestAll {
		testAll(h, &c)
		return
	}

	enpoint, ok := h.Config[c.EntryPoint]
	if !ok {
		log.Fatal("entry pont not found")
	}
	param := map[string]string{}

	if c.Query != "" {
		param["query"] = c.Query
	}
	if c.Item != "" {
		param["item"] = c.Item
	}

	result, err := enpoint.Execute(param)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}
