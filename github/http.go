package github

import "github.com/go-resty/resty/v2"

var defaultClient *resty.Client

func init() {
	defaultClient = resty.New()
}

// DownloadFile implements a web file download.
func DownloadFile(src, dest string) error {
	_, err := defaultClient.R().SetOutput(dest).Get(src)
	return err
}
