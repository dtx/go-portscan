Port scanning for Go with support for concurrency.

Much faster than nmap for basic port checking, currently
only supports TCP networks. But can be easily modified to check
UDP, IPV6 and other networks.

There is a similar tool https://github.com/Sinute/golang-portScan
that I tried to use/debug unsuccessfuly. So this one supports YAML
for configuration and better concepts for concurrency, IP validation 
and calculation.
