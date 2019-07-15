package main

import (
	"bytes"
	"fmt"
	"github.com/smm-goddess/pressure-test/config"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	wg  sync.WaitGroup
	cfg config.Config
)

func gatherResult(ch chan bool) {
	totalCount := 0
	success := 0
	fail := 0
	for totalCount < cfg.TotalTestCount {
		if <-ch {
			success++
		} else {
			fail++
		}
		totalCount++
		fmt.Println("total", totalCount)
	}
	fmt.Println("success:", success)
	fmt.Println("fail", fail)
	wg.Done()
}

func runTest(ch chan int, outputChan chan bool) {
	gid := getGID()
	var client *http.Client
	requestMethod := strings.ToUpper(cfg.RequestMethod)
	if cfg.Timeout > 5 {
		client = &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		}
	} else {
		client = http.DefaultClient
	}

	var uriReplacement []config.Replace   //路径替换
	var bodyReplacement []config.Replace  //body替换
	var queryReplacement []config.Replace //query替换

	for _, replace := range cfg.Replace {
		switch replace.Location {
		case "body":
			bodyReplacement = append(bodyReplacement, replace)
		case "query":
			queryReplacement = append(queryReplacement, replace)
		case "uri":
			uriReplacement = append(uriReplacement, replace)
		}
	}

	for {
		index := <-ch
		target := cfg.TargetLink
		var data string
		data = string(cfg.GetBody())
		if len(data) > 0 {
			for _, replace := range bodyReplacement {
				data = replace.Replace(data, index)
			}
		}
		for _, replace := range uriReplacement {
			target = replace.Replace(target, index)
		}
		var reader *strings.Reader
		if len(data) > 0 {
			reader = strings.NewReader(data)
		} else {
			reader = nil
		}
		if req, err := http.NewRequest(requestMethod, target, reader); err == nil {
			for key, value := range cfg.Headers {
				req.Header.Add(key, value)
			}
			q := req.URL.Query()
			for key, value := range cfg.Query {
				for _, replace := range queryReplacement {
					value = replace.Replace(value, index)
				}
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			for key, value := range cfg.GetCookies() {
				req.AddCookie(&http.Cookie{
					Name:     key,
					Value:    value,
					HttpOnly: true,
				})
			}
			if resp, err := client.Do(req); err == nil {
				if out, err := ioutil.ReadAll(resp.Body); err == nil {
					outputChan <- true
					fmt.Println(gid, ",success", string(out))
				} else {
					outputChan <- false
					fmt.Println(gid, ",response", err.Error())
				}
				_ = resp.Body.Close()
			} else {
				outputChan <- false
				fmt.Println(gid, ",request", err.Error())
			}
		} else {
			outputChan <- false
			fmt.Println(gid, ",make request", err.Error())
		}
	}
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func main() {
	err := config.LoadConfig("config.json", &cfg)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if !cfg.Legal() {
		fmt.Println("config file is illegal")
		os.Exit(1)
	}

	input := make(chan int, cfg.TotalGoRoutine)
	output := make(chan bool, cfg.TotalGoRoutine)
	wg.Add(1)
	go gatherResult(output)
	for i := 0; i < cfg.TotalGoRoutine; i++ {
		go runTest(input, output)
	}
	for i := 0; i < cfg.TotalTestCount; i++ {
		input <- i
	}
	wg.Wait()
}
