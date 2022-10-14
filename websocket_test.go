package gopty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRmrf(t *testing.T) {
	matches := dangerRmrf.FindAll([]byte("rm -rf .\n"), -1)
	assert.Equal(t, 1, len(matches))
	matched := dangerRmrf.MatchString("rm -rf /\n")
	assert.True(t, matched)

	tips := make([]byte, len(dangerCmdTips), len(dangerCmdTips)+len(matches[0]))
	copy(tips, dangerCmdTips)
	tips = append(tips, matches[0]...)
	assert.Equal(t, "Dangerous commands disabled: rm -rf .\n", string(tips))
}
