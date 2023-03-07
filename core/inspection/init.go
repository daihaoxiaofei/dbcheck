package inspection

func Check() {
	DatabaseConfigCheck()                 // 配置参数检查功能
	BaselineCheckIndexColumnDesign()      // 索引设计合规性
	tableNoPrimaryKey()                   // 数据库的基线检查功能--检查表设计合规性 检查表字符集是否为utf8
	InformationSchemaKeyColumnUsage()     // 检测是否存在外键约束
	BaselineCheckProcedureTriggerDesign() // 存储过程、存储函数、触发器检查限制
	UserPrivileges()                      // 用户权限检查
	BaselineCheckPortDesign()             // 开始检查当前数据库是否使用默认端口3306
	DatabasePerformanceStatusCheck()      // 磁盘使用量
	DatabasePerformanceTableIndexCheck()  // 表字符集检查 utf8mb4...等
}
