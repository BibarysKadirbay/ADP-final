package telemetry

import (
	"context"
	"log"
)

func StartSpan(
	ctx context.Context,
	name string,
) context.Context {

	log.Println("[TRACE START]", name)

	return ctx
}

func EndSpan(name string) {
	log.Println("[TRACE END]", name)
}
