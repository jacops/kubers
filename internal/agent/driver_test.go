package agent

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jacops/azure-keyvault-k8s/internal/agent/driver"
)

func TestDriverCanBeFoundInMapByURL(t *testing.T) {
	tests := []struct {
		url        string
		driverType string
	}{
		{url: "keyvault://kv/key", driverType: "*driver.KeyVaultDriver"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			driver, err := getSecretsDriverFromMapByURL(tt.url, driver.Config{})
			driverTypeStr := reflect.TypeOf(driver).String()
			if driverTypeStr != tt.driverType {
				t.Errorf("Unexpected driver type (%s) from URL (%s)", driverTypeStr, tt.url)
			}
			if err != nil {
				t.Errorf("did not expect any errors: %s", err)
			}
		})
	}
}

func TestGetSecretsDriverFromMapByURLErrors(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{url: "txt:///path/key", expected: "Unable to find a driver"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			_, err := getSecretsDriverFromMapByURL(tt.url, driver.Config{})

			if !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("Did not expect any other errors: %s", err)
			}
		})
	}
}
