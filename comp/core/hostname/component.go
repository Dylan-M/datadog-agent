// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

// Package hostname exposes hostname.Get() as a component.
package hostname

import (
	"context"

	"github.com/DataDog/datadog-agent/pkg/util/hostname"
)

// team: agent-shared-components

// Component is the component type.
type Component interface {
	// Get returns the host name for the agent.
	Get(context.Context) (string, error)
	// GetWithProvider returns the hostname for the Agent and the provider that was use to retrieve it.
	GetWithProvider(ctx context.Context) (hostname.Data, error)
	// GetSafe is Get(), but it returns 'unknown host' if anything goes wrong.
	GetSafe(context.Context) string
}
