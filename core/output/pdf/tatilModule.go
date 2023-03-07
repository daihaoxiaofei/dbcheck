package pdf

import (
	"fmt"
)

// 设置标题
func (o *OutputWay) titleModule(titleStr string) {
	o.pdf.SetTitle(titleStr, false)
	o.pdf.SetHeaderFuncMode(func() {
		o.pdf.SetFont("simfang", "", 20)
		wd := o.pdf.GetStringWidth(titleStr) + 6
		o.pdf.SetY(3)                   // 先要设置 Y，然后再设置 X。否则，会导致 X 失效
		o.pdf.SetX((210 - wd) / 2)      // 水平居中的算法
		o.pdf.SetDrawColor(0, 80, 180)  // frame color
		o.pdf.SetFillColor(230, 230, 0) // background color
		o.pdf.SetTextColor(220, 50, 50) // text color
		o.pdf.SetLineWidth(1)
		o.pdf.CellFormat(wd, 10, titleStr, "1", 50, "CM", true, 0, "") // 第 5 个参数，实际效果是:指定下一行的位置
		o.pdf.Ln(1)
	}, false)
}

// 设置一级标题
func (o *OutputWay) primaryTitleModule(titleStr string) {
	o.pdf.SetFont("simfang", "", 12)
	o.pdf.SetFillColor(200, 220, 255) // background color
	o.pdf.CellFormat(0, 6, titleStr, "", 1, "L", true, 0, "")
	o.pdf.Ln(2)
}

// 设置页眉页脚
func (o *OutputWay) headerModule() {
	o.pdf.SetFooterFunc(func() {
		o.pdf.SetY(-15)
		o.pdf.SetFont("simfang", "", 8)
		o.pdf.SetTextColor(128, 128, 128)
		o.pdf.CellFormat(
			0, 5,
			fmt.Sprintf("Page %d", o.pdf.PageNo()),
			"", 0, "C", false, 0, "",
		)
	})
}
