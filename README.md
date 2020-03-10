# conntest
Connection testing tool

**This is in very early alpha, that's why the "versioning" is the commit hash used to generate the image**

## What is conntest?
Currently most networking testing tools focus on throughput, and the related packetloss (iPerf3 for example.) We discovered that what we really needed for our kubernetes clusters was a tool that would send low amounts of traffic, report on the path MTU (PMTU) and round time trip (RTT) time. This and some other metrics can be scraped by prometheus and used to build dashboards showing that all nodes can talk to each other, and that there are no underlying connectivity issues.
This was originally started as an intern project by [David Dai](https://github.com/daianzhuo) during his internship at Thought Machine.

## What protocols are currently supported
We currently only support TCP as the protocol, but hope to add HTTP and GRPC later. This is so that we can check the effects of protocol aware CNIs (Istio for example.)

## How to get started
TODO

```Usage:
conntest [OPTIONS]

Application Options:
      --host_port=          Port to host on (default: 8080)
      --dst_hst=            Destination host to target for tests (default:
                            127.0.0.1:8080)
      --wait_time=          Minimum time between individual tests (default: 5)
      --rand_secs=          Maximum random time to be added to TimeBetTests
                            (default: 5.0)
      --short_test_bytes=   Bytes to use for short tests (default: 10)
      --long_test_bytes=    Bytes to use for long tests (default: 10000)
      --times_to_send=      Number of times to send bytes (default: 0)
      --DNS_retry_interval= Time between attempts to re-discover SRV records
                            (default: 5.0)
      --max_DNS_retries=    Maximum number of retries when attmpting to
                            re-discover SRV records, use -1 for infinite
                            retries (default: -1)

Help Options:
  -h, --help                Show this help message
```

## How to build
```plz build //...```

If you'd prefer to use docker to build and run all tests use

```docker build -f Dockerfile-builder . --rm=false```
