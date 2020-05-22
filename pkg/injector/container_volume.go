package injector

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/strutil"
	corev1 "k8s.io/api/core/v1"
)

const (
	secretVolumeName = "vault-secrets"
	secretVolumePath = "/vault/secrets"
)

func (a *AgentInjector) getUniqueMountPaths() []string {
	var mountPaths []string

	for _, secret := range a.Secrets {
		if !strutil.StrListContains(mountPaths, secret.MountPath) && secret.MountPath != a.Annotations[AnnotationVaultSecretVolumePath] {
			mountPaths = append(mountPaths, secret.MountPath)
		}
	}
	return mountPaths
}

// ContainerVolumes returns the volume data to add to the pod. This volumes
// are used for shared data between containers.
func (a *AgentInjector) ContainerVolumes() []corev1.Volume {
	containerVolumes := []corev1.Volume{
		{
			Name: secretVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: "Memory",
				},
			},
		},
	}
	for index := range a.getUniqueMountPaths() {
		containerVolumes = append(
			containerVolumes,
			corev1.Volume{
				Name: fmt.Sprintf("%s-custom-%d", secretVolumeName, index),
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						Medium: "Memory",
					},
				},
			},
		)
	}
	return containerVolumes
}

// ContainerVolumeMounts mounts the shared memory volume where secrets
// will be rendered.
func (a *AgentInjector) ContainerVolumeMounts() []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      secretVolumeName,
			MountPath: a.Annotations[AnnotationVaultSecretVolumePath],
			ReadOnly:  false,
		},
	}
	for index, mountPath := range a.getUniqueMountPaths() {
		volumeMounts = append(
			volumeMounts,
			corev1.VolumeMount{
				Name:      fmt.Sprintf("%s-custom-%d", secretVolumeName, index),
				MountPath: mountPath,
				ReadOnly:  false,
			},
		)
	}
	return volumeMounts
}
