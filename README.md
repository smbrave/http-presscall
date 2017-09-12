#http压测工具

```bash
./http-presscall -file=./requst.dat -host="http://10.95.116.47:8010/es/search/all" -qps=10 -worker=100

```

* worker:并发请求的协程数
* qps:每秒请求量
* file:压测的数据一行一个循环使用
* method:采用post或get请求
* host:压测的服务器完整url