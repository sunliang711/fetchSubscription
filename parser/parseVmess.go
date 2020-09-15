package parser

import (
	"encoding/json"
	"fetchSubscription/decoder"
	"fmt"
	"strings"

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
	err := json.NewDecoder(reader).Decode(&vnode)
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
		tcp string
		kcp string
		ws  string
		h2  string
		tls string
	)

	if vnode.Tls == "tls" {
		tls = fmt.Sprintf(`{"allowInsecure":true,"serverName":"%v"}`, vnode.Host)
	} else {
		tls = "null"
	}

	if strings.Contains(vnode.Host, ",") {
		vnode.Host = strings.ReplaceAll(vnode.Host, ",", `","`)
	}

	switch vnode.Net {
	case "tcp":
		if vnode.Type == "http" {
			//TODO replace with go template
			tcp = fmt.Sprintf(`
			{
				"connectionReuse": true,
				"header": {
					"type": "http",
					"request": {
						"version": "1.1",
						"method": "GET",
						"path": ["/"],
						"headers": {
							"Host": ["%v"],
							"User-Agent": ["Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36","Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/53.0.2785.109 Mobile/14A456 Safari/601.1.46"],
							"Accept-Encoding": ["gzip, deflate"],
							"Connection": ["keep-alive"],
							"Pragma": "no-cache"
						}
					},
					"response": {
						"version": "1.1",
						"status": "200",
						"reason": "OK",
						"headers": {
							"Content-Type": ["application/octet-stream","video/mpeg"],
							"Transfer-Encoding": ["chunked"],
							"Connection": ["keep-alive"],
							"Pragma": "no-cache"
						}
					}
				}
			}
			`, vnode.Host)
		} else {
			tcp = "null"
		}
	case "kcp":
		kcp = fmt.Sprintf(`
		{
			"mtu": 1350,
			"tti": 50,
			"uplinkCapacity": 12,
			"downlinkCapacity": 100,
			"congestion": false,
			"readBufferSize": 2,
			"writeBufferSize": 2,
			"header": {
				"type": "%v",
				"request": null,
				"response": null
			}
		}
		`, vnode.Type)
	case "ws":
		ws = fmt.Sprintf(`
		{
			"connectionReuse": true,
			"path": "%v",
			"headers": {
				"Host": "%v"
			}
		}
		`, vnode.Path, vnode.Host)
	case "h2":
		h2 = fmt.Sprintf(`
		{
			"path": "%v",
			"headers": {
				"Host": "%v"
			}
		}
		`, vnode.Path, vnode.Host)
	default:
		return "", "", fmt.Errorf("unknow net field: %v", vnode.Net)
	}

	outbound := fmt.Sprintf(`
		{
			"outbound": {
				"protocol": "vmess",
				"settings": {
					"vnext": [
						{
							"address": "%v",
							"port": %v,
							"users": [
								{
									"id": "%v",
									"alterId": %v,
									"security": "auto"
								}
							]
						}
					]
				},
				"streamSettings": {
					"network": "%v",
					"security": "%v",
					"tlsSettings": %v,
					"tcpSettings": %v,
					"kcpSettings": %v,
					"wsSettings": %v,
					"httpSettings": %v
				},
				"mux": {
					"enabled": true
				}
			}
		}"
	`, vnode.Add, vnode.Port, vnode.Id, vnode.Aid, vnode.Net, vnode.Tls, tls, tcp, kcp, ws, h2)
	logrus.Infof("vmess node: %v", outbound)

	return vnode.Ps, outbound, nil
}
