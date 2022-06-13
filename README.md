Port scanning for Go with support for concurrency.

Much faster than nmap for basic port checking, currently
only supports TCP networks. But can be easily modified to check
UDP, IPV6 and other networks.

There is a similar tool https://github.com/Sinute/golang-portScan
that I tried to use/debug unsuccessfuly. So this one supports YAML
for configuration and better concepts for concurrency, IP validation 
and calculation.

Usage:
```
$> go get github.com/dtx/go-portscan
$> cat << EOF > config
--- 
ip: 
  - 
    range: "127.0.0.1/32:20000-50000"
  - 
    range: "8.8.8.8/32:53"
    
proxy: socks5://user:pass@proxy.server
timeout: 5
concurrency: 100
EOF
$> go-portscan
