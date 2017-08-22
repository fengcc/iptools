package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

const (
	SSH_CONFIG  = ".ssh/config"
	MY_IP_URL   = "http://64.182.208.181"                          // Google 的 icanhazip.com ip
	IP_FIND_URL = "http://ip.taobao.com/service/getIpInfo.php?ip=" // 淘宝IP地址查询
)

var (
	ipreg = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
)

type TaobaoData struct {
	Country string `json:"country"`
	Area    string `json:"area"`
	Region  string `json:"region"`
	City    string `json:"city"`
	Isp     string `json:"isp"`
}
type TaobaoResp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

func getMySelfIp() {
	resp, err := http.Get(MY_IP_URL)
	if err != nil {
		fmt.Printf("Send to %s error: %s\n", MY_IP_URL, err.Error())
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read response %s error: %s\n", data, err.Error())
		return
	}
	resp.Body.Close()

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
	resp, err := http.Get(IP_FIND_URL + ip)
	if err != nil {
		fmt.Printf("Send to %s error: %s\n", IP_FIND_URL, err.Error())
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read response %s error: %s\n", data, err.Error())
		return
	}
	var taobaoResp TaobaoResp
	if err = json.Unmarshal(data, &taobaoResp); err != nil {
		fmt.Printf("Json unmarshal %s error: %s\n", data, err.Error())
		return
	}
	if taobaoResp.Code == 1 {
		fmt.Printf("Find ip location failed: %s\n", taobaoResp.Data)
		return
	}
	var taobaoData TaobaoData
	if err = json.Unmarshal(taobaoResp.Data, &taobaoData); err != nil {
		fmt.Printf("Json unmarshal %s error: %s\n", taobaoResp.Data, err.Error())
		return
	}

	fmt.Println(taobaoData.Country, taobaoData.Area, taobaoData.Region, taobaoData.City, taobaoData.Isp)
}

func main() {
	flag.Parse()

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
