package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/EDDYCJY/fake-useragent"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"net"

	"github.com/apoorvam/goterminal"
)

func goFun(postContent string, Referer string, XforwardFor bool, customIP ipArray, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			go goFun(postContent, Referer, XforwardFor, customIP, wg)
		}
	}()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	for true {
		if customIP != nil && len(customIP) > 0 {
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				rand.Seed(time.Now().Unix())
				ip := customIP[rand.Intn(len(customIP))]
				if strings.HasPrefix(addr, "https") {
					addr = ip + ":443"
				} else if strings.HasPrefix(addr, "http") {
					addr = ip + ":80"
				} else {
					addr = ip + ":80"
				}
				return dialer.DialContext(ctx, network, addr)
			}
			transport.DialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				rand.Seed(time.Now().Unix())
				ip := customIP[rand.Intn(len(customIP))]
				if strings.HasPrefix(addr, "https") {
					addr = ip + ":443"
				} else if strings.HasPrefix(addr, "http") {
					addr = ip + ":80"
				} else {
					addr = ip + ":80"
				}
				return dialer.DialContext(ctx, network, addr)
			}
		}

		var request *http.Request
		var err1 error = nil
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Second * 10,
		}
		if len(postContent) > 0 {
			request, err1 = http.NewRequest("POST", TargetUrl, strings.NewReader(postContent))
		} else {
			request, err1 = http.NewRequest("GET", TargetUrl, nil)
		}
		if err1 != nil {
			continue
		}
		if len(Referer) == 0 {
			Referer = TargetUrl
		}
		request.Header.Add("Cookie", RandStringBytesMaskImpr(12))
		request.Header.Add("User-Agent", browser.Random())
		request.Header.Add("Referer", Referer)
		if XforwardFor {
			randomip := generateRandomIPAddress()
			request.Header.Add("X-Forwarded-For", randomip)
			request.Header.Add("X-Real-IP", randomip)
		}

		if len(headers) > 0 {
			for _, head := range headers {
				headKey := head.key
				headValue := head.value
				if strings.HasPrefix(head.key, "Random") {
					count, convErr := strconv.Atoi(strings.ReplaceAll(head.value, "Random", ""))
					if convErr == nil {
						headKey = RandStringBytesMaskImpr(count)
					}
				}
				if strings.HasPrefix(head.value, "Random") {
					count, convErr := strconv.Atoi(strings.ReplaceAll(head.value, "Random", ""))
					if convErr == nil {
						headValue = RandStringBytesMaskImpr(count)
					}
				}
				request.Header.Del(headKey)
				request.Header.Set(headKey, headValue)
			}
		}

		resp, err2 := client.Do(request)
		if err2 != nil {
			continue
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	wg.Done()
}

var h = flag.Bool("h", false, "this help")
var count = flag.Int("c", 16, "concurrent thread for download,default 16")
var url = flag.String("s", "", "target url")
var postContent = flag.String("p", "", "post content")
var referer = flag.String("r", "", "referer url")
var detectLocation = flag.Bool("d", false, "detect Real link from the Location in http header")
var xforwardfor = flag.Bool("f", true, "randomized X-Forwarded-For and X-Real-IP address")
var subscribe = flag.String("sub", "", "subscribe url")
var TerminalWriter = goterminal.New(os.Stdout)
var customIP ipArray
var headers headersList
var TargetUrl string

func usage() {
	fmt.Fprintf(os.Stderr,
		`webBenchmark version: /0.6
Usage: webBenchmark [-c concurrent] [-s target] [-p] [-r refererUrl] [-f] [-i ip]

Options:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		`
Advanced Example:
webBenchmark -c 16 -s https://some.website -r https://referer.url -i 10.0.0.1 -i 10.0.0.2 
	16 concurrent to benchmark https://some.website with https://referer.url directly to ip 10.0.0.1 and 10.0.0.2
webBenchmark -c 16 -s https://some.website -r https://referer.url
	16 concurrent to benchmark https://some.website with https://referer.url to dns resolved ip address

`)
}

func main() {
	flag.Var(&customIP, "i", "custom ip address for that domain, multiple addresses automatically will be assigned randomly")
	flag.Var(&headers, "H", "custom header")
	//flag.BoolVar(&detectLocation, "d", true, "detect Real link from the Location in http header")
	flag.Usage = usage
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	routines := *count

	if customIP != nil && len(customIP) > 0 && routines < len(customIP) {
		routines = len(customIP)
	}
	// subscribe mode
	if len(*subscribe) > 0 {
		subs := Subscribe(*subscribe)
		if *detectLocation {
			location := GetHttpLocation(subs)
			if len(location) > 0 {
				TargetUrl = location
			} else {
				TargetUrl = subs
			}
		} else {
			TargetUrl = subs
		}
		go subscribeUpdate(*subscribe)
	}
	// local detect location
	if len(*subscribe) == 0 && *detectLocation && len(*url) > 0 {
		location := GetHttpLocation(*url)
		if len(location) > 0 {
			TargetUrl = location
		} else {
			TargetUrl = *url
		}

		go RefreshHttpLocation(*url)
	}

	go showStat()
	var waitgroup sync.WaitGroup
	if routines <= 0 {
		routines = 16
	}

	for i := 0; i < routines; i++ {
		waitgroup.Add(1)
		go goFun(*postContent, *referer, *xforwardfor, customIP, &waitgroup)
	}
	waitgroup.Wait()
	TerminalWriter.Reset()
}
