package test

import (
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplates(t *testing.T) {
	t.Parallel()

	chartPath, err := filepath.Abs("../")
	require.NoError(t, err)

	suite.Run(t, &TemplateTest{
		ChartPath: chartPath,
		Release:   "ld-relay-test",
		Namespace: "ld-relay-" + strings.ToLower(random.UniqueId()),
	})
}
