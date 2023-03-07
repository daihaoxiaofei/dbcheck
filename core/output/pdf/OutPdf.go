package pdf

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"strconv"

	"dbcheck/core/result"
	"dbcheck/pkg/config"
)

func OutPdf() {
	var OutputPdf OutputWay
	OutputPdf.OutPdf()
}

type OutputWay struct {
	pdf  *gofpdf.Fpdf
	data map[string][]map[string]string
}

func (o *OutputWay) pdfInit() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("simfang", "", "./core/output/ttf/simfang.ttf")
	return pdf
}

func (o *OutputWay) OutPdf() {
	// 注入数据
	o.data = getData()

	// 设置页面参数
	var tmpCheckTypeSlice []string
	o.pdf = o.pdfInit()

	// // 标题头
	// o.titleModule(`MySQL 巡检报告`)

	o.pdf.AddPage()

	// 标题一
	o.primaryTitleModule("一、巡检介绍")
	title := []float64{50.0, 50.0, 40.0, 40.0}
	var dc1 = []string{"巡检时间:", result.BeginTime.Format("2006-01-02 15:04:05"), "巡检人员:", config.C.ResultOutput.InspectionPersonnel}
	var dc2 = []string{"巡检级别:", config.C.ResultOutput.InspectionLevel, "巡检耗时(s):", strconv.FormatFloat(result.ConsumingTime, 'f', -1, 64)}
	var dc = [][]string{dc1, dc2}
	o.tableInsert(title, dc)
	o.pdf.CellFormat(0, 2, "", "", 1, "LM", false, 0, "")
	o.pdf.Ln(-1)

	// 标题二
	o.primaryTitleModule("二、巡检结果概览")
	w := []float64{17.0, 70.0, 30.0, 38.0, 30.0} // 定义每行表格的宽度
	c := []string{"编号", "检测项", "检测数量", "正常", "异常"}
	o.tableBodyColorsFormat(w, c)
	o.tableInsert(w, o.ResultSummaryStringSlice())
	o.pdf.CellFormat(0, 2, "", "", 1, "LM", false, 0, "")
	o.pdf.Ln(-1)

	// 标题三
	o.primaryTitleModule("三、巡检结果详情")            // TNP
	w3 := []float64{10, 60.0, 35.0, 20.0, 70.0} // 定义每行表格的宽度
	cd := []string{" ", "巡检项名称", "阈值", "错误码", "异常相关信息"}
	o.pdf.SetFont("simfang", "", 10)
	o.pdf.MultiCell(0, 5, string("3.1 巡检数据库环境"), "", "", false)

	// 子标题3.2内容
	o.pdf.MultiCell(0, 5, string("3.2 巡检数据库配置"), "", "", false)
	o.tableBodyColorsFormat(w3, cd)
	// tmpResultConfig := o.tmpConfigCheckResultSummary("configParameter",result.R.Config.ConfigParameter)
	tmpCheckTypeSlice = []string{"configParameter"}
	tmpResultConfig := o.tmpResultSummary(tmpCheckTypeSlice)
	o.tableInsert(w3, tmpResultConfig)
	o.pdf.Ln(-1)

	// 子标题3.3内容
	o.pdf.MultiCell(0, 5, string("3.3 巡检数据库性能"), "", "", false)
	o.tableBodyColorsFormat(w3, cd)
	tmpCheckTypeSlice = []string{"binlogDiskUsageRate", "historyConnectionMaxUsageRate", "tmpDiskTableUsageRate",
		"tmpDiskfileUsageRate", "innodbBufferPoolUsageRate", "innodbBufferPoolDirtyPagesRate", "innodbBufferPoolHitRate",
		"openFileUsageRate", "openTableCacheUsageRate", "openTableCacheOverflowsUsageRate", "selectScanUsageRate", "selectfullJoinScanUsageRate",
		"tableAutoPrimaryKeyUsageRate", "tableRows", "diskFragmentationRate", "bigTable", "coldTable"}
	tmpResultPerformance := o.tmpResultSummary(tmpCheckTypeSlice)
	o.tableInsert(w3, tmpResultPerformance)
	o.pdf.Ln(-1)

	// 子标题3.4内容
	o.pdf.MultiCell(0, 5, string("3.4 巡检数据库基线"), "", "", false)
	o.tableBodyColorsFormat(w3, cd)
	tmpCheckTypeSlice = []string{"tableCharset", "tableEngine", "tableForeign", "tableNoPrimaryKey", "tableAutoIncrement",
		"tableBigColumns", "indexColumnIsNull", "indexColumnType", "tableIncludeRepeatIndex",
		"tableProcedureFunc", "tableTrigger"}
	tmpResultBaselineResult := o.tmpResultSummary(tmpCheckTypeSlice)
	o.tableInsert(w3, tmpResultBaselineResult)
	o.pdf.Ln(-1)

	// 子标题3.5内容
	o.pdf.MultiCell(0, 5, string("3.5 巡检数据库安全"), "", "", false)
	o.tableBodyColorsFormat(w3, cd)
	tmpCheckTypeSlice = []string{"anonymousUsers", "emptyPasswordUser", "rootUserRemoteLogin", "normalUserConnectionUnlimited",
		"userPasswordSame", "normalUserDatabaseAllPrivilages", "normalUserSuperPrivilages", "databasePort"}
	tmpResultUserSecurityResult := o.tmpResultSummary(tmpCheckTypeSlice)
	o.tableInsert(w3, tmpResultUserSecurityResult)
	o.pdf.Ln(-1)

	o.pdf.MultiCell(0, 5, string("3.6 巡检数据库空间"), "", "", false)
	o.pdf.MultiCell(0, 5, string("3.7 巡检数据库备份"), "", "", false)
	o.pdf.Ln(-1)

	// 将内容写入到pdf中
	if err := o.pdf.OutputFileAndClose(config.C.ResultOutput.OutputPath + config.C.ResultOutput.OutputFile); err != nil {
		panic(err.Error())
	}

}

func (o *OutputWay) ResultSummaryStringSlice() [][]string {
	var resultProfile [][]string

	var tempInt int
	for k, v := range o.data {
		if v != nil {
			tempInt++
			tmpRes := o.tempUnfold(fmt.Sprintf("%02d", tempInt), k, v)
			resultProfile = append(resultProfile, tmpRes)
		}
	}
	return resultProfile
}
