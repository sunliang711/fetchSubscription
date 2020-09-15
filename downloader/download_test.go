package downloader

import "testing"

func TestDownload(t *testing.T) {
	subURL := "https://stc-spadesdns.com/link/PUoCxdJB5l4ploeF?sub=3&extend=1"
	headers := make(map[string][]string)

	content, err := Download(subURL, headers)
	if err != nil {
		t.Fatalf("download error: %v", err)
	}
	t.Logf("download content: %v", content)

}
