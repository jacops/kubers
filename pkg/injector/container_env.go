package injector

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/jacops/kubers/pkg/agent"
	"github.com/jacops/kubers/pkg/agent/driver"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
)

func (a *AgentInjector) newConfig() ([]byte, error) {
	var driverName string
	var renderedConfig []byte

	if _, ok := a.Annotations[AnnotationAgentDriver]; !ok {
		return renderedConfig, errors.New("Driver not specified")
	}

	driverName = a.Annotations[AnnotationAgentDriver]

	agentConfig := agent.Config{
		Secrets:      a.Secrets,
		DriverName:   driverName,
		DriverConfig: a.getDriverConfig(driverName),
	}

	return agentConfig.Render()
}

func (a *AgentInjector) getDriverConfig(driverName string) *driver.Config {
	driverConfig := make(driver.Config)
	driverConfigAnnotationPrefix := fmt.Sprintf("%s-%s-", AnnotationAgentDriver, driverName)

	for name, val := range a.Annotations {
		if strings.Contains(name, driverConfigAnnotationPrefix) {
			annotationConfigKey := strings.ReplaceAll(name, driverConfigAnnotationPrefix, "")
			driverConfig[strings.ReplaceAll(annotationConfigKey, "-", "_")] = val
		}
	}

	return &driverConfig
}

// ContainerEnvVars adds the applicable environment vars
// for the Vault Agent sidecar.
func (a *AgentInjector) ContainerEnvVars(init bool) ([]corev1.EnvVar, error) {
	var envs []corev1.EnvVar

	config, err := a.newConfig()
	if err != nil {
		return envs, err
	}

	b64Config := base64.StdEncoding.EncodeToString(config)
	envs = append(envs, corev1.EnvVar{
		Name:  "AGENT_CONFIG",
		Value: b64Config,
	})

	if credentialsSecret, ok := a.Annotations[AnnotationAgentDriverAzureCredentialsSecret]; ok {
		envs = append(envs, corev1.EnvVar{
			Name:      "AZURE_TENANT_ID",
			ValueFrom: getEnvVarFromSecret(credentialsSecret, "tenantId"),
		})

		envs = append(envs, corev1.EnvVar{
			Name:      "AZURE_SUBSCRIPTION_ID",
			ValueFrom: getEnvVarFromSecret(credentialsSecret, "subscriptionId"),
		})

		envs = append(envs, corev1.EnvVar{
			Name:      "AZURE_CLIENT_ID",
			ValueFrom: getEnvVarFromSecret(credentialsSecret, "clientId"),
		})

		envs = append(envs, corev1.EnvVar{
			Name:      "AZURE_CLIENT_SECRET",
			ValueFrom: getEnvVarFromSecret(credentialsSecret, "clientSecret"),
		})
	}

	return envs, nil
}

func getEnvVarFromSecret(credentialsSecret string, key string) *corev1.EnvVarSource {
	return &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: v1.LocalObjectReference{Name: credentialsSecret},
			Key:                  key,
		},
	}
}
