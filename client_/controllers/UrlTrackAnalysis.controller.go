package controllers

import (

	// Takma ad (alias) vermeden direkt paket yollarını yazıyoruz

	"encoding/json"
	dto "shared/models"
	"shared/utils"
	"time"
)

type UrlTrackAnalysis struct {
	BaseController
}

func (c *UrlTrackAnalysis) Get() {

	result, _ := utils.SendRequest[dto.UrlTrackAnalysisResponseBatch](nil, "panel/TotalReading", "POST", c.Ctx, "")
	c.Data["data"] = result.Analysis
	c.Data["LastReading"] = "Henüz okuma yok"
	if result.LastReading != nil {
		c.Data["LastReading"] = result.LastReading.LastReading
	}

	req := dto.UsageAnalysisRequest{}
	if err := c.ParseForm(&req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Geçersiz tarih formatı"}
		c.ServeJSON()
		return
	}
	result1, _ := utils.SendRequest[[]*dto.UsageAnalysisResponse](req, "panel/urlTrackAnalysis", "POST", c.Ctx, "")
	jsons, _ := json.Marshal(result1)
	c.Data["graphic"] = string(jsons)
	general, _ := c.Data["conf"].(dto.GeneralSettings)
	c.Data["menus"] = utils.GetPanelMenus(c.Lang, general)
	if !req.Start.IsZero() {
		c.Data["start"] = req.Start.Format("2006-01-02")
	} else {
		c.Data["start"] = time.Now().Add(-7 * 24 * time.Hour).Format("2006-01-02")
	}
	if !req.End.IsZero() {
		c.Data["end"] = req.End.Format("2006-01-02")
	} else {
		c.Data["end"] = time.Now().Format("2006-01-02")
	}
	c.Data["headerActive"] = "active"
	c.Data["active"] = "analysis"
	c.Layout = "inc/layout.html"
	c.TplName = "panel/analysis/urltrack.html"

}
func (c *UrlTrackAnalysis) LogDatatable() {
	draw, _ := c.GetInt("draw", 1)
	start, _ := c.GetInt("start", 0)
	length, _ := c.GetInt("length", 10)
	searchVal := c.GetString("search[value]")

	orderColumnIdx := c.GetString("order[0][column]")
	orderDir := c.GetString("order[0][dir]") // "asc" veya "desc"
	if orderDir == "" {
		orderDir = "desc"
	}

	orderColumnName := c.GetString("columns[" + orderColumnIdx + "][data]")

	startDateStr := c.GetString("startdate")
	endDateStr := c.GetString("enddate")

	// Go'nun parse edebilmesi için beklediğimiz format (Örn: 2026-05-31)
	// Eğer frontend'den "31.05.2026" geliyorsa burayı "02.01.2006" yapmalısın!
	const dateFormat = "2006-01-02"

	var startDate, endDate time.Time

	if startDateStr != "" {
		startDate, _ = time.Parse(dateFormat, startDateStr)
	}

	if endDateStr != "" {
		endDate, _ = time.Parse(dateFormat, endDateStr)
	}

	dtRequest := dto.DataTableRequest{
		Draw:           draw,
		Start:          start,
		Length:         length,
		SearchValue:    searchVal,
		OrderColumnIdx: orderColumnName,
		OrderDir:       orderDir,
		StartDate:      startDate,
		EndDate:        endDate,
	}

	result, _ := utils.SendRequest[*dto.DataTableResponse](dtRequest, "panel/LogDatatable", "POST", c.Ctx, "")
	c.Data["json"] = result
	c.ServeJSON()

}
