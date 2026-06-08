package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"shared/helpers"
	dto "shared/models"
	"shared/utils"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

type PanelController struct {
	BaseController
}

func (c *PanelController) Panel() {
	general, _ := c.Data["conf"].(dto.GeneralSettings)
	c.Data["menus"] = utils.GetPanelMenus(c.Lang, general)
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	req := dto.UsageAnalysisRequest{}

	// 1. Formdan gelen string tarih verilerini manuel olarak alıyoruz
	// HTML input type="date" için format genelde "2006-01-02" olur.
	// Eğer saat de varsa (datetime-local) format "2006-01-02T15:04" olabilir.
	const formDateFormat = "2006-01-02 15:04:05"

	startStr := c.GetString("start") + " 00:00:00"
	endStr := c.GetString("end") + " 23:59:59"

	// 2. String değerleri time.Time nesnesine çevirip req struct'ına atıyoruz
	if startStr != "" {
		parsedStart, err := time.Parse(formDateFormat, startStr)
		if err != nil {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]string{"error": "Başlangıç tarihi geçersiz formatta"}
			c.ServeJSON()
			return
		}
		req.Start = parsedStart
	}

	if endStr != "" {
		parsedEnd, err := time.Parse(formDateFormat, endStr)
		if err != nil {
			c.Ctx.Output.SetStatus(400)
			c.Data["json"] = map[string]string{"error": "Bitiş tarihi geçersiz formatta"}
			c.ServeJSON()
			return
		}
		req.End = parsedEnd
	}

	if err := c.ParseForm(&req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Geçersiz tarih formatı"}
		c.ServeJSON()
		return
	}

	usage, _ := utils.SendRequest[[]dto.UsageAnalysisResponse](req, "panel/GetURLStats", "POST", c.Ctx, tokenn)
	jsons, _ := json.Marshal(usage)
	c.Data["usage"] = string(jsons)
	if !req.Start.IsZero() {
		c.Data["start"] = req.Start.Format("2006-01-02")
	}
	if !req.End.IsZero() {
		c.Data["end"] = req.End.Format("2006-01-02")
	}

	c.Data["headerActive"] = "active"
	c.Data["active"] = "home"
	c.Layout = "inc/layout.html"
	c.TplName = "panel/index.html"
}

func (c *PanelController) Url() {
	general, _ := c.Data["conf"].(dto.GeneralSettings)
	c.Data["menus"] = utils.GetPanelMenus(c.Lang, general)
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	c.Data["headerActive"] = "active"
	c.Data["active"] = "urls"
	urls, _ := utils.SendRequest[[]dto.UrlResponse](nil, "panel/history", "GET", c.Ctx, tokenn)
	c.Data["urls"] = urls
	beego.ReadFromRequest(&c.Controller)
	c.Layout = "inc/layout.html"
	c.TplName = "panel/url/index.html"
}

func (c *PanelController) UrlAdd() {
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	id := c.Ctx.Input.Param(":id")
	c.Data["action"] = "panel/urls/save"
	c.Data["urls"] = nil
	if id != "" {
		urls, _ := utils.SendRequest[dto.UrlResponse](nil, "panel/historybyid/"+id, "GET", c.Ctx, tokenn)
		c.Data["urls"] = urls
		c.Data["action"] = "panel/urls/save/" + id
	}

	c.TplName = "panel/url/add.html"
}

func (c *PanelController) UrlSaveHandler() {
	id := c.Ctx.Input.Param(":id")
	idd, err := strconv.Atoi(id)
	if idd > 0 && err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.Body([]byte("Invalid ID format"))
		return
	}
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	var req any
	var action, method string
	if id == "" {
		req = &dto.CreateUrlRequest{}
		action = "panel/create"
		method = "POST"
	} else {
		req = &dto.UpdateUrlRequest{
			ID: idd,
		}
		action = "panel/update"
		method = "PUT"
	}

	if err := c.ParseForm(req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}
	status, err := utils.SendRequest[dto.GeneralResponse[any]](req, action, method, c.Ctx, tokenn)
	flash := beego.NewFlash()
	if status.Status {
		flash.Data["Success"] = "Başarılı"
	} else {
		flash.Data["Err"] = status.Message
	}
	flash.Store(&c.Controller)
	c.Redirect(helpers.Baseurl(c.Data["lang"], "panel/urls"), 302)
}

func (c *PanelController) UrlDelete() {
	id := c.Ctx.Input.Param(":id")
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	status, _ := utils.SendRequest[dto.GeneralResponse[any]](nil, "panel/delete/"+id, "DELETE", c.Ctx, tokenn)
	c.Data["json"] = status
	c.ServeJSON()

}

func (c *PanelController) BatchExcel() {
	file, _, err := c.GetFile("excel_file")
	flash := beego.NewFlash()
	if err != nil {
		flash.Data["err"] = "Dosya alınamadı"
		flash.Store(&c.Controller)
		c.Redirect(helpers.Baseurl(c.Data["lang"], "panel/urls"), 302)
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		flash.Data["err"] = "Excel okunamadı"
		flash.Store(&c.Controller)
		c.Redirect(helpers.Baseurl(c.Data["lang"], "panel/urls"), 302)
		return
	}
	defer f.Close()

	rows, err := f.GetRows("Sayfa1")
	if err != nil {
		flash.Data["err"] = "Sayfa bulunamadı"
		flash.Store(&c.Controller)
		c.Redirect(helpers.Baseurl(c.Data["lang"], "panel/urls"), 302)
		return
	}

	// 1. Struct tutacak slice'ımızı tanımlıyoruz
	var request []dto.ExcelUrlRequest
	layout := "2006-01-02 15:04:05"

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue // Başlık satırını atla
		}

		// Tamamen boş satırları atlamak için basit bir kontrol
		if len(row) == 0 {
			continue
		}

		// Hücreleri güvenli bir şekilde okumak için değişkenleri tanımlayalım
		var longUrl, alias, password, dateStr string

		// Kaç sütun dolu gelirse gelsin, indeks kontrolü yaparak değerleri atıyoruz:
		if len(row) > 0 {
			longUrl = strings.TrimSpace(row[0])
		}
		if len(row) > 1 {
			alias = strings.TrimSpace(row[1])
		}
		if len(row) > 2 {
			password = strings.TrimSpace(row[2])
		}
		if len(row) > 3 {
			dateStr = strings.Join(strings.Fields(row[3]), " ")
		}

		// Eğer zorunlu bir alanın varsa (örneğin LongUrl boş olamazsa) burada kontrol edebilirsin:
		if longUrl == "" {
			continue // Link yoksa bu satırı pas geç
		}

		var parsedDatePointer *time.Time

		// Tarih alanı doluysa ve güvenli indeksle alındıysa parse et
		if dateStr != "" {
			t, err := time.Parse(layout, dateStr)
			if err != nil {
				fmt.Println("Hatalı tarih formatı atlandı:", err.Error())
				continue
			}
			parsedDatePointer = &t
		}

		// Güvenli değişkenlerimizi struct modeline eşitle
		user := dto.ExcelUrlRequest{
			LongUrl:        longUrl,
			Alias:          alias,
			Password:       &password,
			ExpirationDate: parsedDatePointer,
		}

		request = append(request, user)
	}

	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	status, _ := utils.SendRequest[dto.GeneralResponse[any]](request, "panel/bulkcreate", "POST", c.Ctx, tokenn)
	if status.Status {
		flash.Data["Success"] = "Başarılı"
	} else {
		flash.Data["Err"] = status.Message
	}
	flash.Store(&c.Controller)
	c.Redirect(helpers.Baseurl(c.Data["lang"], "panel/urls"), 302)
}

// DownloadTemplate - static/files altındaki taslak dosyayı indirir
func (c *PanelController) DownloadTemplate() {
	// Dosyanın tam adını ve static klasöründeki yolunu belirliyoruz
	fileName := "url_taslak.xlsx"
	filePath := filepath.Join("static", "files", fileName)

	// Opsiyonel Güvenlik Kontrolü: Dosya gerçekten orada var mı?
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Dosya yoksa flash mesajı set edip yönlendirebilirsin
		flash := beego.NewFlash()
		flash.Data["Err"] = "İndirme başarısız: Taslak dosya sunucuda bulunamadı."
		flash.Store(&c.Controller)

		// panel/urls sayfasına geri postala
		c.Redirect(c.URLFor("PanelController.Url"), 302)
		return
	}

	// Beego'nun built-in download fonksiyonu.
	// İlk parametre: Dosyanın sunucudaki yolu
	// İkinci parametre (Opsiyonel): Kullanıcının bilgisayarına inerken görünecek dosya adı
	c.Ctx.Output.Download(filePath, "url_toplu_yukleme_taslagi.xlsx")
}
func (c *PanelController) Qr() {
	id := c.Ctx.Input.Param(":id")
	c.Data["id"] = id
	c.Data["href"] = helpers.Baseurl(c.Data["slug"], "panel/urls/qrcreate/"+id)
	c.Data["logohref"] = helpers.Baseurl(c.Data["slug"], "panel/urls/logoUpload")
	c.TplName = "panel/url/qr.html"

}
func (c *PanelController) QrCreate() {
	id := c.Ctx.Input.Param(":id")
	token := c.Ctx.Input.GetData("token")
	tokenn := token.(string)
	resultUrl, _ := utils.SendRequest[dto.UrlResponse](nil, "panel/historybyid/"+id, "GET", c.Ctx, tokenn)
	shortURL := helpers.GetShortUrl(resultUrl.Alias, resultUrl.Password, resultUrl.ShortCode, "qr")
	var png []byte
	png, err := qrcode.Encode(shortURL, qrcode.Highest, 400)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("QR kod üretilemedi")
		return
	}

	c.Ctx.Output.Header("Content-Type", "image/png")
	c.Ctx.Output.Header("Content-Length", fmt.Sprintf("%d", len(png)))
	c.Ctx.Output.Body(png)
}

func (c *PanelController) UploadLogo() {
	// 1. Formdan "logoUpload" isimli dosyayı yakala
	file, header, err := c.GetFile("logoUpload")
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Dosya alınamadı: " + err.Error()}
		c.ServeJSON()
		return
	}
	defer file.Close()

	// 2. Dosya uzantısını kontrol et (Sadece PNG ve JPG/JPEG)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Sadece PNG, JPG veya JPEG formatında logolar yüklenebilir."}
		c.ServeJSON()
		return
	}

	// 3. Benzersiz bir dosya adı oluştur (Çakışmaları önlemek için MD5 + Timestamp)
	hasher := md5.New()
	hasher.Write([]byte(time.Now().String() + header.Filename))
	fileName := fmt.Sprintf("%x%s", hasher.Sum(nil), ext)

	// 4. Kaydedilecek hedef klasörü belirle (Klasörün var olduğundan emin ol)
	uploadDir := "static/uploads/logos/"
	uploadPath := filepath.Join(uploadDir, fileName)

	// 5. Dosyayı sunucuya kaydet
	err = c.SaveToFile("logoUpload", uploadPath)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "Dosya kaydedilirken bir hata oluştu."}
		c.ServeJSON()
		return
	}
	webPath := "/" + uploadPath
	c.Data["json"] = map[string]interface{}{"success": true, "logo_url": webPath, "message": "Logo başarıyla yüklendi."}
	c.ServeJSON()
}
