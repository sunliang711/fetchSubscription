package parser

import (
	"bytes"
	"encoding/json"
	"fetchSubscription/decoder"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

type VmessNode struct {
	Host string
	Path string
	Tls  string
	Add  string
	Port int
	Aid  int
	Net  string
	Type string
	V    string
	Ps   string
	Id   string
	// Class int
}

var (
	err  error
	tmpl *template.Template
)

func init() {
	tmplFile := "v2ray.tmpl"
	var err error
	tmpl, err = template.ParseFiles(tmplFile)
	if err != nil {
		logrus.Fatalf("parse template file error: %v", err)
	}
}

func parse_vmess(node string) (string, string, error) {
	node = node[len(PrefixVmess):]
	logrus.WithField("node", node).Infof("parse_vmess")

	decoded, err := decoder.Decode(node)
	if err != nil {
		logrus.Errorf("decode vmess node error: %v", err)
		return "", "", err
	}
	logrus.Infof("decoded vmess node: %s", decoded)

	// decoded data format:
	// {"host":"","path":"","tls":"","add":"14.17.97.145","port":5010,"aid":2,"net":"tcp","type":"none","v":"2","ps":"台湾 01 [D3/VR/IPLC]","id":"32470e14-85fb-3bf0-aa0c-1f7ba46b58b7","class":3}

	return convert_vmess(decoded)
}

func convert_vmess(node string) (string, string, error) {

	var (
		v2Path string
		v2Host string
	)

	reader := strings.NewReader(node)

	var vnode VmessNode
	err = json.NewDecoder(reader).Decode(&vnode)
	if err != nil {
		logrus.Errorf("decode data to json object error: %v", err)
		return "", "", err
	}
	//TODO vnode.Ps vnode.Add vnode.Id 去掉空白字符
	if vnode.Tls != "tls" {
		vnode.Tls = "none"
	}

	// get v2Host v2Path
	if vnode.V == "2" {
		v2Host = vnode.Host
		v2Path = vnode.Path
	} else {
		switch vnode.Net {
		case "tcp":
			v2Host = vnode.Host
			v2Path = ""
		case "kcp":
			v2Host = ""
			v2Path = ""
		case "ws":
			v2HostTmp := vnode.Host
			if v2HostTmp != "" {
				slice := strings.Split(v2HostTmp, ";")
				if len(slice) > 0 {
					v2Host = slice[0]
					v2Path = slice[0]
				} else {
					v2Host = ""
					v2Path = v2Host
				}
			}
		case "h2":
			v2Host = ""
			v2Path = vnode.Path
		default:
			return "", "", fmt.Errorf("unknow net filed: %v", vnode.Net)
		}
	}

	vnode.Host = v2Host
	vnode.Path = v2Path

	var (
		tcp string = "null"
		kcp string = "null"
		ws  string = "null"
		h2  string = "null"
		tls string = "null"
	)

	if vnode.Tls == "tls" {
		tls = fmt.Sprintf(`{"allowInsecure":true,"serverName":"%v"}`, vnode.Host)
	} else {
		tls = "null"
	}

	if strings.Contains(vnode.Host, ",") {
		vnode.Host = strings.ReplaceAll(vnode.Host, ",", `","`)
	}

	w := bytes.Buffer{}
	switch vnode.Net {
	case "tcp":
		if vnode.Type == "http" {
			tmpl.ExecuteTemplate(&w, "tcp", vnode)
			tcp = w.String()
		} else {
			tcp = "null"
		}
	case "kcp":
		tmpl.ExecuteTemplate(&w, "kcp", vnode)
		kcp = w.String()
	case "ws":
		tmpl.ExecuteTemplate(&w, "ws", vnode)
		ws = w.String()
	case "h2":
		tmpl.ExecuteTemplate(&w, "h2", vnode)
		h2 = w.String()
	default:
		return "", "", fmt.Errorf("unknow net field: %v", vnode.Net)
	}
	ioutil.ReadAll(&w)

	m := map[string]string{
		"address":      vnode.Add,
		"port":         fmt.Sprintf("%d", vnode.Port),
		"id":           vnode.Id,
		"alterId":      fmt.Sprintf("%d", vnode.Aid),
		"network":      vnode.Net,
		"security":     vnode.Tls,
		"tlsSettings":  tls,
		"tcpSettings":  tcp,
		"kcpSettings":  kcp,
		"wsSettings":   ws,
		"httpSettings": h2,
	}
	tmpl.ExecuteTemplate(&w, "outbound", m)
	outbound := w.String()
	ioutil.ReadAll(&w)
	logrus.Infof("--------vmess node: %v", outbound)

	tmpl.ExecuteTemplate(&w, "full", map[string]string{"outbound": outbound})
	all := w.String()
	ioutil.ReadAll(&w)

	return vnode.Ps, all, nil
}
