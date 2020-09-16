package parser

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type FilterConfig struct {
	// 是否启用黑名单，不启用则用白名单
	EnableBlackList bool
	// 黑名单或者白名单匹配时的正则
	Regex string
}

func Parse(nodesContent string, cfg *FilterConfig, full bool) (map[string]string, error) {
	//TODO delete
	logrus.SetLevel(logrus.DebugLevel)

	// name => config
	ret := make(map[string]string)

	reader := strings.NewReader(nodesContent)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		node := scanner.Text()
		if node == "" {
			continue
		}
		logrus.Debugf("node data: %v", node)
		name, parsed, err := parse(node, full)
		if err != nil {
			logrus.Errorf("parse node:%v error: %v,skip...", node, err)
		} else {
			if filter(name, cfg) {
				ret[name] = parsed
			}
		}
	}
	return ret, nil
}

const (
	PrefixVmess = "vmess://"
	PrefixSs    = "ss://"
)

func parse(node string, full bool) (string, string, error) {
	switch {
	case strings.HasPrefix(node, PrefixVmess):
		return parse_vmess(node, full)
	case strings.HasPrefix(node, PrefixSs):
		return parse_ss(node, full)
	default:
		return "", "", fmt.Errorf("Only support vmess ss")
	}
}

func filter(name string, cfg *FilterConfig) bool {
	//TODO
	return true
}
