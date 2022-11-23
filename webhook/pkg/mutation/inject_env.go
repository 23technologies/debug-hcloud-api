package mutation

import (
	"errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"os"
)

// injectCert is a container for the mutation injecting environment vars
type injectEnv struct {
	Logger logrus.FieldLogger
}

// injectCert implements the podMutator interface
var _ podMutator = (*injectEnv)(nil)

// Name returns the struct name
func (se injectEnv) Name() string {
	return "inject_env"
}

// Mutate returns a new mutated pod according to set env rules
func (se injectEnv) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	se.Logger = se.Logger.WithField("mutation", se.Name())
	if len(os.Args) != 3 {
		se.Logger.Warning("Can't mutate Pods, invalid number of arguments provided (needs 2) :/")
		return pod, errors.New("need exactly two arguments to mutate the pods")
	}
	mutatedPod := pod.DeepCopy()
	// build out env var slice
	envVars := []corev1.EnvVar{{
		Name:  os.Args[1],
		Value: os.Args[2],
	}}

	// inject env vars into pod
	for _, envVar := range envVars {
		se.Logger.Debugf("pod env injected %s", envVar)
		injectEnvVar(mutatedPod, envVar)
	}

	return mutatedPod, nil
}

// injectEnvVar injects a var in both containers and init containers of a pod
func injectEnvVar(pod *corev1.Pod, envVar corev1.EnvVar) {
	for i, container := range pod.Spec.Containers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.Containers[i].Env = append(container.Env, envVar)
		}
	}
	for i, container := range pod.Spec.InitContainers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.InitContainers[i].Env = append(container.Env, envVar)
		}
	}
}

// HasEnvVar returns true if environment variable exists false otherwise
func HasEnvVar(container corev1.Container, checkEnvVar corev1.EnvVar) bool {
	for _, envVar := range container.Env {
		if envVar.Name == checkEnvVar.Name {
			return true
		}
	}
	return false
}
