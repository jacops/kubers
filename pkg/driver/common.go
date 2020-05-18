package driver

import (
	"net/url"
	"strings"
)

//GetServiceTypeFromURL returns service type
func GetServiceTypeFromURL(secretURL string) string {
	u, err := url.Parse(secretURL)
	if err != nil {
		return ""
	}

	return u.Scheme
}

//GetServiceKeyNamesPair returns serviceName and keyName
func GetServiceKeyNamesPair(secretURL string) (serviceName string, keyName string, err error) {
	u, err := url.Parse(secretURL)
	if err != nil {
		return
	}

	serviceName = strings.Trim(u.Host, "/")
	keyName = strings.Trim(u.Path, "/")

	return
}

//GetKeyName returns key name from URL
func GetKeyName(secretURL string) (keyName string, err error) {
	u, err := url.Parse(secretURL)
	if err != nil {
		return
	}

	keyName = strings.Trim(u.Host, "/")

	return
}
