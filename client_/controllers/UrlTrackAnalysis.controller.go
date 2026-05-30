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
