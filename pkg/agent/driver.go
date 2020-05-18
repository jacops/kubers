package agent

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	awsDriver "github.com/jacops/kubers/pkg/agent/driver/aws"
	azureDriver "github.com/jacops/kubers/pkg/agent/driver/azure"
	"github.com/jacops/kubers/pkg/driver"
)

func newDriver(name string, config *driver.Config, logger hclog.Logger) (driver.Driver, error) {
	switch name {
	case "azure":
		return azureDriver.New(config, logger), nil
	case "aws":
		return awsDriver.New(config, logger), nil
	}

	return nil, fmt.Errorf("Driver not found: %s", name)
}
