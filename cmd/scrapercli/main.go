package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/boypt/scraper"
)

var (
	configFile = flag.String("c", "", "config file")
	entry      = flag.String("e", "", "entry point (key of the config)")
	query      = flag.String("q", "", "query string")
	testAll    = flag.Bool("testall", false, "test all keys in the config using the query")
)

func main() {
	flag.Parse()

	h := &scraper.Handler{
		Log:   true,
		Debug: true,
		Headers: map[string]string{
			//we're a trusty browser :)
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
		},
	}

	if err := h.LoadConfigFile(*configFile); err != nil {
		log.Fatalln(*configFile, err)
	}

	param := map[string]string{}
	m, err := url.ParseQuery(*query)
	if err != nil {
		log.Fatalln(err)
	}

	for k, v := range m {
		param[k] = v[0]
	}

	if *testAll {
		for name, endpoint := range h.Config {
			if strings.Contains(name, "/") {
				log.Println("skip entpoint:", name)
				continue
			}
			result, err := endpoint.Execute(param)
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			log.Printf("endpoint %s returned %d results\n", name, len(result))
		}

		return
	}

	enpoint, ok := h.Config[*entry]
	if !ok {
		log.Fatal("entry pont not found")
	}
	result, err := enpoint.Execute(param)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}
