package api

import (
	"fmt"
	"github.com/yunhanshu-net/pkg/typex"
	"github.com/yunhanshu-net/pkg/typex/files"
	"github.com/yunhanshu-net/pkg/x/jsonx"
	"testing"
	"time"
)

func Test1111(t *testing.T) {
	s := jsonx.String(files.Files{
		Files: []*files.File{
			{
				CreatedAt: typex.NewTime(time.Now()),
			},
		},
		Options:   map[string]interface{}{},
		Config:    map[string]interface{}{},
		CreatedAt: typex.Time(time.Now()),
	})
	fmt.Println(s)
}
