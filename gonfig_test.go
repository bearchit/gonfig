package gonfig_test

import (
	"os"
	"testing"

	"github.com/bearchit/gonfig"

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

func TestEnvScanner_ErrorOnFail(t *testing.T) {
	newGonfig := func(breakOnError bool) *gonfig.Engine {
		return gonfig.New(
			gonfig.WithScanners(
				gonfig.NewYMLScanner("config/test.env", breakOnError)),
		)
	}

	t.Run("skip on error", func(t *testing.T) {
		gfg := newGonfig(false)
		assert.NoError(t, gfg.Unmarshal(&struct{}{}))
	})

	t.Run("break on error", func(t *testing.T) {
		gfg := newGonfig(true)
		assert.Error(t, gfg.Unmarshal(&struct{}{}))
	})
}
