package injector

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jacops/azure-keyvault-k8s/internal/agent"
	"github.com/mattbaird/jsonpatch"
	corev1 "k8s.io/api/core/v1"
)

const (
	DefaultVaultImage = "jacops/kubersctl:0.1.1"
)

// AgentInjector is the top level structure holding all the
// configurations for the Vault Agent container.
type AgentInjector struct {

	// Annotations are the current pod annotations used to
	// configure the Vault Agent container.
	Annotations map[string]string

	// ImageName is the name of the Vault image to use for the
	// sidecar container.
	ImageName string

	// Inject is the flag used to determine if a container should be requested
	// in a pod request.
	Inject bool

	// Patches are all the mutations we will make to the pod request.
	Patches []*jsonpatch.JsonPatchOperation

	// Pod is the original Kubernetes pod spec.
	Pod *corev1.Pod

	// Secrets are all the templates, the path in Vault where the secret can be
	//found, and the unique name of the secret which will be used for the filename.
	Secrets []*agent.Secret

	// Status is the current injection status.  The only status considered is "injected",
	// which prevents further mutations.  A user can patch this annotation to force a new
	// mutation.
	Status string
}

// New creates a new instance of Agent by parsing all the Kubernetes annotations.
func New(pod *corev1.Pod, patches []*jsonpatch.JsonPatchOperation) (*AgentInjector, error) {
	agentInjector := &AgentInjector{
		Annotations: pod.Annotations,
		ImageName:   pod.Annotations[AnnotationAgentImage],
		Patches:     patches,
		Pod:         pod,
		Status:      pod.Annotations[AnnotationAgentStatus],
	}

	var err error

	agentInjector.Secrets = agentInjector.secrets()
	agentInjector.Inject, err = agentInjector.inject()
	if err != nil {
		return agentInjector, err
	}

	return agentInjector, nil
}

// Validate the instance of Agent to ensure we have everything needed
// for basic functionality.
func (a *AgentInjector) Validate() error {
	return nil
}

// Patch creates the necessary pod patches to inject the Vault Agent
// containers.
func (a *AgentInjector) Patch() ([]byte, error) {
	var patches []byte

	// Add our volume that will be shared by the containers
	// for passing data in the pod.
	a.Patches = append(a.Patches, addVolumes(
		a.Pod.Spec.Volumes,
		a.ContainerVolumes(),
		"/spec/volumes")...)

	//Add Volume Mounts
	for i, container := range a.Pod.Spec.Containers {
		a.Patches = append(a.Patches, addVolumeMounts(
			container.VolumeMounts,
			a.ContainerVolumeMounts(),
			fmt.Sprintf("/spec/containers/%d/volumeMounts", i))...)
	}

	container, err := a.ContainerInitSidecar()
	if err != nil {
		return patches, err
	}

	containers := a.Pod.Spec.InitContainers

	a.Patches = append(a.Patches, addContainers(
		a.Pod.Spec.InitContainers,
		[]corev1.Container{container},
		"/spec/initContainers")...)

	//Add Volume Mounts
	for i, container := range containers {
		if container.Name == "vault-agent-init" {
			continue
		}
		a.Patches = append(a.Patches, addVolumeMounts(
			container.VolumeMounts,
			a.ContainerVolumeMounts(),
			fmt.Sprintf("/spec/initContainers/%d/volumeMounts", i))...)
	}

	// Add annotations so that we know we're injected
	a.Patches = append(a.Patches, updateAnnotations(
		a.Pod.Annotations,
		map[string]string{AnnotationAgentStatus: "injected"})...)

	// Generate the patch
	if len(a.Patches) > 0 {
		var err error
		patches, err = json.Marshal(a.Patches)
		if err != nil {
			return patches, err
		}
	}
	return patches, nil
}

// ShouldInject checks whether the pod in question should be injected
// with Vault Agent containers.
func ShouldInject(pod *corev1.Pod) (bool, error) {
	raw, ok := pod.Annotations[AnnotationAgentInject]
	if !ok {
		return false, nil
	}

	inject, err := strconv.ParseBool(raw)
	if err != nil {
		return false, err
	}

	if !inject {
		return false, nil
	}

	// This shouldn't happen so bail.
	raw, ok = pod.Annotations[AnnotationAgentStatus]
	if !ok {
		return true, nil
	}

	// "injected" is the only status we care about.  Don't do
	// anything if it's set.  The user can update the status
	// to force a new mutation.
	if raw == "injected" {
		return false, nil
	}

	return true, nil
}
