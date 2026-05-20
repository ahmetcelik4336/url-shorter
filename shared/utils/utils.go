package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

var apiurl []byte

func InitUtils() error {
	key := os.Getenv("APIURL")
	apiurl = []byte(key)
	return nil
}

func CallAPI(method, url, token string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	apikey := os.Getenv("APIKEY")
	req.Header.Set("APIKEY", apikey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}
func GenerateDates(date string) (time.Time, time.Time) {
	parsedDate, _ := time.Parse("2006-01-02", date)
	start := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 1).Add(-time.Nanosecond)
	return start, end
}
func GetTokenFromHeader(header string) string {
	return strings.Replace(header, "Bearer ", "", 1)
}

func SendRequest[T any](req any, path, method string, ctx *context.Context, token string) (T, error) {
	var res T

	rawToken := ctx.Input.GetData("token")
	if rawToken == nil {
		rawToken = token
	}
	tokenStr, _ := rawToken.(string)

	apiUrl := string(apiurl)

	var bodyReader io.Reader
	if req != nil && method != "GET" {
		jsonData, err := json.Marshal(req)
		if err != nil {
			return res, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	resp, err := CallAPI(method, apiUrl+path, tokenStr, bodyReader)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	//fmt.Println("Gelen Ham Veri:", string(bodyBytes))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Println("hata", string(bodyBytes))
		return res, fmt.Errorf("API hatası! Status")
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		//return res, fmt.Errorf("decode error: %v", err)
	}
	//log.Println(string(bodyBytes))
	return res, nil
}
