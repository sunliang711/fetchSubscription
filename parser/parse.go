package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

type ListMode int

func (l ListMode) String() string {
	switch l {
	case ModeNone:
		return "ModeNone"
	case ModeBlackList:
		return "ModeBlackList"
	case ModeWhiteList:
		return "ModeWhiteList"
	}
	return "Unknown mode"
}

const (
	ModeNone ListMode = iota
	ModeWhiteList
	ModeBlackList
)

type FilterConfig struct {
	Mode ListMode
	// 黑名单或者白名单匹配时的keyword
	Lists []string
}

var (
	err  error
	tmpl *template.Template
)

func initTmpl(tmplFile string) {
	// tmplFile := "v2ray.tmpl"
	var err error
	tmpl, err = template.ParseFiles(tmplFile)
	if err != nil {
		logrus.Fatalf("parse template file error: %v", err)
	}
}

// Parse生成outbound的map或者含有单个outbound配置文件(full=true)字符串的map
func Parse(nodesContent string, cfgs []*FilterConfig, full bool, tmplFile string) (map[string]string, error) {
	logrus.Infof("Parse...")
	initTmpl(tmplFile)

	// name => config
	ret := make(map[string]string)

	reader := strings.NewReader(nodesContent)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		node := scanner.Text()
		if node == "" {
			continue
		}
		logrus.Infof("node raw data: %v", node)
		name, parsed, err := parse(node, full)
		if err != nil {
			// logrus.Errorf("parse node: %v error: %v,skip...", node, err)
			logrus.WithFields(logrus.Fields{"node": node, "error": err}).Errorf("parse node error")
		} else {
			ret[name] = parsed
		}
	}

	filter(ret, cfgs)
	logrus.Infof("Pased.")
	return ret, nil
}

const (
	MaxPortNo = 65535
)

var (
	ErrPortRange = errors.New("port range not enough")
)

type Multi struct {
	Ps             string
	InPort         string
	InboundString  string
	OutboundString string
}

// PaParseMultiV2ray 生成所有所有节点都放在同一个配置文件中
func ParseMultiV2ray(nodesContent string, cfg []*FilterConfig, startPort int, tmplFile string) (string, error) {
	// full = false,来获取所有outbound的map
	logrus.Infof("ParseMultiV2ray...")
	outbounds, err := Parse(nodesContent, cfg, false, tmplFile)
	if err != nil {
		return "", err
	}

	if MaxPortNo-startPort < len(outbounds) {
		logrus.WithFields(logrus.Fields{"startPort": startPort, "outbounds len": len(outbounds)}).Errorf(ErrPortRange.Error())
		return "", ErrPortRange
	}

	multiObjs := []Multi{}

	var b bytes.Buffer
	inPort := startPort
	// 根据outbound来生成同等数量的inbound，这些inboud从startPort开始，每次累加1
	// 并且没有被使用(listen),如果被使用了，则用下一个
	for ps, outbound := range outbounds {
		for portInUse(inPort) {
			inPort += 1
			if inPort > MaxPortNo {
				return "", ErrPortRange
			}
		}
		err = tmpl.ExecuteTemplate(&b, "inbound", map[string]string{"ps": ps, "port": fmt.Sprintf("%d", inPort)})
		if err != nil {
			return "", err
		}
		inboundString := b.String()
		ioutil.ReadAll(&b)
		logrus.Debugf("inbound: %v", inboundString)

		multi := Multi{Ps: ps, InPort: fmt.Sprintf("%d", inPort), InboundString: inboundString, OutboundString: outbound}

		multiObjs = append(multiObjs, multi)
		inPort += 1
	}

	err = tmpl.ExecuteTemplate(&b, "multi-outbounds", multiObjs)
	if err != nil {
		return "", err
	}
	logrus.Infof("ParseMultiV2ray done.")
	return b.String(), nil
}

func portInUse(port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Debugf("port: %d not in use", port)
		return false
	}
	conn.Close()
	logrus.Debugf("port: %d in use", port)
	return true
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

func filter(nodes map[string]string, cfgs []*FilterConfig) {
	for _, cfg := range cfgs {
		switch cfg.Mode {
		case ModeBlackList:
			for ps, _ := range nodes {
				for _, black := range cfg.Lists {
					if strings.Contains(ps, black) {
						logrus.Infof("delete black list item: %s", ps)
						delete(nodes, ps)
					}
				}
			}
		case ModeWhiteList:
			for ps, _ := range nodes {
				exist := false
			INNER:
				for _, white := range cfg.Lists {
					if strings.Contains(ps, white) {
						exist = true
						break INNER
					}
				}
				if !exist {
					logrus.Infof("delete non white list item: %s", ps)
					delete(nodes, ps)
				}
			}
		default:
		}
	}
}
