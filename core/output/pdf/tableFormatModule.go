package pdf

import (
	"github.com/jung-kurt/gofpdf"
)

const colWd = 20 // 35.0

// 设置表格头部信息
func (o *OutputWay) pdfTableHeaderFormat(pdf *gofpdf.Fpdf, insertContent []string) *gofpdf.Fpdf {
	// 设置表格第一行的样式（也就是在设置 <thead> 标签的样式）
	pdf.SetTextColor(224, 224, 224)
	pdf.SetFillColor(64, 64, 64)
	// for colJ := 0; colJ <  len(insertContent); colJ++{
	pdf.CellFormat(colWd, 10, insertContent[0], "1", 0, "CM", true, 0, "")
	pdf.CellFormat(55.0, 10, insertContent[1], "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWd, 10, insertContent[2], "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWd, 10, insertContent[3], "1", 0, "CM", true, 0, "")
	// }
	pdf.Ln(-1)
	return pdf
}

// 设置表格内容信息
func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}

type countryType struct {
	nameStr, capitalStr, areaStr, popStr string
}

// Colored table
func (o *OutputWay) tableBodyColorsFormat(w []float64, insertContent []string) {
	countryList := make([]countryType, 0, 8)
	// Colors, line width and bold font
	o.pdf.SetFillColor(24, 24, 24)
	o.pdf.SetTextColor(255, 255, 255)
	// pdf.SetDrawColor(128, 0, 0)
	o.pdf.SetLineWidth(.3)
	o.pdf.SetFont("simfang", "", 10)
	// 	Header
	// w := []float64{17.0, 70.0, 30.0, 38.0,30}
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}
	left := (210 - wSum) / 2
	o.pdf.SetX(left)
	for j, str := range insertContent {
		o.pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
	}
	o.pdf.Ln(-1)
	// Color and font restoration 颜色和字体恢复
	o.pdf.SetFillColor(224, 235, 255)
	o.pdf.SetTextColor(0, 0, 0)
	o.pdf.SetFont("", "", 0)
	// 	Data
	fill := false
	for _, c := range countryList {
		o.pdf.SetX(left)
		o.pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", fill, 0, "")
		o.pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", fill, 0, "")
		o.pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
			"LR", 0, "R", fill, 0, "")
		o.pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
			"LR", 0, "R", fill, 0, "")
		o.pdf.Ln(-1)
		fill = !fill
	}
	o.pdf.SetX(left)
	o.pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
}

// 没有颜色的表格
func (o *OutputWay) pdfTableBodyFormat(pdf *gofpdf.Fpdf, w []float64, insertContent []string) *gofpdf.Fpdf {
	countryList := make([]countryType, 0, 8)
	// Column widths
	wSum := 0.0
	for _, v := range w {
		wSum += v
	}
	left := (210 - wSum) / 2
	// 	Header
	pdf.SetX(left)
	for j, str := range insertContent {
		pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	// Data
	for _, c := range countryList {
		pdf.SetX(left)
		pdf.CellFormat(w[0], 6, c.nameStr, "LR", 0, "", false, 0, "")
		pdf.CellFormat(w[1], 6, c.capitalStr, "LR", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 6, strDelimit(c.areaStr, ",", 3),
			"LR", 0, "R", false, 0, "")
		pdf.CellFormat(w[3], 6, strDelimit(c.popStr, ",", 3),
			"LR", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.SetX(left)
	pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")
	return pdf
}

func (o *OutputWay) tableInsert(w []float64, insertContent [][]string) {
	var pdf1 = o.pdf
	if insertContent != nil {
		for i := range insertContent {
			pdf1 = o.pdfTableBodyFormat(pdf1, w, insertContent[i])
		}
	}
}
