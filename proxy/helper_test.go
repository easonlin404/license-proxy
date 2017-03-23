package proxy

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGenerateSignature(*testing.T) {
	key, _ := hex.DecodeString("1ae8ccd0e7985cc0b6203a55855a1034afc252980e970ca90e5202689f947ab9")
	iv, _ := hex.DecodeString("d58ce954203b7c9a9a9d467f59839249")

	plaintextBase64 := "ew0KICAicGF5bG9hZCIgOiAiQ0FFU2hBRUtUQWdBRWtnQUFBQUNBQUFRV1BYYmh0Yi9xNDNmM1NmdUMyVlAzcTBqZUFFQ1czZW1Ra1duMndYQ1lWT252bFdQRE5xaDhWVklCNEdtc05BOGVWVkZpZ1hrUVdJR04wR2xnTUtqcFVFU0xBb3FDaFFJQVJJUUpNUEN6bDJiVml5TVFFdHlLL2d0bVJBQkdoQXlOV1kzT0RNek1UY3lNbUpqTTJFeUdBRWd2NWlRa0FVYUlDM09OMXpWZ2VWMHJQN3cyVm1WTEdvcnFDbGNNUU80QmRiSFB5azNHc25ZIiwNCiAgInByb3ZpZGVyIiA6ICJ3aWRldmluZV90ZXN0IiwNCiAgImNvbnRlbnRfaWQiOiAiYWEiLA0KICAiY29udGVudF9rZXlfc3BlY3MiOiBbDQogICAgeyAidHJhY2tfdHlwZSI6ICJIRCIgfSwNCiAgICB7ICJ0cmFja190eXBlIjogIkFVRElPIiB9DQogIF0NCn0="
	message, _ := base64.StdEncoding.DecodeString(plaintextBase64)
	fmt.Println("Message:")
	fmt.Println(string(message))

	signature := GenerateSignature(key, iv, message)

	fmt.Println("signature:" + signature)

}
