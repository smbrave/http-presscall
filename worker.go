package main

import (
	"io/ioutil"
	"strings"
	"time"

	"log"

	"github.com/juju/ratelimit"
	"github.com/toolkits/net/httplib"
)

//测试结果
type result_t struct {
	Err  error
	Cost time.Duration
}

//压测工作
func Worker(host string, method string, data chan string, result chan *result_t) {

	for {
		select {
		case <-data:
			start := time.Now()
			var client *httplib.BeegoHttpRequest = nil
			if strings.ToUpper(method) == "POST" {
				client = httplib.Post(host)
			} else {
				client = httplib.Get(host)
			}
			_, err := client.Bytes()
			result <- &result_t{Err: err, Cost: time.Now().Sub(start)}
		}
	}
}

//读取压测数据
func Reader(ch chan string, qps int, data_file string) {
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

//统计结果
func Result(result chan *result_t) {
	var succ_num int64
	var fail_num int64
	var cost int64

	var pre_succ_num int64
	var pre_fail_num int64
	var pre_cost int64

	var prestat = time.Now()

	for {
		select {
		case res := <-result:
			if res.Err != nil {
				fail_num += 1
				continue
			}

			succ_num += 1
			cost += res.Cost.Nanoseconds()
			cur := time.Now()

			if cur.Sub(prestat).Seconds() >= 1 {
				prestat = cur
				log.Printf("total[%d] error[%d] avg[%.2f]ms",
					succ_num-pre_succ_num, fail_num-pre_fail_num,
					float64(cost-pre_cost)/(float64(succ_num-pre_succ_num)*float64(1000000)),
				)
				pre_succ_num = succ_num
				pre_fail_num = fail_num
				pre_cost = cost
			}
		}
	}
}
