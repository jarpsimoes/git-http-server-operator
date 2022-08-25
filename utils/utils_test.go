package utils

import (
	githttpserver1alpha1 "github.com/jarpsimoes/git-http-server-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestMergeConfigurationWithEnvironmentVariables(t *testing.T) {
	dummyObjectIncomplete := githttpserver1alpha1.GitHttpServer{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: githttpserver1alpha1.GitHttpServerSpec{
			Image:      "dummy-image",
			PathClone:  "_clone1",
			PathPull:   "_pull1",
			RepoBranch: "https://anybranch",
			RepoTarget: "main",
		},
		Status: githttpserver1alpha1.GitHttpServerStatus{},
	}
	variables := MergeConfigurationWithEnvironmentVariables(dummyObjectIncomplete)

	for _, variable := range variables {
		if variable.Name == "PATH_CLONE" {
			assert.Equal(t, "_clone1", variable.Value, "Check customized clone path")
		}
		if variable.Name == "PATH_VERSION" {
			assert.Equal(t, "_version", variable.Value, "Check default version path")
		}

	}
	dummyObjectWithPasswordIncomplete := githttpserver1alpha1.GitHttpServer{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: githttpserver1alpha1.GitHttpServerSpec{
			Image:        "dummy-image",
			PathClone:    "_clone1",
			PathPull:     "_pull1",
			RepoBranch:   "https://anybranch",
			RepoTarget:   "main",
			RepoPassword: "AnyPassword",
			RepoUsername: "AnyUsername",
			HttpPort:     8082,
		},
		Status: githttpserver1alpha1.GitHttpServerStatus{},
	}

	variablesWithPassword := MergeConfigurationWithEnvironmentVariables(dummyObjectWithPasswordIncomplete)

	for _, variable := range variablesWithPassword {
		if variable.Name == "PATH_PULL" {
			assert.Equal(t, "_pull1", variable.Value, "Check customized pull path")
		}
		if variable.Name == "REPO_USERNAME" {
			assert.Equal(t, "AnyUsername", variable.Value, "Check username")
		}
		if variable.Name == "REPO_PASSWORD" {
			assert.Equal(t, "AnyPassword", variable.Value, "Check password")
		}
		if variable.Name == "HTTP_PORT" {
			assert.Equal(t, "8082", variable.Value, "Check to string from int32")
		}
	}
}
func TestGetProbe(t *testing.T) {
	dummyObjectIncomplete := githttpserver1alpha1.GitHttpServer{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: githttpserver1alpha1.GitHttpServerSpec{
			PathHealth: "_health1",
			HttpPort:   8082,
		},
		Status: githttpserver1alpha1.GitHttpServerStatus{},
	}

	probe := GetProbe(dummyObjectIncomplete)

	assert.Equal(t, "/_health1", probe.HTTPGet.Path)
	assert.Equal(t, int32(8082), probe.HTTPGet.Port.IntVal)
}
