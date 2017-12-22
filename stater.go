package main

import (
	"log"
	"time"
)

//统计结果
func Result(result chan *result_t) {
	var succ_num int64
	var fail_num int64
	var cost int64

	var pre_succ_num int64
	var pre_fail_num int64
	var pre_cost int64
	var max_cost, min_cost int64

	var prestat = time.Now()
	var preerror = time.Now()
	max_cost = 0
	min_cost = 1e18

	//统计函数
	stat := func(err error) {
		cur := time.Now()
		if cur.Sub(prestat).Seconds() >= 1 {
			prestat = cur
			log.Printf("total[%d] error[%d] max[%.2f] min[%.2f] avg[%.2f]ms",
				succ_num-pre_succ_num, fail_num-pre_fail_num,
				float64(max_cost)/float64(1000000),
				float64(min_cost)/float64(1000000),
				float64(cost-pre_cost)/(float64(succ_num-pre_succ_num)*float64(1000000)),
			)
			pre_succ_num = succ_num
			pre_fail_num = fail_num
			pre_cost = cost
			max_cost = 0
			min_cost = 1e18
		}

		if err != nil {
			if cur.Sub(preerror).Seconds() >= 10 {
				log.Printf("last error:%s", err.Error())
				preerror = cur
			}
		}
	}

	for {
		select {
		case res := <-result:
			if res.Err != nil {
				fail_num += 1
				stat(res.Err)
				continue
			}
			stat(nil)
			succ_num += 1
			tmp := res.Cost.Nanoseconds()
			cost += tmp
			if tmp > max_cost {
				max_cost = tmp
			}
			if tmp < min_cost {
				min_cost = tmp
			}

		}
	}
}
