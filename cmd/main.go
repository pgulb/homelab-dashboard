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

func main() {
	log.Println("Starting dashboard...")
	config := getConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := `
		<style>
		th {
			border: 1px dashed gray;
			padding: 5px;
		}
		td {
			border: 1px dashed gray;
			padding: 5px;
		}
		table {
			border: 1px dashed green;
		}
		body {
			background-color: black;
			color: green;
			display: flex;
			justify-content: center;
			align-items: center;
		}
		</style>
		
		<table>
		<tr><th>Service name</th><th>LAN link</th><th>VPN link</th></tr>`
		for _, v := range config.Entries {
			resp = fmt.Sprintf(
				"%s<tr><td>%s</td><td>%s</td><td>%s</td></tr>",
				resp,
				v.Name,
				v.Lan,
				v.Vpn,
			)
		}
		resp = fmt.Sprintf(
			"%s%s",
			resp,
			"</table>",
		)
		w.Write([]byte(resp))
	})

	s := &http.Server{
		Addr:         ":9091",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
