package injector

import (
	corev1 "k8s.io/api/core/v1"
)

// ContainerInitSidecar creates a new init container to be added
// to the pod being mutated.  After Vault 1.4 is released, this can
// be removed because an exit_after_auth environment variable is
// available for the agent.  This means we won't need to generate
// two config files.
func (a *AgentInjector) ContainerInitSidecar() (corev1.Container, error) {
	volumeMounts := []corev1.VolumeMount{}
	volumeMounts = append(volumeMounts, a.ContainerVolumeMounts()...)

	envs, err := a.ContainerEnvVars(true)
	if err != nil {
		return corev1.Container{}, err
	}

	return corev1.Container{
		Name:            "kubers-init",
		Image:           a.ImageName,
		ImagePullPolicy: "IfNotPresent",
		Env:             envs,
		VolumeMounts:    volumeMounts,
		Args:            []string{"agent", "fetch-secrets"},
	}, nil
}
