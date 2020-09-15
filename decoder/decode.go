package decoder

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func Decode(content string) (string, error) {
	//TODO: delete SetLevel
	logrus.SetLevel(logrus.DebugLevel)

	// get length of content
	lenOfContent := len(content)
	logrus.Debugf("original decode content: %v", content)
	logrus.Debugf("lenOfContent: %v", lenOfContent)

	// padding with "=" if lenOfContent %4 !=0
	if lenOfContent%4 != 0 {
		padstring := string("===="[lenOfContent%4:])

		content = fmt.Sprintf("%v%v", content, padstring)
		logrus.Debugf("padding with: %v", padstring)
	}
	// why?
	// replace '-' with '+'
	// replace '_' with '/'
	content = strings.ReplaceAll(content, "-", "+")
	content = strings.ReplaceAll(content, "_", "/")
	logrus.Debugf("decoded decode content: %v", content)

	ret, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		logrus.Errorf("decode error: %v", err)
		return "", err
	}

	return string(ret), nil
}
