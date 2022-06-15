// Author: yann
// Date: 2022/5/28
// Desc: consistent

package consistent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	consistent := New()
	assert.NotEqual(t, consistent, nil, "实例为空")
}

func TestConsistent_Add(t *testing.T) {
	consistent := New()
	assert.NotEqual(t, consistent, nil, "实例为空")
	consistent.add("192.168.10.13", "192.168.10.14", "192.168.10.15", "192.168.10.16", "192.168.10.17", "192.168.10.18")

	users := []string{"xxx1", "TpJ-0-1jl/ab2b642468b18b7d.hblock", "TpJ-0-1jl/aa299567d32e93de.hblock", "TpJ-0-1jl/ac29956xx32e93de.hblock", "TpJ-0-1jl/ad299567d32xx3de.hblock", "TpJ-0-1jl/ae299567d32e93de.hblock"}
	for _, key := range users {
		get := consistent.Get(key)
		t.Logf("节点: %v   key: %v", get, key)
	}

}
