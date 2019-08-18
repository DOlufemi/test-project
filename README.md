[![Build Status](https://travis-ci.com/iavael/test-project.svg?branch=master)](https://travis-ci.com/iavael/test-project)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/32fa5514075a46d2b8e0605992baa3bd)](https://www.codacy.com/app/iavael/test-project?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=iavael/test-project&amp;utm_campaign=Badge_Grade)
[![GolangCI](https://golangci.com/badges/github.com/iavael/test-project.svg)](https://golangci.com/r/github.com/iavael/test-project)

Input format
------------
POST data in specified format at URL /api/v1/write
```json
{"ts": 1000000, "key": "metric1", "val": 123}
```
  * ts — timestamp (unix time)
  * key — metric name
  * val — metric value

Stored format
-------------
Stored format was inspired by graphite metrics protocol

    {timestamp} {name} {value}

Install & Run
-------
```bash
# https://hub.docker.com/r/iavael/test-project
docker run iavael/test-project
```

Run tests
---------
```bash
go test .
```

Run benchmark
-------------
```bash
go test -bench .
```

Run loadtest
------------
```bash
TMPFILE=$(mktemp) && echo "${TMPFILE}" 
echo '{"ts": 1000000, "key": "metric1", "val": 123}' > "${TMPFILE}" 
ab -n 1000000 -c 100 -p "${TMPFILE}" http://localhost:8080/api/v1/write
```

Performance
-----------
Thousands rps via loopback with HDD as backing store
On my computer service shows this performance with default settings

    $ ab -n 1000000 -c 100 -p /tmp/test http://localhost:8080/api/v1/write
    This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
    Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
    Licensed to The Apache Software Foundation, http://www.apache.org/
    
    Benchmarking localhost (be patient)
    Completed 100000 requests
    Completed 200000 requests
    Completed 300000 requests
    Completed 400000 requests
    Completed 500000 requests
    Completed 600000 requests
    Completed 700000 requests
    Completed 800000 requests
    Completed 900000 requests
    Completed 1000000 requests
    Finished 1000000 requests
    
    
    Server Software:        
    Server Hostname:        localhost
    Server Port:            8080
    
    Document Path:          /api/v1/write
    Document Length:        0 bytes
    
    Concurrency Level:      100
    Time taken for tests:   110.640 seconds
    Complete requests:      1000000
    Failed requests:        0
    Total transferred:      75000000 bytes
    Total body sent:        187000000
    HTML transferred:       0 bytes
    Requests per second:    9038.31 [#/sec] (mean)
    Time per request:       11.064 [ms] (mean)
    Time per request:       0.111 [ms] (mean, across all concurrent requests)
    Transfer rate:          661.99 [Kbytes/sec] received
                            1650.55 kb/s sent
                            2312.54 kb/s total
    
    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0    4   1.8      4      41
    Processing:     0    7   4.7      6     251
    Waiting:        0    5   4.5      5     249
    Total:          0   11   5.1     10     256
    
    Percentage of the requests served within a certain time (ms)
      50%     10
      66%     11
      75%     11
      80%     12
      90%     14
      95%     17
      98%     23
      99%     30
     100%    256 (longest request)
