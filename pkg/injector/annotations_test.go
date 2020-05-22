package injector

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/mattbaird/jsonpatch"
)

func TestInitCanSet(t *testing.T) {
	annotations := make(map[string]string)
	pod := testPod(annotations)

	err := Init(pod, &AgentInjectorConfig{"foobar-image", "", "", "", ""})
	if err != nil {
		t.Errorf("got error, shouldn't have: %s", err)
	}

	tests := []struct {
		annotationKey   string
		annotationValue string
	}{
		{annotationKey: AnnotationAgentImage, annotationValue: "foobar-image"},
	}

	for _, tt := range tests {
		raw, ok := pod.Annotations[tt.annotationKey]
		if !ok {
			t.Errorf("Default annotation %s not set, it should be.", tt.annotationKey)
		}

		if raw != tt.annotationValue {
			t.Errorf("Default annotation confiured value incorrect, wanted %s, got %s", tt.annotationValue, raw)

		}
	}
}

func TestInitDefaults(t *testing.T) {
	annotations := make(map[string]string)
	pod := testPod(annotations)

	err := Init(pod, &AgentInjectorConfig{DefaultVaultImage, "", "", "", ""})
	if err != nil {
		t.Errorf("got error, shouldn't have: %s", err)
	}

	tests := []struct {
		annotationKey   string
		annotationValue string
	}{
		{annotationKey: AnnotationAgentImage, annotationValue: DefaultVaultImage},
	}

	for _, tt := range tests {
		raw, ok := pod.Annotations[tt.annotationKey]
		if !ok {
			t.Errorf("Default annotation %s not set, it should be.", tt.annotationKey)
		}

		if raw != tt.annotationValue {
			t.Errorf("Default annotation value incorrect, wanted %s, got %s", tt.annotationValue, raw)
		}
	}
}

func TestCouldErrorAnnotations(t *testing.T) {
	tests := []struct {
		key   string
		value string
		valid bool
	}{
		{AnnotationAgentInject, "true", true},
		{AnnotationAgentInject, "false", true},
		{AnnotationAgentInject, "TRUE", true},
		{AnnotationAgentInject, "FALSE", true},
		{AnnotationAgentInject, "0", true},
		{AnnotationAgentInject, "1", true},
		{AnnotationAgentInject, "t", true},
		{AnnotationAgentInject, "f", true},
		{AnnotationAgentInject, "tRuE", false},
		{AnnotationAgentInject, "fAlSe", false},
		{AnnotationAgentInject, "", false},
	}

	for i, tt := range tests {
		annotations := map[string]string{tt.key: tt.value}
		pod := testPod(annotations)
		var patches []*jsonpatch.JsonPatchOperation

		err := Init(pod, &AgentInjectorConfig{"image", "", "", "", ""})
		if err != nil {
			t.Errorf("got error, shouldn't have: %s", err)
		}

		_, err = New(pod, patches)
		if err != nil && tt.valid {
			t.Errorf("[%d] got error, shouldn't have: %s", i, err)
		} else if err == nil && !tt.valid {
			t.Errorf("[%d] got no error, should have: %s", i, err)
		}
	}
}

func TestInitEmptyPod(t *testing.T) {
	var pod *corev1.Pod

	err := Init(pod, &AgentInjectorConfig{"foobar-image", "", "", "", ""})
	if err == nil {
		t.Errorf("got no error, should have")
	}
}
