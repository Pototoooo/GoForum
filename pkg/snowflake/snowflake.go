package snowflake

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	if st, err = time.Parse("2006-01-02", startTime); err != nil {
		return err
	}
	snowflake.Epoch = st.UnixMilli()
	Node, err = snowflake.NewNode(machineID)
	return err
}
func GenerateID() (id int64) {
	return Node.Generate().Int64()
}
func main() {
	// 初始化雪花算法节点
	Init("2023-01-01", 1)
	fmt.Println(GenerateID())
}
