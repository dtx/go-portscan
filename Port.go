package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v2"
)

var proxyFlag = false
var proxyDialer proxy.Dialer

type ipTables map[string]bool

type IPRange struct {
	Range string `yaml:"range"`
}

type Config struct {
	Concurrency int
	TimeOut     int       `yaml:"timeout"`
	IPRanges    []IPRange `yaml:"ip"`
	Proxy       string    `yaml:"proxy"`
}

var config Config

func main() {

	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	yamlFile, _ := filepath.Abs("./config")
	yamlData, err := ioutil.ReadFile(yamlFile)
	config = Config{}
	err = yaml.Unmarshal(yamlData, &config)

	if err != nil {
		fmt.Printf("error: %v", err)
	}
	if config.Proxy != "" {
		Proxy_setUp(&config)
		proxyFlag = true
	}
	CheckPort(&config)
}

func isValidPortNumber(n int) bool {
	if 0 < n && n <= 65535 {
		return true
	}
	return false
}

func getPortsInRange(rangeStr string) ([]int, error) {
	ports := []int{}
	portRangeStr := strings.Split(rangeStr, "-")
	if len(portRangeStr) == 1 {
		portNum, err := strconv.Atoi(portRangeStr[0])
		if err != nil {
			return ports, errors.New("Port Parse Error")
		}
		if isValidPortNumber(portNum) {
			ports = append(ports, portNum)
		}
	} else {
		portRangeStart, err := strconv.Atoi(portRangeStr[0])
		if err != nil {
			return ports, errors.New("Port Parse Error")
		}
		portRangeEnd, err := strconv.Atoi(portRangeStr[1])
		if err != nil {
			return ports, errors.New("Port Parse Error")
		}
		if portRangeEnd < portRangeStart {
			fmt.Println("Invalid port range")
		}

		for portNum := portRangeStart; portNum <= portRangeEnd; portNum++ {
			if isValidPortNumber(portNum) {
				ports = append(ports, portNum)
			}
		}
	}
	return ports, nil
}

func CheckPort(c *Config) {
	concurrency := 10
	blockers := make(chan bool, concurrency)

	for _, ipRange := range c.IPRanges {
		rangeStr := ipRange.Range

		components := strings.Split(rangeStr, ":")

		if len(components) != 2 {
			fmt.Println("Invalid range string provided:", rangeStr)
		}

		ip, ipnet, err := net.ParseCIDR(components[0])
		if err != nil {
			fmt.Println(err)
		}

		ports, err := getPortsInRange(components[1])
		if err != nil {
			fmt.Println(err)
		}

		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			for _, port := range ports {
				blockers <- true
				go checkTCP(ip.String(), strconv.Itoa(port), blockers, c.TimeOut)
			}
		}
	}

	for i := 0; i < cap(blockers); i++ {
		blockers <- true
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func checkTCP(ip string, port string, blocker chan bool, timeout int) {
	defer func() { <-blocker }()
	if !proxyFlag  {
		connection, err := net.DialTimeout("tcp", ip+":"+port, time.Duration(timeout)*time.Second)
		if err == nil {
			fmt.Printf("%s:%s - true\n", ip, port)
			connection.Close()
		}
	} else {

		connection, err := proxyDialer.Dial("tcp", ip+":"+port)
		if err == nil {
			connection.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			fmt.Printf("%s:%s - true\n", ip, port)
			connection.Close()
		}
	}
}

func Proxy_setUp(config *Config) {

	uri, err := url.Parse(config.Proxy)
	if err != nil {
		panic(err)
	}

	dialer, err := proxy.FromURL(uri, &net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	proxyDialer = dialer
	return

}
