# k8s-test

## The idea behind this application is to test one of the security aspects of k8s.

While the application accepts two parameters "cpu" and "inet"
**cpu** - to emulate a load test for the processor,
**inet** - to test connections to the addresses specified in the inet.cfg file 

```shell
> cat inet.cfg
example.com:443
192.168.1.100:80
```shell
