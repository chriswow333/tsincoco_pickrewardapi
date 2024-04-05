package ctx

import (
	"context"
	"sync/atomic"
	"time"

	"pickrewardapi/base/pointer"
	"pickrewardapi/base/validator"

	"github.com/sirupsen/logrus"
)

var (
	// logExcept set keys that should be except from log
	logExcept = []string{"accessToken"}
)

// CTX extends Google's context to support logging methods.
type CTX struct {
	context.Context
	logrus.FieldLogger
	// Following counts are used for counting redis requests per http request
	RedisCacheCount   *int32
	RedisPersistCount *int32
}

// Background returns a non-nil, empty Context. It is never canceled, has no values, and
// has no deadline. It is typically used by the main function, initialization, and tests,
// and as the top-level Context for incoming requests
func Background() CTX {
	return CTX{
		Context:           context.Background(),
		FieldLogger:       logrus.StandardLogger(),
		RedisCacheCount:   pointer.Int32(0),
		RedisPersistCount: pointer.Int32(0),
	}
}

// TODO returns a Background.
func TODO() CTX {
	return Background()
}

// WithValue returns a copy of parent in which the value associated with key is val.
func WithValue(parent CTX, key string, val interface{}) CTX {
	if validator.IsInStringSlice(logExcept, key) {
		return CTX{
			Context:           context.WithValue(parent, key, val),
			FieldLogger:       parent.FieldLogger,
			RedisCacheCount:   parent.RedisCacheCount,
			RedisPersistCount: parent.RedisPersistCount,
		}
	}
	return CTX{
		Context:           context.WithValue(parent, key, val),
		FieldLogger:       parent.FieldLogger.WithField(key, val),
		RedisCacheCount:   parent.RedisCacheCount,
		RedisPersistCount: parent.RedisPersistCount,
	}
}

// WithValues returns a copy of parent in which the values associated with keys are vals.
func WithValues(parent CTX, kvs map[string]interface{}) CTX {
	c := parent
	for k, v := range kvs {
		c = WithValue(c, k, v)
	}
	return c
}

// WithCancel returns a copy of parent with added cancel function
func WithCancel(parent CTX) (CTX, context.CancelFunc) {
	newCtx, cFunc := context.WithCancel(parent)
	return CTX{
		Context:           newCtx,
		FieldLogger:       parent.FieldLogger,
		RedisCacheCount:   parent.RedisCacheCount,
		RedisPersistCount: parent.RedisPersistCount,
	}, cFunc
}

// WithTimeout returns a copy of parent with timeout condition
// and cancel function
func WithTimeout(parent CTX, d time.Duration) (CTX, context.CancelFunc) {
	newCtx, cFunc := context.WithTimeout(parent, d)
	return CTX{
		Context:           newCtx,
		FieldLogger:       parent.FieldLogger,
		RedisCacheCount:   parent.RedisCacheCount,
		RedisPersistCount: parent.RedisPersistCount,
	}, cFunc
}

// IncrRedisCount add cnt on ctx.RedisCacheCount or ctx.RedisPersistCount
func (c CTX) IncrRedisCount(name string, cnt int32) {
	switch name {
	case "cache":
		if c.RedisCacheCount != nil {
			atomic.AddInt32(c.RedisCacheCount, cnt)
		}
	case "persistent":
		if c.RedisPersistCount != nil {
			atomic.AddInt32(c.RedisPersistCount, cnt)
		}
	}
}

// LoadRedisCacheCount load ctx.RedisCacheCount if exists
func (c CTX) LoadRedisCacheCount() (int32, bool) {
	if c.RedisCacheCount == nil {
		return 0, false
	}
	return atomic.LoadInt32(c.RedisCacheCount), true
}

// LoadRedisPersistCount load ctx.RedisPersistCount if exists
func (c CTX) LoadRedisPersistCount() (int32, bool) {
	if c.RedisPersistCount == nil {
		return 0, false
	}
	return atomic.LoadInt32(c.RedisPersistCount), true
}
