package utils

import "github.com/bwmarrin/snowflake"

// defaultNode 是全局的 Snowflake 节点实例，用于生成唯一ID。
var defaultNode *snowflake.Node

// InitDefaultSnowflakeNode 初始化全局的 Snowflake 节点实例。
// 参数:
//
//	node: 节点 ID，必须在 0 到 1023 之间。
//
// 异常:
//
//	如果节点初始化失败，将抛出 panic。
func InitDefaultSnowflakeNode(node int64) {
	var err error
	defaultNode, err = snowflake.NewNode(node)
	if err != nil {
		panic(err)
	}
}

// GetSnowflakeId 生成并返回一个唯一的 Snowflake ID。
// 返回值:
//
//	一个唯一的 64 位整数 ID。
func GetSnowflakeId() int64 {
	return defaultNode.Generate().Int64()
}
