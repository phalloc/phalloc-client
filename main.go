package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

const (
	Author      = "Tai Groot"
	Version     = "0.0.1"
	ReleaseDate = "2019-06-12"
	url         = "https://api.phalloc.com/event"
)

type Metrics struct {
	Data []Metric `json:"data"`
}
type Metric struct {
	Mac       string `json:"mac"`
	TimeStamp string `json:"ts"`
	RSSI      string `json:"rssi"`
}

var (
	debug = flag.Bool("debug", false, "debug")
	cmd   []string
)

// Parse flags
func init() {
	flag.Parse()
}

func process(jstring string) {
	var arr Metrics
	if err := json.Unmarshal([]byte(jstring), &arr); err == nil {
		b, _ := json.Marshal(arr)
		var jsonStr = []byte(b)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

	} else {
		log.Fatal(err)
	}

}

// main loop
func main() {
	cmd = append(cmd, "/usr/bin/phalloc-sniffer")
	cmd = append(cmd, "-m -t 3 wlan1")
	for {

		var stout, sterr bytes.Buffer
		cmd := exec.Command(cmd[0], cmd[1:]...)

		cmd.Stdout = &stout
		cmd.Stderr = &sterr

		cmd.Run()
		if *debug {
			log.Printf("Stdout: %q\n", stout.String())
			log.Printf("Stderr: %q\n", sterr.String())
		}
		process(stout.String())
	}
}
