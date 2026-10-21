package tools

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func MakeQrcode(text string) string {
	qrBytes, _ := qrcode.Encode(text, qrcode.Low, 256)

	qrString := base64.StdEncoding.EncodeToString(qrBytes)
	return qrString
}
