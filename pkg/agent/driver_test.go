package agent

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jacops/kubers/pkg/driver"
)

func TestDriverCanBeFoundByName(t *testing.T) {
	tests := []struct {
		name       string
		driverType string
	}{
		{name: "azure", driverType: "*azure.Driver"},
		{name: "aws", driverType: "*aws.Driver"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver, err := newDriver(tt.name, &driver.Config{}, getLogger())
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
			driver, err := newDriver(tt.name, &driver.Config{}, getLogger())

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
