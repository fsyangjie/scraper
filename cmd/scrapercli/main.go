package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/fsyangjie/scraper"
)

var (
	configFile = flag.String("c", "", "config file")
	entry      = flag.String("e", "", "entry point (key of the config)")
	query      = flag.String("q", "", "query string")
	testAll    = flag.Bool("testall", false, "test all keys in the config using the query")
	hide       = flag.Bool("hide", false, "hide the scraper debug info")
)

func main() {
	flag.Parse()

	h := &scraper.Handler{
		Log:   !*hide,
		Debug: !*hide,
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
		paral := make(chan struct{}, 5)
		var wg sync.WaitGroup
		for name, endpoint := range h.Config {
			if strings.Contains(name, "/item") {
				log.Println("entpoint skiped (per item entpoint):", name)
				continue
			}

			wg.Add(1)
			go func(n string, e *scraper.Endpoint) {
				paral <- struct{}{}
				result, err := e.Execute(param)
				if err != nil {
					log.Println(n, err)
				}
				log.Printf("endpoint %s returned %d results\n", n, len(result))
				<-paral
				wg.Done()
			}(name, endpoint)
		}

		wg.Wait()
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
