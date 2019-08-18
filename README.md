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
    Completed 10000 requests
    Completed 20000 requests
    Completed 30000 requests
    Completed 40000 requests
    Completed 50000 requests
    Completed 60000 requests
    Completed 70000 requests
    Completed 80000 requests
    Completed 90000 requests
    Completed 100000 requests
    Finished 100000 requests
    
    
    Server Software:        
    Server Hostname:        localhost
    Server Port:            8080
    
    Document Path:          /api/v1/write
    Document Length:        0 bytes
    
    Concurrency Level:      100
    Time taken for tests:   9.977 seconds
    Complete requests:      100000
    Failed requests:        0
    Total transferred:      7500000 bytes
    Total body sent:        18700000
    HTML transferred:       0 bytes
    Requests per second:    10023.28 [#/sec] (mean)
    Time per request:       9.977 [ms] (mean)
    Time per request:       0.100 [ms] (mean, across all concurrent requests)
    Transfer rate:          734.13 [Kbytes/sec] received
                            1830.42 kb/s sent
                            2564.55 kb/s total
    
    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0    4   1.0      4      16
    Processing:     0    6   1.8      6      46
    Waiting:        0    5   1.8      4      44
    Total:          0   10   2.0     10      50
    
    Percentage of the requests served within a certain time (ms)
      50%     10
      66%     10
      75%     11
      80%     11
      90%     12
      95%     13
      98%     15
      99%     17
     100%     50 (longest request)
