package provider

import "testing"

func TestUrlCanReturnServiceKeyNamesPair(t *testing.T) {
	tests := []struct {
		url                 string
		expectedServiceName string
		expectedKey         string
	}{
		{url: "keyvault://keyvault/key", expectedServiceName: "keyvault", expectedKey: "key"},
		{url: "keyvault://secretmanager/somekey", expectedServiceName: "secretmanager", expectedKey: "somekey"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			serviceName, keyName, err := GetServiceKeyNamesPair(tt.url)

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

func TestUrlCanReturnServiceType(t *testing.T) {
	tests := []struct {
		url                 string
		expectedServiceType string
	}{
		{url: "keyvault://name/key", expectedServiceType: "keyvault"},
		{url: "secretmanager://name/somekey", expectedServiceType: "secretmanager"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			serviceType := GetServiceTypeFromURL(tt.url)

			if serviceType != tt.expectedServiceType {
				t.Errorf("expected service type: %s; got: %s", tt.expectedServiceType, serviceType)
			}
		})
	}
}

func TestUrlCanReturnKeyName(t *testing.T) {
	tests := []struct {
		url         string
		expectedKey string
	}{
		{url: "keyvault://name", expectedKey: "name"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			keyName, err := GetKeyName(tt.url)

			if keyName != tt.expectedKey {
				t.Errorf("expected key name: %s; got: %s", tt.expectedKey, keyName)
			}

			if err != nil {
				t.Errorf("did not expect any errors: %s", err)
			}
		})
	}
}
