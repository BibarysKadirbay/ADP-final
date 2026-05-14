package integration

import (
	"testing"
)

func TestRepositoryIntegrationRequiresDocker(t *testing.T) {
	t.Skip("run in CI with delivery postgres and migrations applied")
}
