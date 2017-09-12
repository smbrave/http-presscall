package main

import (
	"flag"
)

var (
	flag_method = flag.String("method", "get", "http method get/post")
	flag_qps    = flag.Int("qps", 10, "qps for http presscall")
	flag_file   = flag.String("file", "", "file for presscall data")
	flag_worker = flag.Int("worker", 100, "worker for http request")
	flag_host   = flag.String("host", "", "host for http request")
)

func main() {
	flag.Parse()
	data := make(chan string, *flag_worker)
	result := make(chan *result_t, *flag_worker)

	go Reader(data, *flag_qps, *flag_file)
	go Result(result)

	for i := 0; i < *flag_worker; i++ {
		go Worker(*flag_host, *flag_method, data, result)
	}

	select {}

}
