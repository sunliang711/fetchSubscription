package main

import (
	"fetchSubscription/decoder"
	"fetchSubscription/downloader"
	"fetchSubscription/parser"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	subURL := pflag.StringP("url", "u", "", "subscription url")
	outputFile := pflag.StringP("output", "o", "config.json", "output config file")
	startPort := pflag.Int16P("sport", "p", 13000, "start port")
	level := pflag.StringP("level", "l", "warning", "log level: debug, info, warning, error, fatal")
	tmplFile := pflag.StringP("tmpl", "t", "v2ray.tmpl", "template file")
	whiteList := pflag.StringSliceP("white-list", "w", nil, "white list keywords: -w keyword1,keyword2")
	blackList := pflag.StringSliceP("black-list", "b", nil, "black list keywords: -b keyword1 -b keyword2")

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

	if len(*whiteList) > 0 && len(*blackList) > 0 {
		logrus.Fatalf("Can not use white list and black list simultaneously")
	}

	cfg := &parser.FilterConfig{}
	if len(*whiteList) > 0 {
		cfg.Mode = parser.ModeWhiteList
		cfg.Lists = *whiteList
	} else if len(*blackList) > 0 {
		cfg.Mode = parser.ModeBlackList
		cfg.Lists = *blackList
	} else {
		cfg.Mode = parser.ModeNone
	}
	logrus.Debugf("filterConfig: %+v", cfg)

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

	// nodes, err := parser.Parse(decoded, nil, false)
	// if err != nil {
	// 	logrus.Fatalf("parse error: %v", err)
	// }
	// for name, node := range nodes {
	// 	logrus.Infof("name: %v node: %v", name, node)
	// }

	config, err := parser.ParseMulti(decoded, cfg, int(*startPort), *tmplFile)
	if err != nil {
		fmt.Printf("ParseMulti error: %v", err)
		return
	}
	// fmt.Printf("config: %v", config)
	logrus.Debugf("config: %v", config)

	f, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("open file error: %v", err)
		return
	}
	defer f.Close()
	f.WriteString(config)
	logrus.Infof("config file has written to '%v'", *outputFile)

}
