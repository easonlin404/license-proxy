package main

import (
	"bytes"
	"catchplay.com/license-proxy/proxy"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {

	r := gin.Default()

	r.POST("/proxy", func(c *gin.Context) {
		body, status := generateLicence()
		fmt.Println(body)

		s, _ := strconv.Atoi(status)

		c.JSON(s, unmarshal(body))

		//unmarshal()
	})

	r.Run() // listen and serve on 0.0.0.0:8080

}

func unmarshal(jsonText string) map[string]interface{} {
	var f map[string]interface{}
	//str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(jsonText), &f)
	fmt.Println(f)
	return f
}



type LicenseRequest struct {
	Request   string `json:"request"`
	Signature string `json:"signature"`
	Signer    string `json:"signer"`
}

func genrateLicenseRequest() LicenseRequest {
	//TODO
	key, _ := hex.DecodeString("1ae8ccd0e7985cc0b6203a55855a1034afc252980e970ca90e5202689f947ab9")
	iv, _ := hex.DecodeString("d58ce954203b7c9a9a9d467f59839249")

	plaintextBase64 := "ew0KICAicGF5bG9hZCIgOiAiQ0FFU2hBRUtUQWdBRWtnQUFBQUNBQUFRV1BYYmh0Yi9xNDNmM1NmdUMyVlAzcTBqZUFFQ1czZW1Ra1duMndYQ1lWT252bFdQRE5xaDhWVklCNEdtc05BOGVWVkZpZ1hrUVdJR04wR2xnTUtqcFVFU0xBb3FDaFFJQVJJUUpNUEN6bDJiVml5TVFFdHlLL2d0bVJBQkdoQXlOV1kzT0RNek1UY3lNbUpqTTJFeUdBRWd2NWlRa0FVYUlDM09OMXpWZ2VWMHJQN3cyVm1WTEdvcnFDbGNNUU80QmRiSFB5azNHc25ZIiwNCiAgInByb3ZpZGVyIiA6ICJ3aWRldmluZV90ZXN0IiwNCiAgImNvbnRlbnRfaWQiOiAiYWEiLA0KICAiY29udGVudF9rZXlfc3BlY3MiOiBbDQogICAgeyAidHJhY2tfdHlwZSI6ICJIRCIgfSwNCiAgICB7ICJ0cmFja190eXBlIjogIkFVRElPIiB9DQogIF0NCn0="
	message, _ := base64.StdEncoding.DecodeString(plaintextBase64)

	signature := proxy.GenerateSignature(key, iv, message)

	var licenseRequest LicenseRequest
	licenseRequest.Request = plaintextBase64
	licenseRequest.Signature = signature
	licenseRequest.Signer = "widevine_test"

	return licenseRequest
}

func generateLicence() (string, string) {
	licenseRequest:= genrateLicenseRequest()
	jsonStr, err := json.Marshal(licenseRequest)
	if err != nil {
		fmt.Println("json err:", err)
	}


	fmt.Println(string(jsonStr))
	//var jsonStr = []byte(`{
	//"request": "ewogICJwYXlsb2FkIiA6ICJDQUVTaEFFS1RBZ0FFa2dBQUFBQ0FBQVFXUFhiaHRiL3E0M2YzU2Z1QzJWUDNxMGplQUVDVzNlbVFrV24yd1hDWVZPbnZsV1BETnFoOFZWSUI0R21zTkE4ZVZWRmlnWGtRV0lHTjBHbGdNS2pwVUVTTEFvcUNoUUlBUklRSk1QQ3psMmJWaXlNUUV0eUsvZ3RtUkFCR2hBeU5XWTNPRE16TVRjeU1tSmpNMkV5R0FFZ3Y1aVFrQVVhSUMzT04xelZnZVYwclA3dzJWbVZMR29ycUNsY01RTzRCZGJIUHlrM0dzblkiLAogICJwcm92aWRlciIgOiAid2lkZXZpbmVfdGVzdCIsCiAgImNvbnRlbnRfaWQiOiAiWm10cU0yeHFZVk5rWm1Gc2EzSXphZz09IiwKICAiY29udGVudF9rZXlfc3BlY3MiOiBbCiAgICB7ICJ0cmFja190eXBlIjogIlNEIiB9LAogICAgeyAidHJhY2tfdHlwZSI6ICJIRCIgfSwKICAgIHsgInRyYWNrX3R5cGUiOiAiQVVESU8iIH0KICBdCn0K",
	//"signature":"xPkAbb3tjOY/ybdz0tmJMq9erH9ILnS5natMZr3QEW8=",
	//"signer": "widevine_test"
//}`)

	url := "https://license.uat.widevine.com/cenc/getlicense/widevine_test"
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return string(body), resp.Status
}
