package main

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/juju/ratelimit"
)

//读取压测数据
func Reader(ch chan string) {
	qps := *flag_qps
	data_file := *flag_file
	data, err := ioutil.ReadFile(data_file)
	if err != nil {
		panic(err)
	}

	backet := ratelimit.NewBucket(time.Second/time.Duration(qps), 10)
	lines := strings.Split(string(data), "\n")
	total := len(lines)
	idx := 0
	for {
		line := lines[idx]
		d := backet.Take(1)
		time.Sleep(d)
		ch <- line
		idx += 1
		if idx >= total {
			idx = 0
		}
	}
}
