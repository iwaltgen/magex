package github

import "github.com/go-resty/resty/v2"

var defaultClient *resty.Client

func init() {
	defaultClient = resty.New()
}

func DownloadFile(url, dest string) error {
	_, err := defaultClient.R().SetOutput(dest).Get(url)
	return err
}
