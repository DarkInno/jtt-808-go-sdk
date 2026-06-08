package core

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
)

// ShardedConnPool 分片连接池
type ShardedConnPool struct {
	shards     []*connShard
	shardCount int
	totalCount int64
}

// connShard 连接分片
type connShard struct {
	mu    sync.RWMutex
	conns map[string]Connection
	count int
}

// NewShardedConnPool 创建分片连接池
func NewShardedConnPool(shardCount int) *ShardedConnPool {
	if shardCount <= 0 {
		shardCount = 16
	}

	pool := &ShardedConnPool{
		shards:     make([]*connShard, shardCount),
		shardCount: shardCount,
	}

	for i := 0; i < shardCount; i++ {
		pool.shards[i] = &connShard{
			conns: make(map[string]Connection),
		}
	}

	return pool
}

// getShard 获取分片
func (p *ShardedConnPool) getShard(key string) *connShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return p.shards[h.Sum32()%uint32(p.shardCount)]
}

// Put 添加连接
func (p *ShardedConnPool) Put(conn Connection) {
	shard := p.getShard(conn.DeviceID())
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.conns[conn.DeviceID()]; !exists {
		shard.count++
		atomic.AddInt64(&p.totalCount, 1)
	}
	shard.conns[conn.DeviceID()] = conn
}

// Get 获取连接
func (p *ShardedConnPool) Get(deviceID string) (Connection, bool) {
	shard := p.getShard(deviceID)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	conn, exists := shard.conns[deviceID]
	return conn, exists
}

// Delete 删除连接
func (p *ShardedConnPool) Delete(deviceID string) {
	shard := p.getShard(deviceID)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.conns[deviceID]; exists {
		delete(shard.conns, deviceID)
		shard.count--
		atomic.AddInt64(&p.totalCount, -1)
	}
}

// Size 获取连接数
func (p *ShardedConnPool) Size() int64 {
	return atomic.LoadInt64(&p.totalCount)
}

// Range 遍历连接
func (p *ShardedConnPool) Range(fn func(deviceID string, conn Connection) bool) {
	for _, shard := range p.shards {
		shard.mu.RLock()
		for deviceID, conn := range shard.conns {
			if !fn(deviceID, conn) {
				shard.mu.RUnlock()
				return
			}
		}
		shard.mu.RUnlock()
	}
}

// Clear 清空连接池
func (p *ShardedConnPool) Clear() {
	for _, shard := range p.shards {
		shard.mu.Lock()
		shard.conns = make(map[string]Connection)
		shard.count = 0
		shard.mu.Unlock()
	}
	atomic.StoreInt64(&p.totalCount, 0)
}
