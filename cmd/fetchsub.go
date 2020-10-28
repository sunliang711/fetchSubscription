package main

import (
	"fetchSubscription/decoder"
	"fetchSubscription/downloader"
	"fetchSubscription/parser"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	SEP_LIST      = "|"
	SEP_KEY_VALUE = ":"
	SEP_ITEM      = ","
)

func main() {
	subURL := pflag.StringP("url", "u", "", "subscription url")
	outputFile := pflag.StringP("output", "o", "config.json", "output config file")
	startPort := pflag.Int16P("sport", "p", 13000, "start port")
	level := pflag.StringP("level", "l", "warning", "log level: debug, info, warning, error, fatal")
	tmplFile := pflag.StringP("tmpl", "t", "v2ray.tmpl", "template file")
	filterList := pflag.StringP("filter", "f", "", "unix-like pipe filter; specify filter list,format: 'w:VIP2,VIP3|b:game,tv'; 'w' for white list, 'b' for black list. filter rule execute one by one")

	pflag.Parse()

	switch *level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}

	logrus.Infof("sub url: %v", *subURL)
	logrus.Infof("output config file: %v", *outputFile)
	logrus.Infof("start port: %v", *startPort)
	logrus.Infof("template file: %v", *tmplFile)
	logrus.Infof("filter list: %v", *filterList)

	var cfgs []*parser.FilterConfig
	if len(*filterList) > 0 {
		// split by '|'
		lists := strings.Split(*filterList, SEP_LIST)
		for _, list := range lists {
			if len(list) == 0 {
				continue
			}
			logrus.Debugf("list item: %v", list)
			// split by ':'
			items := strings.Split(list, SEP_KEY_VALUE)
			if len(items) != 2 {
				logrus.Fatalf("filter list format error")
			}
			cfg := &parser.FilterConfig{}
			switch items[0] {
			case "b":
				cfg.Mode = parser.ModeBlackList
			case "w":
				cfg.Mode = parser.ModeWhiteList
			default:
				logrus.Fatalf("filter list format error: unknow filter type")
			}
			cfg.Lists = strings.Split(items[1], SEP_ITEM)
			cfgs = append(cfgs, cfg)
		}
	}

	logrus.Debugf("filterConfig: %+v", cfgs)

	if *subURL == "" {
		logrus.Fatalf("no sub subscription url")
	}
	content, err := downloader.Download(*subURL, nil)
	if err != nil {
		logrus.Fatalf("download error: %v", err)
	}
	decoded, err := decoder.Decode(content)
	if err != nil {
		logrus.Fatalf("decode error: %v", err)
	}

	config, err := parser.ParseMultiV2ray(decoded, cfgs, int(*startPort), *tmplFile)
	if err != nil {
		logrus.Fatalf("ParseMulti error: %v", err)
	}
	logrus.Debugf("config: %v", config)

	f, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Fatalf("open file error: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(config)
	if err != nil {
		logrus.Fatalf("Write file: %v error: %v", *outputFile, err)
	}
	logrus.Infof("config file has written to '%v'", *outputFile)

}
