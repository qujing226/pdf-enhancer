package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 模拟雪花ID生成器（简化版）
// 格式：timestamp (毫秒) + nodeID (10位以内) + sequence (0-999)
const nodeBits = 10                           // 可支持最多 2^10 = 1024 个节点
const sequenceBits = 10                       // 每毫秒最多支持 2^10 = 1024 个ID
const maxSequence = -1 ^ (-1 << sequenceBits) // = 1023

var nodeId int64
var lastTimestamp int64
var sequence int64
var mu sync.Mutex

func init() {
	rand.Seed(time.Now().UnixNano())
	nodeId = int64(rand.Intn(1024)) // 随机分配一个节点ID
}

// GenerateSnowflakeID 生成模拟的雪花风格ID
func GenerateSnowflakeID() string {
	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	if timestamp < lastTimestamp {
		// 时钟回拨
		panic(fmt.Sprintf("时钟回拨 %d 毫秒", lastTimestamp-timestamp))
	}

	if timestamp == lastTimestamp {
		sequence = (sequence + 1) % maxSequence
		if sequence == 0 {
			// 等待下一毫秒
			timestamp = tilNextMillis(lastTimestamp)
		}
	} else {
		sequence = 0
	}

	lastTimestamp = timestamp

	id := (timestamp << (nodeBits + sequenceBits)) |
		(nodeId << sequenceBits) |
		sequence

	return fmt.Sprintf("%d", id)
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	for timestamp <= lastTimestamp {
		time.Sleep(1 * time.Millisecond)
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}
