package main

import (
	"bytes"
	"net"
	"net/url"
	"strings"
	"time"

	"fmt"

	"github.com/toolkits/net/httplib"
)

//测试结果
type result_t struct {
	Err  error
	Cost time.Duration
}

func getConn(reqUrl string) (net.Conn, error) {
	u, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	port := u.Port()
	if port == "" || port == "0" {
		port = "80"
	}

	//fmt.Printf("host:%s port:%s path:%s param:%s\n", u.Host, port, u.Path, u.RawQuery)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", u.Host) //获取一个TCP地址信息,TCPAddr
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr) //创建一个TCP连接:TCPConn
	if err != nil {
		panic(err)
	}
	return conn, nil
}

func buildRequest(reqUrl, method, body string, header map[string]string) ([]byte, error) {
	u, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if strings.ToUpper(method) == "GET" {
		uri := u.Path
		if u.RawQuery != "" {
			uri += "&" + u.RawQuery
		}
		buf.WriteString(fmt.Sprintf("GET %s HTTP/1.1\r\n", uri))
	} else if strings.ToUpper(method) == "POST" {
		buf.WriteString(fmt.Sprintf("POST %s HTTP/1.1\r\n", u.Path))
	} else {
		return nil, fmt.Errorf("method:%s not support", method)
	}

	buf.WriteString(fmt.Sprintf("Host: %s\r\n", u.Hostname()))

	if header != nil {
		for k, v := range header {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}

	if strings.ToUpper(method) == "POST" {
		buf.WriteString("\r\n")
		buf.WriteString(body)
	}
	return buf.Bytes(), nil
}

//压测工作
func Worker(data chan string, result chan *result_t) {
	reqUrl := *flag_url
	method := *flag_method
	timeout := *flag_timeout

	for {
		select {
		case line := <-data:
			start := time.Now()

			if strings.ToUpper(method) == "POST" {
				var client *httplib.BeegoHttpRequest = nil
				client = httplib.Post(reqUrl)
				client.Header("Content-Type", "application/x-www-form-urlencoded")
				client.Body(line)
				client.SetTimeout(time.Duration(timeout)*time.Second, time.Duration(timeout)*time.Second)
				_, err := client.Bytes()
				result <- &result_t{Err: err, Cost: time.Now().Sub(start)}
			} else {
				var realUrl string
				if line != "" {
					realUrl = fmt.Sprintf("%s%s", reqUrl, line)
				} else {
					realUrl = reqUrl
				}
				client := httplib.Get(realUrl)
				_, err := client.Bytes()
				result <- &result_t{Err: err, Cost: time.Now().Sub(start)}
			}

		}
	}
}
