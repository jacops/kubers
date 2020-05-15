package driver

import (
	"reflect"
	"strings"
	"testing"
)

func TestDriverCanBeFoundByName(t *testing.T) {
	tests := []struct {
		name       string
		driverType string
	}{
		{name: "keyvault", driverType: "*driver.KeyVaultDriver"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewSimpleFactory()
			driver, err := factory.GetDriver(tt.name, Config{})
			driverTypeStr := reflect.TypeOf(driver).String()
			if driverTypeStr != tt.driverType {
				t.Errorf("Unexpected driver type (%s) by name (%s)", driverTypeStr, tt.name)
			}
			if err != nil {
				t.Errorf("did not expect any errors: %s", err)
			}
		})
	}
}

func TestNotFoundDriverThrowingError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "nonext-driver"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewSimpleFactory()
			driver, err := factory.GetDriver(tt.name, Config{})

			if driver != nil {
				t.Errorf("did not expect any driver: %s", tt.name)
			}

			if strings.Contains(err.Error(), "Driver not found") {
				return
			}

			if err != nil {
				t.Errorf("did not expect any other errors: %s", err)
			}
		})
	}
}

func TestUrlCanReturnServiceKeyNamesPair(t *testing.T) {
	tests := []struct {
		url                 string
		expectedServiceName string
		expectedKey         string
	}{
		{url: "keyvault://keyvault/key", expectedServiceName: "keyvault", expectedKey: "key"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			serviceName, keyName, err := getServiceKeyNamesPair(tt.url)

			if serviceName != tt.expectedServiceName {
				t.Errorf("expected service name: %s; got: %s", tt.expectedServiceName, serviceName)
			}

			if keyName != tt.expectedKey {
				t.Errorf("expected key name: %s; got: %s", tt.expectedKey, keyName)
			}

			if err != nil {
				t.Errorf("did not expect any errors: %s", err)
			}
		})
	}
}
