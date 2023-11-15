package pshealthz

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	errorz "github.com/kunitsucom/util.go/errors"
	cliz "github.com/kunitsucom/util.go/exp/cli"
	signalz "github.com/kunitsucom/util.go/os/signal"
	"github.com/kunitsucom/util.go/version"

	"github.com/kunitsucom/pshealthz/internal/config"
	"github.com/kunitsucom/pshealthz/internal/pshealthz"
)

//nolint:cyclop
func PSHealthz(ctx context.Context) error {
	if _, err := config.Load(ctx); err != nil {
		if errors.Is(err, cliz.ErrHelp) {
			return nil
		}
		return fmt.Errorf("config.Load: %w", err)
	}

	if config.Version() {
		fmt.Printf("version: %s\n", version.Version())           //nolint:forbidigo
		fmt.Printf("revision: %s\n", version.Revision())         //nolint:forbidigo
		fmt.Printf("build branch: %s\n", version.Branch())       //nolint:forbidigo
		fmt.Printf("build timestamp: %s\n", version.Timestamp()) //nolint:forbidigo
		return nil
	}

	server := &http.Server{
		Addr:              config.Addr(),
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           http.HandlerFunc(pshealthz.PSHealthz),
	}

	ctx, stop := signalz.NotifyContext(ctx, func(signal os.Signal, stop context.CancelCauseFunc) {
		defer stop(fmt.Errorf("signal received: %s", signal.String())) //nolint:goerr113 // signal は error interface を実装していない
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			err = errorz.Errorf("server.Shutdown: %w", err)
			panic(err)
		}
	}, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer stop(nil)

	errChan := make(chan error, 1)
	go func(errChan chan<- error) {
		errChan <- server.ListenAndServe()
	}(errChan)

	<-ctx.Done()

	if err := <-errChan; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errorz.Errorf("server.ListenAndServe: %w", err)
	}

	return nil
}
