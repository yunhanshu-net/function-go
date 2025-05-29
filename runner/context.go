package runner

import (
	"context"
	"fmt"
	"github.com/yunhanshu-net/pkg/constants"
)

type Context struct {
	context.Context
	user    string
	name    string
	version string
}

func (c *Context) getDb() string {
	return fmt.Sprintf("%s_%s.db", c.user, c.name)
}

func (c *Context) getTraceId() string {
	value := c.Context.Value(constants.TraceID)
	if value == nil {
		return ""
	}
	v, ok := value.(string)
	if ok {
		return v
	}
	return ""
}

func (c *Context) GetUsername() string {
	return ""
}

func (c *Context) GetFile() string {
	return ""
}
