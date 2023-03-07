package result

import (
	"time"
)

type Rs struct {
	VariableName string // 项目
	Variable     string // 当前值
	Value        string // 建议值
	Status       string // 正常 normal 异常 abnormal
	Type         string // 状态
	Threshold    string // 阈值
	CurrentValue string // 现行的值
}

// InspectionResults 记录检测结果
type InspectionResults struct {
	// 数据库配置参数检查
	Config struct {
		ConfigParameter []map[string]string
	}
	// 数据库基线检查
	Baseline struct {
		TableDesign struct {
			TableCharset      []map[string]string
			TableEngine       []map[string]string
			TableForeign      []map[string]string
			TableNoPrimaryKey []map[string]string
		}
		ColumnDesign struct {
			TableAutoIncrement []map[string]string
			TableBigColumns    []map[string]string
		}
		IndexColumnsDesign struct {
			IndexColumnType          []map[string]string
			IndexColumnIsNull        []map[string]string
			IndexColumnIsRepeatIndex []map[string]string
		}
		ProcedureTriggerDesign struct {
			TableProcedure []map[string]string
			TableTrigger   []map[string]string
			TableFunc      []map[string]string
		}
	}
	// 数据库性能检查
	Performance struct {
		PerformanceStatus struct {
			BinlogDiskUsageRate              []map[string]string
			HistoryConnectionMaxUsageRate    []map[string]string
			TmpDiskTableUsageRate            []map[string]string
			TmpDiskfileUsageRate             []map[string]string
			InnodbBufferPoolUsageRate        []map[string]string
			InnodbBufferPoolDirtyPagesRate   []map[string]string
			InnodbBufferPoolHitRate          []map[string]string
			OpenFileUsageRate                []map[string]string
			OpenTableCacheUsageRate          []map[string]string
			OpenTableCacheOverflowsUsageRate []map[string]string
			SelectScanUsageRate              []map[string]string
			SelectfullJoinScanUsageRate      []map[string]string
		} // 检查状态
		PerformanceTableIndex struct {
			TableAutoPrimaryKeyUsageRate []map[string]string
			TableRows                    []map[string]string
			DiskFragmentationRate        []map[string]string
			BigTable                     []map[string]string
			ColdTable                    []map[string]string
		}
		// 检查表和索引
	}
	// 数据库安全检查
	Security struct {
		UserPriDesign struct {
			AnonymousUsers                  []map[string]string
			EmptyPasswordUser               []map[string]string
			RootUserRemoteLogin             []map[string]string
			NormalUserConnectionUnlimited   []map[string]string
			UserPasswordSame                []map[string]string
			NormalUserDatabaseAllPrivilages []map[string]string
			NormalUserSuperPrivilages       []map[string]string
		}
		PortDesign struct {
			DatabasePort []map[string]string
		}
	}
}

// 检测耗时相关变量
var (
	BeginTime     time.Time              // 开始时间
	ConsumingTime float64                // 耗时
	R             = &InspectionResults{} // 一级结构体初始化
)
