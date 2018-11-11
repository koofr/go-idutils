// The idutils package generates unique IDs
//
// Implementation is based on Twitter Snowflake
// https://github.com/twitter/snowflake/blob/snowflake-2010/src/main/scala/com/twitter/service/snowflake/IdWorker.scala
package idutils

import (
	"fmt"
	"sync"
	"time"
)

const (
	CustomEpoch        = int64(1288834974657) // Thu, 04 Nov 2010 01:42:54 GMT
	WorkerIdBits       = 5
	DatacenterIdBits   = 5
	MaxWorkerId        = -1 ^ (-1 << 5) // -1 ^ (-1 << workerIdBits)
	MaxDatacenterId    = -1 ^ (-1 << 5) // -1 ^ (-1 << datacenterIdBits)
	SequenceBits       = 12
	WorkerIdShift      = 12              // sequenceBits
	DatacenterIdShift  = 17              // sequenceBits + workerIdBits
	TimestampLeftShift = 22              // sequenceBits + workerIdBits + datacenterIdBits
	SequenceMask       = -1 ^ (-1 << 12) // -1 ^ (-1 << sequenceBits)
)

type TimeId = int64
type Timestamp = int64

func timeGen() int64 {
	return time.Now().UnixNano() / 1000000
}

type Generator struct {
	WorkerId      int64
	DatacenterId  int64
	sequence      int64
	lastTimestamp Timestamp
	timeGen       func() Timestamp
	mutex         sync.Mutex
}

func NewGenerator(workerId int64, datacenterId int64) (*Generator, error) {
	if workerId > MaxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("worker Id cannot be greater than %d or less than 0", MaxWorkerId)
	}

	if datacenterId > MaxDatacenterId || datacenterId < 0 {
		return nil, fmt.Errorf("datacenter Id cannot be greater than %d or less than 0", MaxDatacenterId)
	}

	g := &Generator{
		WorkerId:      workerId,
		DatacenterId:  datacenterId,
		sequence:      0,
		lastTimestamp: -1,
		timeGen:       timeGen,
	}

	return g, nil
}

func (g *Generator) buildId() TimeId {
	return ((g.lastTimestamp - CustomEpoch) << TimestampLeftShift) |
		(g.DatacenterId << DatacenterIdShift) |
		(g.WorkerId << WorkerIdShift) |
		g.sequence
}

func (g *Generator) NextId() (TimeId, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	timestamp := g.timeGen()

	if timestamp < g.lastTimestamp {
		return 0, fmt.Errorf("Clock moved backwards. Refusing to generate id for %d milliseconds", g.lastTimestamp-timestamp)
	}

	if g.lastTimestamp == timestamp {
		g.sequence = (g.sequence + 1) & SequenceMask

		if g.sequence == 0 {
			timestamp = g.tilNextMillis(g.lastTimestamp)
		}
	} else {
		g.sequence = 0
	}

	g.lastTimestamp = timestamp

	return g.buildId(), nil
}

func (g *Generator) tilNextMillis(lastTimestamp Timestamp) Timestamp {
	timestamp := g.timeGen()

	for timestamp <= lastTimestamp {
		timestamp = g.timeGen()
	}

	return timestamp
}

func IdToTimestamp(id TimeId) Timestamp {
	return (id >> TimestampLeftShift) + CustomEpoch
}

func IdToTime(id TimeId) time.Time {
	return time.Unix(0, IdToTimestamp(id)*1000000).UTC()
}

func IdEndOfTimestamp(timestamp Timestamp) TimeId {
	return (timestamp-CustomEpoch)<<TimestampLeftShift | ((1 << TimestampLeftShift) - 1)
}

func IdEndOfTime(t time.Time) TimeId {
	return IdEndOfTimestamp(t.UnixNano() / 1000000)
}

func IdStartOfTimestamp(timestamp Timestamp) TimeId {
	return (timestamp - CustomEpoch) << TimestampLeftShift
}

func IdStartOfTime(t time.Time) TimeId {
	return IdStartOfTimestamp(t.UnixNano() / 1000000)
}
