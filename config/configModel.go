package config

import (
	"fmt"
	"github.com/smm-goddess/pressure-test/library/text"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	Task   `json:"task"`
	Target `json:"target"`
}

func (config *Config) Legal() bool {
	config.RequestMethod = strings.ToUpper(config.RequestMethod)
	if config.RequestMethod != "GET" && config.RequestMethod != "POST" {
		return false
	}

	if config.TotalGoRoutine < 10 {
		config.TotalGoRoutine = 10
	}
	if count := config.TotalTestCount % config.TotalGoRoutine; count != 0 {
		config.TotalTestCount = count * config.TotalGoRoutine
	}
	return true
}

type Task struct {
	TotalGoRoutine int `json:"totalGoRoutine"`
	TotalTestCount int `json:"totalTestCount"`
	Timeout        int `json:"timeout"`
}

type Target struct {
	TargetLink    string            `json:"targetLink"`
	RequestMethod string            `json:"requestMethod"`
	Headers       map[string]string `json:"headers"`
	Cookies       string            `json:"cookies"`
	Body          string            `json:"body"`
	Query         map[string]string `json:"query"`
	Replace       []Replace         `json:"replace"`
	cookieMap     map[string]string
	body          []byte
}

type Replace struct {
	Origin   string `json:"origin"`
	Format   string `json:"format"`
	Params   string `json:"params"`
	Location string `json:"location"`
}

func (target *Target) GetCookies() map[string]string {
	if len(target.cookieMap) == 0 && len(target.Cookies) > 0 {
		cookies := strings.Split(target.Cookies, ";")
		cookieMap := make(map[string]string)
		for _, cookie := range cookies {
			cook := strings.Split(cookie, "=")
			if len(cook) >= 2 {
				cookieMap[cook[0]] = cook[1]
			}
		}
		target.cookieMap = cookieMap
	}
	return target.cookieMap
}

func (target *Target) GetBody() []byte {
	if s, err := os.Stat(target.Body); err != nil || s.IsDir() {
		return []byte{}
	} else {
		target.body, _ = ioutil.ReadFile(target.Body)
	}
	return target.body
}

var randSRegex = regexp.MustCompile("{{randS\\((\\d+)\\)}}")
var randSLRegex = regexp.MustCompile("{{randSL\\((\\d+)\\)}}")

func (replace *Replace) Replace(s string, index int) string {
	params := strings.Split(replace.Params, ",")
	formatParams := make([]interface{}, len(params))
	for i, param := range params {
		if param == "{{index}}" {
			formatParams[i] = index
		} else if match := randSRegex.FindStringSubmatch(param); len(match) > 0 {
			length, _ := strconv.Atoi(match[1])
			formatParams[i] = text.GenerateRandomString(length)
		} else if match := randSLRegex.FindStringSubmatch(param); len(match) > 0 {
			length, _ := strconv.Atoi(match[1])
			formatParams[i] = text.GenerateRandomStringLower(length)
		} else {
			formatParams[i] = param
		}
	}
	return strings.Replace(s, replace.Origin, fmt.Sprintf(replace.Format, formatParams...), -1)
}
