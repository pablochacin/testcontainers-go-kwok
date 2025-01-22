package kwok_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kwok "github.com/pablochacin/testcontainers-go-kwok"
	"github.com/testcontainers/testcontainers-go"
)

func TestKwok(t *testing.T) {
	ctx := context.Background()

	kwokContainer, err := kwok.Run(ctx, "ghcr.io/pablochacin/kwok:latest")
	require.NoError(t, err)
	testcontainers.CleanupContainer(t, kwokContainer)
    
	k8s, err := kwokContainer.GetClient(ctx)
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

	testcontainers.CleanupContainer(t, kwokContainer)
}
