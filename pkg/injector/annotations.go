package injector

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jacops/kubers/pkg/agent"
	corev1 "k8s.io/api/core/v1"
)

const (
	// AnnotationAgentDriver sets a driver for the agent used to retrieve secrets
	AnnotationAgentDriver = "kubers.jacops.pl/agent-driver"

	// AnnotationAgentDriverAzureCredentialsSecret enables authentication via Azure service principal
	AnnotationAgentDriverAzureCredentialsSecret = "kubers.jacops.pl/agent-driver-azure-credentials-secret"

	// AnnotationAgentDriverAWSCredentialsSecret enables authentication via AWS keys
	AnnotationAgentDriverAWSCredentialsSecret = "kubers.jacops.pl/agent-driver-aws-credentials-secret"

	// AnnotationAgentDriverAWSRegion overrides aws region passed to the operator
	AnnotationAgentDriverAWSRegion = "kubers.jacops.pl/agent-driver-aws-region"

	// AnnotationAgentStatus is the key of the annotation that is added to
	// a pod after an injection is done.
	// There's only one valid status we care about: "injected".
	AnnotationAgentStatus = "kubers.jacops.pl/agent-inject-status"

	// AnnotationAgentInject is the key of the annotation that controls whether
	// injection is explicitly enabled or disabled for a pod. This should
	// be set to a true or false value, as parseable by strconv.ParseBool
	AnnotationAgentInject = "kubers.jacops.pl/agent-inject"

	// AnnotationAgentInjectSecret is the key annotation that configures Vault
	// Agent to retrieve the secrets from Vault required by the app.  The name
	// of the secret is any unique string after "vault.hashicorp.com/agent-inject-secret-",
	// such as "vault.hashicorp.com/agent-inject-secret-foobar".  The value is the
	// path in Vault where the secret is located.
	AnnotationAgentInjectSecret = "kubers.jacops.pl/agent-inject-secret"

	// AnnotationAgentImage is the name of the Vault docker image to use.
	AnnotationAgentImage = "kubers.jacops.pl/agent-image"

	// AnnotationVaultSecretVolumePath specifies where the secrets are to be
	// Mounted after fetching.
	AnnotationVaultSecretVolumePath = "kubers.jacops.pl/secret-volume-path"

	// AnnotationPreserveSecretCase if enabled will preserve the case of secret name
	// by default the name is converted to lower case.
	AnnotationPreserveSecretCase = "kubers.jacops.pl/preserve-secret-case"
)

// AgentInjectorConfig ...
type AgentInjectorConfig struct {
	Image      string
	DriverName string
	AWSRegion  string
}

// Init configures the expected annotations required to create a new instance
// of Agent.  This should be run before running new to ensure all annotations are
// present.
func Init(pod *corev1.Pod, cfg AgentInjectorConfig) error {
	if pod == nil {
		return errors.New("pod is empty")
	}

	if pod.ObjectMeta.Annotations == nil {
		pod.ObjectMeta.Annotations = make(map[string]string)
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationVaultSecretVolumePath]; !ok {
		pod.ObjectMeta.Annotations[AnnotationVaultSecretVolumePath] = secretVolumePath
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentDriver]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentDriver] = cfg.DriverName
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentDriverAWSRegion]; !ok {
		pod.ObjectMeta.Annotations[AnnotationAgentDriverAWSRegion] = cfg.AWSRegion
	}

	if _, ok := pod.ObjectMeta.Annotations[AnnotationAgentImage]; !ok {
		if cfg.Image == "" {
			return errors.New("No agent image provided")
		}
		pod.ObjectMeta.Annotations[AnnotationAgentImage] = cfg.Image
	}

	return nil
}

// secrets parses annotations with the pattern "kubers.jacops.com/agent-inject-secret-".
// Everything following the final dash becomes the name of the secret,
// and the value is the path in Vault.
func (a *AgentInjector) secrets() []*agent.SecretMetadata {
	var secrets []*agent.SecretMetadata

	for name, path := range a.Annotations {
		secretName := fmt.Sprintf("%s-", AnnotationAgentInjectSecret)
		if strings.Contains(name, secretName) {
			raw := strings.ReplaceAll(name, secretName, "")
			name := raw

			if ok, _ := a.preserveSecretCase(raw); !ok {
				name = strings.ToLower(raw)
			}

			if name == "" {
				continue
			}

			mountPath := a.Annotations[AnnotationVaultSecretVolumePath]
			mountPathAnnotationName := fmt.Sprintf("%s-%s", AnnotationVaultSecretVolumePath, raw)

			if val, ok := a.Annotations[mountPathAnnotationName]; ok {
				mountPath = val
			}

			secrets = append(secrets, &agent.SecretMetadata{
				Name:      name,
				URL:       path,
				MountPath: mountPath,
			})
		}
	}
	return secrets
}

func (a *AgentInjector) inject() (bool, error) {
	raw, ok := a.Annotations[AnnotationAgentInject]
	if !ok {
		return true, nil
	}

	return strconv.ParseBool(raw)
}

func (a *AgentInjector) preserveSecretCase(secretName string) (bool, error) {

	preserveSecretCaseAnnotationName := fmt.Sprintf("%s-%s", AnnotationPreserveSecretCase, secretName)

	var raw string

	if val, ok := a.Annotations[preserveSecretCaseAnnotationName]; ok {
		raw = val
	} else {
		raw, ok = a.Annotations[AnnotationPreserveSecretCase]
		if !ok {
			return false, nil
		}
	}
	return strconv.ParseBool(raw)
}
