package main

import (
	"flag"
)

var (
	flag_method  = flag.String("method", "get", "http method get/post")
	flag_qps     = flag.Int("qps", 10, "qps for http presscall")
	flag_file    = flag.String("file", "", "file for presscall data")
	flag_worker  = flag.Int("worker", 1000, "worker for http request")
	flag_url    = flag.String("url", "", "url for http request")
	flag_timeout = flag.Int("timeout", 1, "timeout for http")
)

func main() {
	flag.Parse()
	data := make(chan string, *flag_worker)
	result := make(chan *result_t, *flag_worker)

	go Reader(data)
	go Result(result)

	for i := 0; i < *flag_worker; i++ {
		go Worker(data, result)
	}

	select {}

}
