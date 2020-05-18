package injector

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/jacops/kubers/pkg/agent"
	"github.com/jacops/kubers/pkg/provider"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
)

func (a *AgentInjector) newConfig() ([]byte, error) {
	var providerName string
	var renderedConfig []byte

	if _, ok := a.Annotations[AnnotationAgentProvider]; !ok {
		return renderedConfig, errors.New("Provider not specified")
	}

	providerName = a.Annotations[AnnotationAgentProvider]

	agentConfig := agent.Config{
		Secrets:      a.Secrets,
		ProviderName:   providerName,
		ProviderConfig: a.getProviderConfig(providerName),
	}

	return agentConfig.Render()
}

func (a *AgentInjector) getProviderConfig(providerName string) *provider.Config {
	providerConfig := make(provider.Config)
	providerConfigAnnotationPrefix := fmt.Sprintf("%s-%s-", AnnotationAgentProvider, providerName)

	for name, val := range a.Annotations {
		if strings.Contains(name, providerConfigAnnotationPrefix) {
			annotationConfigKey := strings.ReplaceAll(name, providerConfigAnnotationPrefix, "")
			providerConfig[strings.ReplaceAll(annotationConfigKey, "-", "_")] = val
		}
	}

	return &providerConfig
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

	if credentialsSecret, ok := a.Annotations[AnnotationAgentProviderAzureCredentialsSecret]; ok {
		envs = append(envs, getAzureCredentialsEnvs(credentialsSecret)...)
	} else if credentialsSecret, ok := a.Annotations[AnnotationAgentProviderAWSCredentialsSecret]; ok {
		envs = append(envs, getAWSCredentialsEnvs(credentialsSecret)...)
	}

	return envs, nil
}

func getAzureCredentialsEnvs(credentialsSecret string) []corev1.EnvVar {
	var envs []corev1.EnvVar

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

	return envs
}

func getAWSCredentialsEnvs(credentialsSecret string) []corev1.EnvVar {
	var envs []corev1.EnvVar

	envs = append(envs, corev1.EnvVar{
		Name:      "AWS_ACCESS_KEY_ID",
		ValueFrom: getEnvVarFromSecret(credentialsSecret, "accessKeyId"),
	})

	envs = append(envs, corev1.EnvVar{
		Name:      "AWS_SECRET_ACCESS_KEY",
		ValueFrom: getEnvVarFromSecret(credentialsSecret, "secretAccessKey"),
	})

	return envs
}

func getEnvVarFromSecret(credentialsSecret string, key string) *corev1.EnvVarSource {
	return &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: v1.LocalObjectReference{Name: credentialsSecret},
			Key:                  key,
		},
	}
}
