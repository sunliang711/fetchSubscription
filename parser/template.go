package parser

const vmessTmpl = `
{{ define "tcp" }}
			{
				"connectionReuse": true,
				"header": {
					"type": "http",
					"request": {
						"version": "1.1",
						"method": "GET",
						"path": ["/"],
						"headers": {
							"Host": ["{{.Host}}"],
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
{{ end }}

{{ define "kcp" }}
		{
			"mtu": 1350,
			"tti": 50,
			"uplinkCapacity": 12,
			"downlinkCapacity": 100,
			"congestion": false,
			"readBufferSize": 2,
			"writeBufferSize": 2,
			"header": {
				"type": "{{.Type}}",
				"request": null,
				"response": null
			}
		}
{{ end }}

{{ define "ws" }}
		{
			"connectionReuse": true,
			"path": "{{.Path}}",
			"headers": {
				"Host": "{{.Host}}"
			}
		}
{{ end }}

{{ define "h2" }}
		{
			"path": "{{.Path}}",
			"headers": {
				"Host": "{{.Host}}"
			}
		}
{{ end }}

{{ define "outbound" }}
		"outbound": {
			"protocol": "vmess",
			"settings": {
				"vnext": [
					{
						"address": "{{.address}}",
						"port": {{.port}},
						"users": [
							{
								"id": "{{.id}}",
								"alterId": {{.alterId}},
								"security": "auto"
							}
						]
					}
				]
			},
			"streamSettings": {
				"network": "{{.network}}",
				"security": "{{.security}}",
				"tlsSettings": {{.tlsSettings}},
				"tcpSettings": {{.tcpSettings}},
				"kcpSettings": {{.kcpSettings}},
				"wsSettings": {{.wsSettings}},
				"httpSettings": {{.httpSettings}}
			},
			"mux": {
				"enabled": true
			}
		}

{{ end }}
`
