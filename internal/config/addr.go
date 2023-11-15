package config

import (
	"context"

	cliz "github.com/kunitsucom/util.go/exp/cli"
)

func loadAddr(_ context.Context, cmd *cliz.Command) string {
	v, _ := cmd.GetOptionString(_OptionAddr)
	return v
}

func Addr() string {
	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig.Addr
}
