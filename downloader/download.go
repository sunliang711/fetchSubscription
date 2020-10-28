package downloader

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// unit: ms
	default_download_timeout  = 5000
	max_download_content_size = 50
)

func Download(subscriptionURL string, headers map[string][]string) (string, error) {
	//TODO add proxy support

	// Note: just support GET
	req, err := http.NewRequest("GET", subscriptionURL, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Method": "GET", "url": subscriptionURL}).Errorf("NewRequest error: %v", err)
		return "", err
	}
	for k, vs := range headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}

	}

	downloadTimeoutEnv := os.Getenv("DOWNLOAD_TIMEOUT")
	downloadTimeout, err := strconv.Atoi(downloadTimeoutEnv)
	if err != nil {
		logrus.Infof("Use default_download_timeout: %v", default_download_timeout)
		downloadTimeout = default_download_timeout
	} else {
		logrus.Infof("use env DOWNLOAD_TIMEOUT: %v", downloadTimeout)
	}

	client := http.Client{Timeout: time.Millisecond * time.Duration(downloadTimeout)}
	logrus.Infof("Downloading url: %v", subscriptionURL)
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Errorf("client.Do error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Errorf("read response error: %v", err)
		return "", nil
	}

	atMost := len(ret)
	if atMost > max_download_content_size {
		atMost = max_download_content_size
	}
	logrus.Infof("Downloaded content: %v...", string(ret)[:atMost])

	return string(ret), nil
}
