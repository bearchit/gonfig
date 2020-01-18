package gonfig_test

import (
	"github.com/bearchit/gonfig"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestNewEngine(t *testing.T) {
	e := gonfig.New()
	assert.NotNil(t, e)
}

func TestEngine_Env(t *testing.T) {
	c := new(struct {
		Name string
	})

	t.Run("env", func(t *testing.T) {
		os.Setenv("NAME", "gonfig")
		e := gonfig.New(
			gonfig.WithScanners(
				gonfig.NewEnvScanner("", false),
			),
		)

		require.NoError(t, e.Unmarshal(c))
		assert.Equal(t, "gonfig", c.Name)
	})
}
