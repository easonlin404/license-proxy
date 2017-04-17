package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"io/ioutil"
	"license-proxy/util"
	"net/http"
	"strconv"
)

var provider = "widevine_test"
var url = "https://license.uat.widevine.com/cenc/getlicense/" + provider
var allowedTrackTypes = "SD_HD"
var key = "1ae8ccd0e7985cc0b6203a55855a1034afc252980e970ca90e5202689f947ab9"
var iv = "d58ce954203b7c9a9a9d467f59839249"

type LicenseRequest struct {
	Request   string `json:"request"`
	Signature string `json:"signature"`
	Signer    string `json:"signer"`
}

type Request struct {
	Payload           string `json:"payload"`
	ContentId         string `json:"content_id"`
	Provider          string `json:"provider"`
	AllowedTrackTypes string `json:"allowed_track_types"`
}

func main() {

	r := gin.Default()

	r.POST("/proxy", func(c *gin.Context) {
		body, _ := ioutil.ReadAll(c.Request.Body)

		// allow cross domain AJAX requests
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		if len(body) == 0 {
			c.String(400, string("Empty Request"))
			return
		}

		fmt.Println("Requst POST body:")
		fmt.Println(string(body))

		resBody, s := generateLicense(body)

		status, _ := strconv.Atoi(s)

		var licenseResponse map[string]interface{}
		json.Unmarshal([]byte(resBody), &licenseResponse)


		indentJson, _ := json.Marshal(licenseResponse)
		fmt.Println("response Body:", string(indentJson))

		jsonStatus := licenseResponse["status"].(string)
		if jsonStatus == "OK" {
			license := licenseResponse["license"].(string)
			fmt.Println(license)

			licenseDecode, _ := base64.StdEncoding.DecodeString(license)
			c.String(status, string(licenseDecode))
		} else {
			c.String(status, jsonStatus)
		}

	})

	r.Run(":9000") // listen and serve on 0.0.0.0:9000

}

func buildMessage(body []byte) []byte {
	var request Request
	request.Payload = base64.StdEncoding.EncodeToString(body)

	//request.ContentId
	request.Provider = provider
	request.AllowedTrackTypes = allowedTrackTypes

	message, _ := json.Marshal(request)
	return message
}

func genrateLicenseRequest(body []byte) LicenseRequest {
	keyByteAry, _ := hex.DecodeString(key)
	ivByteAry, _ := hex.DecodeString(iv)

	message := buildMessage(body)

	var licenseRequest LicenseRequest
	licenseRequest.Request = base64.StdEncoding.EncodeToString(message)
	licenseRequest.Signature = util.GenerateSignature(keyByteAry, ivByteAry, message)
	licenseRequest.Signer = provider

	return licenseRequest
}

func generateLicense(body []byte) (string, string) {
	licenseRequest := genrateLicenseRequest(body)
	jsonStr, err := json.Marshal(licenseRequest)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Println(string(jsonStr))

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
	resBody, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(resBody))

	return string(resBody), resp.Status
}
