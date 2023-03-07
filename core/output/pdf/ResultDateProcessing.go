package pdf

import (
	"dbcheck/core/result"
	"fmt"
	"strconv"
)

// 将结果展开并统计
func (o *OutputWay) tempUnfold(checkNum, checkType string, ast []map[string]string) []string {
	var abnormalCount = 0
	var normalCount = 0
	var checkNumberTotal string
	var aa []string
	for k := range ast {
		if ast[k]["checkStatus"] == "abnormal" && ast[k]["checkType"] == checkType {
			abnormalCount++
		}
		if ast[k]["checkStatus"] == "normal" && ast[k]["checkType"] == checkType {
			normalCount++
		}
	}
	checkNumberTotal = strconv.Itoa(abnormalCount + normalCount)
	aa = []string{checkNum, checkType, checkNumberTotal, strconv.Itoa(normalCount), strconv.Itoa(abnormalCount)}
	return aa
}

func getData() map[string][]map[string]string {
	var data = make(map[string][]map[string]string)

	data["configParameter"] = result.R.Config.ConfigParameter
	data["tableCharset"] = result.R.Baseline.TableDesign.TableCharset
	data["tableEngine"] = result.R.Baseline.TableDesign.TableEngine
	data["tableNoPrimaryKey"] = result.R.Baseline.TableDesign.TableNoPrimaryKey
	data["tableForeign"] = result.R.Baseline.TableDesign.TableForeign
	data["tableAutoIncrement"] = result.R.Baseline.ColumnDesign.TableAutoIncrement
	data["tableBigColumns"] = result.R.Baseline.ColumnDesign.TableBigColumns
	data["indexColumnIsNull"] = result.R.Baseline.IndexColumnsDesign.IndexColumnIsNull
	data["indexColumnType"] = result.R.Baseline.IndexColumnsDesign.IndexColumnType
	data["tableIncludeRepeatIndex"] = result.R.Baseline.IndexColumnsDesign.IndexColumnIsRepeatIndex
	data["tableProcedureFunc"] = result.R.Baseline.ProcedureTriggerDesign.TableProcedure
	data["tableTrigger"] = result.R.Baseline.ProcedureTriggerDesign.TableTrigger
	data["anonymousUsers"] = result.R.Security.UserPriDesign.AnonymousUsers
	data["emptyPasswordUser"] = result.R.Security.UserPriDesign.EmptyPasswordUser
	data["rootUserRemoteLogin"] = result.R.Security.UserPriDesign.RootUserRemoteLogin
	data["normalUserConnectionUnlimited"] = result.R.Security.UserPriDesign.NormalUserConnectionUnlimited
	data["userPasswordSame"] = result.R.Security.UserPriDesign.UserPasswordSame
	data["normalUserDatabaseAllPrivilages"] = result.R.Security.UserPriDesign.NormalUserDatabaseAllPrivilages
	data["normalUserSuperPrivilages"] = result.R.Security.UserPriDesign.NormalUserSuperPrivilages
	data["databasePort"] = result.R.Security.PortDesign.DatabasePort
	data["binlogDiskUsageRate"] = result.R.Performance.PerformanceStatus.BinlogDiskUsageRate
	data["historyConnectionMaxUsageRate"] = result.R.Performance.PerformanceStatus.HistoryConnectionMaxUsageRate
	data["tmpDiskTableUsageRate"] = result.R.Performance.PerformanceStatus.TmpDiskTableUsageRate
	data["tmpDiskfileUsageRate"] = result.R.Performance.PerformanceStatus.TmpDiskfileUsageRate
	data["innodbBufferPoolUsageRate"] = result.R.Performance.PerformanceStatus.InnodbBufferPoolUsageRate
	data["innodbBufferPoolDirtyPagesRate"] = result.R.Performance.PerformanceStatus.InnodbBufferPoolDirtyPagesRate
	data["innodbBufferPoolHitRate"] = result.R.Performance.PerformanceStatus.InnodbBufferPoolHitRate
	data["openFileUsageRate"] = result.R.Performance.PerformanceStatus.OpenFileUsageRate
	data["openTableCacheUsageRate"] = result.R.Performance.PerformanceStatus.OpenTableCacheUsageRate
	data["openTableCacheOverflowsUsageRate"] = result.R.Performance.PerformanceStatus.OpenTableCacheOverflowsUsageRate
	data["selectScanUsageRate"] = result.R.Performance.PerformanceStatus.SelectScanUsageRate
	data["selectfullJoinScanUsageRate"] = result.R.Performance.PerformanceStatus.SelectfullJoinScanUsageRate
	data["tableAutoPrimaryKeyUsageRate"] = result.R.Performance.PerformanceTableIndex.TableAutoPrimaryKeyUsageRate
	data["tableRows"] = result.R.Performance.PerformanceTableIndex.TableRows
	data["diskFragmentationRate"] = result.R.Performance.PerformanceTableIndex.DiskFragmentationRate
	data["bigTable"] = result.R.Performance.PerformanceTableIndex.BigTable
	data["coldTable"] = result.R.Performance.PerformanceTableIndex.ColdTable
	return data
}

func (o *OutputWay) tmpcc(checkRulest []map[string]string) []string {
	var bc []string
	var tmpCheckType, tmpThreshold, tmpAbnormalInformation string
	var tmpeq int
	for i := range checkRulest {
		if checkRulest[i]["checkStatus"] == "abnormal" {
			tmpCheckType = checkRulest[i]["checkType"]
			tmpThreshold = checkRulest[i]["threshold"]
			tmpeq++
			if tmpeq > 1 {
				tmpAbnormalInformation = fmt.Sprintf("%s 等", tmpAbnormalInformation)
				break
			}
			if tmpAbnormalInformation != "" {
				tmpAbnormalInformation = fmt.Sprintf("%s,%s", tmpAbnormalInformation, checkRulest[i]["currentValue"])
			} else {
				tmpAbnormalInformation = fmt.Sprintf("%s", checkRulest[i]["currentValue"])
			}
		}
	}
	if tmpCheckType != "" && tmpThreshold != "" && tmpAbnormalInformation != "" {
		bc = []string{tmpCheckType, tmpThreshold, `ErrorCode`, tmpAbnormalInformation}
	}
	return bc
}

// 总结
func (o *OutputWay) tmpResultSummary(CheckTypeSlice []string) [][]string {
	var bc [][]string
	var cc []string

	var tempInt int
	for i := range CheckTypeSlice {
		if vi, ok := o.data[CheckTypeSlice[i]]; ok {
			if vi != nil {
				tempInt++
				cc = o.tmpcc(vi)
				if cc != nil {
					cd := append([]string{fmt.Sprintf("%02d", tempInt)}, cc...)
					bc = append(bc, cd)
				}
			}
		}
	}
	return bc
}
