package downloader

import (
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
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

	client := http.Client{}
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

	return string(ret), nil
}
