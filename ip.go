package main

import (
	"flag"
	"fmt"
	"github.com/wangtuanjie/ip17mon"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"github.com/atotto/clipboard"
	"time"
)

const (
	SSH_CONFIG = ".ssh/config"
	MY_IP_URL  = "http://64.182.208.181" // Google 的 icanhazip.com ip
)

var (
	ipreg      = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	ipDataFile = flag.String("ipip", "/Users/feng/Projects/go/src/tools/ip/mydata4vipweek2.dat", "location of mydata4vipweek2.dat")
)

func getMySelfIp() {
	resp, err := http.Get(MY_IP_URL)
	if err != nil {
		fmt.Printf("Send to %s error: %s\n", MY_IP_URL, err.Error())
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read response of %s error: %s\n", MY_IP_URL, err.Error())
		return
	}
	ip := strings.TrimRight(string(data[:]), "\n")
	fmt.Println(ip)
	getLocation(ip)
}

func getIpFromSshConfig(host string) {
	homeDir := os.Getenv("HOME")
	data, err := ioutil.ReadFile(homeDir + "/" + SSH_CONFIG)
	if err != nil {
		fmt.Printf("Read file %s error: %s\n", SSH_CONFIG, err.Error())
		return
	}
	lines := strings.Split(string(data[:]), "\n")

	r := make(map[string]string)
	var key string
	for _, l := range lines {
		if strings.HasPrefix(l, "Host ") {
			key = strings.Split(l, " ")[1]
		} else if strings.HasPrefix(l, "HostName ") {
			r[key] = strings.Split(l, " ")[1]
		}
	}

	if r[host] != "" {
		clipboard.WriteAll(r[host])
		fmt.Println(r[host])
	}
}

func getLocation(ip string) {
	loc, err := ip17mon.Find(ip)
	if err != nil {
		fmt.Printf("ip17mon Find error: %s\n", err.Error())
		return
	}
	fmt.Println(loc.Country, loc.City, loc.Region, loc.Isp)
}

func main() {
	flag.Parse()

	fileInfo, err := os.Stat(*ipDataFile)
	if err != nil {
		fmt.Printf("Get info of %s error: %s\n", *ipDataFile, err.Error())
		return
	}

	// IP 库一个月没有更新了，提示更新
	if time.Now().Sub(fileInfo.ModTime()).Hours() > 24 * 15 {
		fmt.Println("warning: ip library updated before one month ago")
	}

	if err := ip17mon.Init(*ipDataFile); err != nil {
		fmt.Printf("Init ip17mon file %s error: %s\n", *ipDataFile, err.Error())
		return
	}

	switch len(os.Args) {
	case 1:
		getMySelfIp()
	case 2:
		if ipreg.MatchString(os.Args[1]) {
			getLocation(os.Args[1])
		} else {
			getIpFromSshConfig(os.Args[1])
		}
	default:
		fmt.Println("Usage:\n	ip [hostname | ip]")
	}
}
