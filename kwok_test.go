package kwok_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kwok "github.com/pablochacin/testcontainers-go-kwok"
	"github.com/testcontainers/testcontainers-go"
)

func TestKwok(t *testing.T) {
	ctx := context.Background()

	ctr, err := kwok.Run(ctx, "ghcr.io/pablochacin/kwok:latest")
	require.NoError(t, err)

	kubeConfigYaml , err := ctr.GetKubeConfig(ctx)
	require.NoError(t, err)

	restcfg, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigYaml)
	require.NoError(t, err)
    
	k8s, err := kubernetes.NewForConfig(restcfg)
	require.NoError(t, err)
    
	pod := &corev1.Pod{
	    TypeMeta: metav1.TypeMeta{
		Kind:       "Pod",
		APIVersion: "v1",
	    },
	    ObjectMeta: metav1.ObjectMeta{
		Name: "test-pod",
	    },
	    Spec: corev1.PodSpec{
		Containers: []corev1.Container{
		    {
			Name:  "nginx",
			Image: "nginx",
		    },
		},
	    },
	}
    
	_, err = k8s.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	require.NoError(t, err)

	testcontainers.CleanupContainer(t, ctr)
}
