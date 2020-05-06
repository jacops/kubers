package injector

import (
	"encoding/base64"

	"github.com/jacops/azure-keyvault-k8s/internal/agent"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
)

func (a *AgentInjector) newConfig() ([]byte, error) {
	agentConfig := agent.Config{
		Secrets: a.Secrets,
	}

	return agentConfig.Render()
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

	if credentialsSecret, ok := a.Annotations[AnnotationAgentAzureCredentialsSecret]; ok {
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
