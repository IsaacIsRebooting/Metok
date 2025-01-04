package filerepohelper

type FileTableShardingConfig interface {
	GetShardingNumber(tableName, domainName, bizName string) int64
}
