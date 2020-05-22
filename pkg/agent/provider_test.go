package agent

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/pkg/provider"
)

func TestProviderCanBeFoundByName(t *testing.T) {
	tests := []struct {
		name         string
		providerType string
	}{
		{name: "azure", providerType: "*azure.Provider"},
		{name: "aws", providerType: "*aws.Provider"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := newProvider(tt.name, &provider.Config{}, hclog.New(&hclog.LoggerOptions{
				Name: "handler",
			}))
			providerTypeStr := reflect.TypeOf(provider).String()
			if providerTypeStr != tt.providerType {
				t.Errorf("Unexpected provider type (%s) by name (%s)", providerTypeStr, tt.name)
			}
			if err != nil {
				t.Errorf("did not expect any errors: %s", err)
			}
		})
	}
}

func TestNotFoundProviderThrowingError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "nonext-provider"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := newProvider(tt.name, &provider.Config{}, hclog.New(&hclog.LoggerOptions{
				Name: "handler",
			}))

			if provider != nil {
				t.Errorf("did not expect any provider: %s", tt.name)
			}

			if strings.Contains(err.Error(), "Provider not found") {
				return
			}

			if err != nil {
				t.Errorf("did not expect any other errors: %s", err)
			}
		})
	}
}
