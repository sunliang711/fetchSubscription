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

	config, err := parser.ParseMulti(decoded, nil, int(*startPort))
	if err != nil {
		fmt.Printf("ParseMulti error: %v", err)
		return
	}
	fmt.Printf("config: %v", config)

	f, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("open file error: %v", err)
		return
	}
	defer f.Close()
	f.WriteString(config)

}
