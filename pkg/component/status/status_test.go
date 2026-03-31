// SPDX-FileCopyrightText: 2026 k0s authors
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"context"
	"errors"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubefake "k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusHandlerProbesOwnNode(t *testing.T) {
	t.Parallel()

	client := kubefake.NewSimpleClientset(&corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "worker0"}})
	client.PrependReactor("list", "nodes", func(ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("unexpected list call")
	})

	handler := statusHandler{
		Status: &Status{
			StatusInformation: K0sStatus{
				Workloads: true,
				NodeName:  "worker0",
			},
		},
		client: client,
	}

	status := handler.getCurrentStatus(context.Background())
	assert.True(t, status.WorkerToAPIConnectionStatus.Success)
	assert.Empty(t, status.WorkerToAPIConnectionStatus.Message)
}

func TestStatusHandlerReportsProbeErrors(t *testing.T) {
	t.Parallel()

	handler := statusHandler{
		Status: &Status{
			StatusInformation: K0sStatus{
				Workloads: true,
				NodeName:  "worker0",
			},
		},
		client: kubefake.NewSimpleClientset(),
	}

	status := handler.getCurrentStatus(context.Background())
	assert.False(t, status.WorkerToAPIConnectionStatus.Success)
	assert.ErrorContains(t, errors.New(status.WorkerToAPIConnectionStatus.Message), `nodes "worker0" not found`)
}
