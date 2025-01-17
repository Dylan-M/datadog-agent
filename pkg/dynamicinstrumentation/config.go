// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package dynamicinstrumentation

import (
	"github.com/DataDog/datadog-agent/cmd/system-probe/config"
	"github.com/DataDog/datadog-agent/pkg/ebpf"
)

// Config holds the configuration for the user tracer system probe module
type Config struct {
	ebpf.Config
	DynamicInstrumentationEnabled bool
}

//nolint:revive // TODO(DEBUG) Fix revive linter
func NewConfig(sysprobeConfig *config.Config) (*Config, error) {
	_, diEnabled := sysprobeConfig.EnabledModules[config.DynamicInstrumentationModule]
	return &Config{
		Config:                        *ebpf.NewConfig(),
		DynamicInstrumentationEnabled: diEnabled,
	}, nil
}
