package main

import (
	"fetchSubscription/decoder"
	"fetchSubscription/downloader"
	"fetchSubscription/parser"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v suburl\n", os.Args[0])
		return
	}

	subURL := os.Args[1]
	content, err := downloader.Download(subURL, nil)
	if err != nil {
		logrus.Fatalf("download error: %v", err)
	}
	decoded, err := decoder.Decode(content)
	if err != nil {
		logrus.Fatalf("decode error: %v", err)
	}

	nodes, err := parser.Parse(decoded, nil, true)
	if err != nil {
		logrus.Fatalf("parse error: %v", err)
	}
	for name, node := range nodes {
		logrus.Infof("name: %v node: %v", name, node)
	}

}
