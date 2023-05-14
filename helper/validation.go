package helper

import "net/url"

func IsValidURL(u string) bool {
	_, err := url.Parse(u)
	return err == nil
}
