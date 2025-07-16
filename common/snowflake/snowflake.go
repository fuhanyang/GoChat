package snowflake

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
	"user/Const"
	"user/settings"
)

var node *snowflake.Node

func Init(config *settings.SnowflakeConfig) (err error) {
	var st time.Time
	if st, err = time.Parse(time.RFC3339, config.StartTime); err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(int64(Const.Machine_code))
	if node == nil {
		return fmt.Errorf("failed to create snowflake node")
	}
	return
}
func GetID() int64 {
	return node.Generate().Int64()
}
