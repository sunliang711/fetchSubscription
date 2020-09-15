package parser

import (
	"fetchSubscription/decoder"
	"fetchSubscription/downloader"
	"testing"
)

func TestParse(t *testing.T) {
	nodesContent := `vmess://eyJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6IiIsImFkZCI6IjE0LjE3Ljk3LjE0NSIsInBvcnQiOjUwMTAsImFpZCI6MiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiIsInBzIjoi5Y+w5rm+IDAxIFtEMy9WUi9JUExDXSIsImlkIjoiMzI0NzBlMTQtODVmYi0zYmYwLWFhMGMtMWY3YmE0NmI1OGI3IiwiY2xhc3MiOjN9
vmess://eyJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6IiIsImFkZCI6IjE0LjE3Ljk3LjE0NSIsInBvcnQiOjUwMDgsImFpZCI6MiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiIsInBzIjoi5pel5pysIDAxIFtEMy9WUi9JUExDXSIsImlkIjoiMzI0NzBlMTQtODVmYi0zYmYwLWFhMGMtMWY3YmE0NmI1OGI3IiwiY2xhc3MiOjN9
vmess://eyJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6IiIsImFkZCI6IjE0LjE3Ljk3LjE0NSIsInBvcnQiOjUwMDYsImFpZCI6MiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiIsInBzIjoi54uu5Z+OIDAxIFtEMy9WUi9JUExDXSIsImlkIjoiMzI0NzBlMTQtODVmYi0zYmYwLWFhMGMtMWY3YmE0NmI1OGI3IiwiY2xhc3MiOjN9
vmess://eyJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6IiIsImFkZCI6IjE0LjE3Ljk3LjE0NSIsInBvcnQiOjUwMTIsImFpZCI6MiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiIsInBzIjoi576O5Zu9IDAxIFtEMy9WUi9JUExDXSIsImlkIjoiMzI0NzBlMTQtODVmYi0zYmYwLWFhMGMtMWY3YmE0NmI1OGI3IiwiY2xhc3MiOjN9
vmess://eyJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6IiIsImFkZCI6ImlwbGMtMDEuc3RjLWJncC5jb20iLCJwb3J0Ijo1MDAwLCJhaWQiOjIsIm5ldCI6InRjcCIsInR5cGUiOiJub25lIiwidiI6IjIiLCJwcyI6Iummmea4ryAwMSBbRDMvVlIvSVBMQ10iLCJpZCI6IjMyNDcwZTE0LTg1ZmItM2JmMC1hYTBjLTFmN2JhNDZiNThiNyIsImNsYXNzIjozfQ==
vmess://eyJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6IiIsImFkZCI6IjE0LjE3Ljk3LjE0NSIsInBvcnQiOjUwMDIsImFpZCI6MiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiIsInBzIjoi6aaZ5rivIDAyIFtEMy9WUi9JUExDXSIsImlkIjoiMzI0NzBlMTQtODVmYi0zYmYwLWFhMGMtMWY3YmE0NmI1OGI3IiwiY2xhc3MiOjN9
vmess://eyJob3N0IjoiIiwicGF0aCI6IiIsInRscyI6IiIsImFkZCI6IjE0LjE3Ljk3LjE0NSIsInBvcnQiOjUwMDQsImFpZCI6MiwibmV0IjoidGNwIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiIsInBzIjoi6aaZ5rivIDAzIFtEMy9WUi9JUExDXSIsImlkIjoiMzI0NzBlMTQtODVmYi0zYmYwLWFhMGMtMWY3YmE0NmI1OGI3IiwiY2xhc3MiOjN9`
	cfg := &FilterConfig{}
	ret, err := Parse(nodesContent, cfg)
	if err != nil {
		t.Fatalf("Parse nodesContent error: %v", err)
	}
	t.Logf("nodesContent parse ret: %v", ret)
}

func TestAll(t *testing.T) {
	subURL := "https://stc-spadesdns.com/link/PUoCxdJB5l4ploeF?sub=3&extend=1"
	content, err := downloader.Download(subURL, nil)
	if err != nil {
		t.Fatalf("download error: %v", err)
	}
	decoded, err := decoder.Decode(content)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	nodes, err := Parse(decoded, nil)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	for name, node := range nodes {
		t.Logf("name: %v node: %v", name, node)
	}
}
