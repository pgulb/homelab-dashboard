package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	yaml "github.com/goccy/go-yaml"
)

type Entry struct {
	Name string
	Lan  string
	Vpn  string
}

type Config struct {
	Entries []Entry
}

func getConfig() Config {
	confFile := os.Getenv("DASHBOARD_CONFIG")
	if confFile == "" {
		confFile = "./config.yml"
	}
	conf, err := os.ReadFile(confFile)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	if err := yaml.Unmarshal(conf, &config); err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", config)
	return config
}

func exitAfterWhile() {
	time.Sleep(time.Millisecond)
	os.Exit(0)
}

func main() {
	log.Println("Starting dashboard...")
	config := getConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := `
		<style>
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body {
			background: #0a0a0a;
			color: #33ff33;
			font-family: "Courier New", Courier, monospace;
			display: flex;
			justify-content: center;
			align-items: center;
			min-height: 100vh;
			padding: 16px;
		}
		table {
			border-collapse: collapse;
			width: 100%;
			max-width: 900px;
			font-size: clamp(14px, 2.5vw, 22px);
			line-height: 1.5;
			text-shadow: 0 0 6px #33ff3380;
		}
		th, td {
			border: 1px solid #33ff3366;
			padding: 10px 14px;
			text-align: left;
		}
		th {
			border-bottom: 2px solid #33ff33;
			font-weight: bold;
			letter-spacing: 1px;
			text-transform: uppercase;
			font-size: 0.8em;
		}
		tr:hover td {
			background: #33ff3310;
		}
		a {
			color: inherit;
			text-decoration: none;
			display: block;
		}
		a:hover {
			text-decoration: underline;
			text-shadow: 0 0 10px #33ff33;
		}
		@media (max-width: 600px) {
			body { padding: 8px; align-items: flex-start; padding-top: 24px; }
			th, td { padding: 8px 10px; }
			th:nth-child(1), td:nth-child(1) { min-width: 0; }
		}
		</style>
		
		<table>
		<tr><th>Service name</th><th>LAN link</th><th>VPN link</th></tr>`
		for _, v := range config.Entries {
			resp = fmt.Sprintf(
				"%s<tr><td>%s</td><td><a href='%s'>%s</a></td><td><a href='%s'>%s</a></td></tr>",
				resp,
				v.Name,
				v.Lan,
				v.Lan,
				v.Vpn,
				v.Vpn,
			)
		}
		resp = fmt.Sprintf(
			"%s%s",
			resp,
			"</table>",
		)
		_, err := w.Write([]byte(resp))
		if err != nil {
			log.Println(err)
		}
	})
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Will exit now!")
		_, err := w.Write([]byte("restarting"))
		if err != nil {
			log.Println(err)
		}
		go exitAfterWhile() // this allows writing response and then exitting
	})

	s := &http.Server{
		Addr:         ":9091",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}