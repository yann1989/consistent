// Author: yann
// Date: 2022/5/28
// Desc: consistent

package consistent

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

//实现sort.sort接口
type uints []uint32

// Len returns the length of the uints array.
func (x uints) Len() int { return len(x) }

// Less returns true if element i is less than element j.
func (x uints) Less(i, j int) bool { return x[i] < x[j] }

// Swap exchanges elements i and j.
func (x uints) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

type Hash func(string) uint32

type Consistent struct {
	hash     Hash              // hash算法
	replicas int               // 虚拟节点
	sortKeys uints             // 已排序的节点哈希切片
	hashMap  map[uint32]string // 节点哈希和KEY的map，键是哈希值，值是节点Key
	sync.RWMutex
}

const (
	replicasNone = iota << 7
	defaultReplicas
)

var (
	// 默认使用CRC32算法
	defaultHash = func(key string) uint32 {
		return crc32.ChecksumIEEE([]byte(key))
	}
)

func New(opts ...Option) *Consistent {
	m := &Consistent{
		hashMap: make(map[uint32]string),
	}
	for _, opt := range opts {
		opt(m)
	}

	if m.replicas == replicasNone {
		m.replicas = defaultReplicas
	}
	if m.hash == nil {
		m.hash = defaultHash
	}
	return m
}

func (c *Consistent) IsEmpty() bool {
	return len(c.sortKeys) == 0
}

// IsExist 判断节点是否已添加
func (c *Consistent) IsExist(node string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.hashMap[c.hash(c.generateKey(node, 0))]
	return ok
}

// Add 方法用来添加缓存节点，参数为节点key，比如使用IP
func (c *Consistent) Add(nodes ...string) {
	c.Lock()
	defer c.Unlock()
	c.add(nodes...)
}

// 调用前请加锁
func (c *Consistent) add(nodes ...string) {
	for _, node := range nodes {
		// 结合复制因子计算所有虚拟节点的hash值，并存入m.keys中，同时在m.hashMap中保存哈希值和key的映射
		for i := 0; i < c.replicas; i++ {
			hash := c.hash(c.generateKey(node, i))
			c.sortKeys = append(c.sortKeys, hash)
			c.hashMap[hash] = node
		}
	}
	// 对所有虚拟节点的哈希值进行排序，方便之后进行二分查找
	sort.Sort(c.sortKeys)
}

// generateKey generates a string key for a node with an index.
func (c *Consistent) generateKey(node string, index int) string {
	return fmt.Sprintf("%s%d", node, index)
}

// Get 获取离给定对象最近的节点hash
func (c *Consistent) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	if c.IsEmpty() {
		return "", false
	}

	hash := c.hash(key)

	// 通过二分查找第一个节点hash值大于对象hash值的节点
	idx := sort.Search(len(c.sortKeys), func(i int) bool { return c.sortKeys[i] >= hash })

	// 如果查找结果大于节点哈希数组的最大索引，则为第一节点
	if idx == len(c.sortKeys) {
		idx = 0
	}

	return c.hashMap[c.sortKeys[idx]], true
}

// Remove 删除一个节点
func (c *Consistent) Remove(node string) {
	c.Lock()
	defer c.Unlock()
	c.remove(node)
}

// 需要在调用前加锁
func (c *Consistent) remove(node string) {
	for i := 0; i < c.replicas; i++ {
		delete(c.hashMap, c.hash(c.generateKey(node, i)))
	}
	c.updateKeys()
}

// 需要在调用前加锁 在删除后重新排序
func (c *Consistent) updateKeys() {
	var sortKeys uints
	for k := range c.hashMap {
		sortKeys = append(sortKeys, k)
	}
	sort.Sort(sortKeys)
	c.sortKeys = sortKeys
}
