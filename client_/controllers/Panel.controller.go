package controllers

type PanelController struct {
	BaseController
}

func (c *PanelController) Panel() {
	/*req := &dto.UrlCountAnalysisRequest{Date: ""}
	_, err := utils.SendRequest[dto.UrlCountAnalysisResponse](req, "panel/GetURLCount", "POST", c.Ctx)

	if err != nil {
		fmt.Println("Hata oluştu:", err)
		return
	}*/

	c.TplName = "auth/register.html"
}
