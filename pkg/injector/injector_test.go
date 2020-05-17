package injector

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func testPod(annotations map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "foo",
			Annotations: annotations,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "foobar",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "foobar",
							MountPath: "serviceaccount/somewhere",
						},
					},
				},
			},
		},
	}
}

func TestShouldInject(t *testing.T) {
	tests := []struct {
		annotations map[string]string
		inject      bool
	}{
		{map[string]string{AnnotationAgentInject: "true", AnnotationAgentStatus: ""}, true},
		{map[string]string{AnnotationAgentInject: "false", AnnotationAgentStatus: ""}, false},
		{map[string]string{AnnotationAgentInject: "true", AnnotationAgentStatus: "injected"}, false},
		{map[string]string{AnnotationAgentInject: "false", AnnotationAgentStatus: "injected"}, false},
		{map[string]string{AnnotationAgentInject: "true", AnnotationAgentStatus: "update"}, true},
	}

	for _, tt := range tests {
		pod := testPod(tt.annotations)
		inject, err := ShouldInject(pod)
		if err != nil {
			t.Errorf("got error, shouldn't have: %s", err)
		}

		if inject != tt.inject {
			t.Errorf("expected should inject to be %v, got %v", tt.inject, inject)
		}
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		agent AgentInjector
		valid bool
	}{
		{
			AgentInjector{
				ImageName: "test",
			}, true,
		},
	}

	for _, tt := range tests {
		err := tt.agent.Validate()

		if err != nil && tt.valid {
			t.Errorf("got error, shouldn't have: %s", err)
		}

		if err == nil && !tt.valid {
			t.Error("got no error, should have")
		}
	}
}
