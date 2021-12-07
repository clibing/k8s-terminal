package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	// load(filepath.Join(dir, execHome, fileName))
	DefLoad()
	t.Logf("current global config: %v", GlobalCfg)
}
