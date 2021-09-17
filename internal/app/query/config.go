package query

import (
	"time"

	"github.com/romberli/go-util/constant"
)

const (
	defaultLimit    = 5
	defaultOffset   = 0
	minRowsExamined = 100000

	maxDuration = 30 * constant.Day
	maxLimit    = 100
)

type OrderType int

type Config struct {
	startTime time.Time
	endTime   time.Time
	limit     int
	offset    int
}

func NewConfig(startTime, endTime time.Time, limit, offset int) *Config {
	return newConfig(startTime, endTime, limit, offset)
}

func NewConfigWithDefault() *Config {
	return newConfig(time.Now().Add(-constant.Week), time.Now(), defaultLimit, defaultOffset)
}

func newConfig(startTime, endTime time.Time, limit, offset int) *Config {
	return &Config{
		startTime: startTime,
		endTime:   endTime,
		limit:     limit,
		offset:    offset,
	}
}

func (c *Config) GetStartTime() time.Time {
	return c.startTime
}

func (c *Config) GetEndTime() time.Time {
	return c.endTime
}

func (c *Config) GetLimit() int {
	return c.limit
}

func (c *Config) GetOffset() int {
	return c.offset
}

func (c *Config) SetStartTime(startTime time.Time) {
	c.startTime = startTime
}

func (c *Config) SetEndTime(endTime time.Time) {
	c.endTime = endTime
}

func (c *Config) SetLimit(limit int) {
	c.limit = limit
}

func (c *Config) SetOffset(offset int) {
	c.offset = offset
}

func (c *Config) IsValid() bool {
	duration := c.GetEndTime().Sub(c.GetStartTime())
	if duration > maxDuration {
		return false
	}

	if c.GetLimit() > maxLimit {
		return false
	}

	return true
}
