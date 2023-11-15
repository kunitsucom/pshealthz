package main

import (
	"context"
	"log"

	"github.com/kunitsucom/pshealthz/pkg/pshealthz"
)

func main() {
	ctx := context.Background()

	if err := pshealthz.PSHealthz(ctx); err != nil {
		log.Fatalf("pshealthz.PSHealthz: %+v", err)
	}
}
