package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func crc16(data string) string {
	crc := 0xFFFF
	for i := 0; i < len(data); i++ {
		crc ^= int(data[i]) << 8
		for j := 0; j < 8; j++ {
			if (crc & 0x8000) > 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc = crc << 1
			}
		}
	}
	crc &= 0xFFFF
	return fmt.Sprintf("%04X", crc)
}

func GenerateDynamicQRIS(staticQR string, amount float64) string {
	dynamicQR := strings.Replace(staticQR, "010211", "010212", 1)

	amountStr := strconv.FormatFloat(amount, 'f', 0, 64)
	amountLen := fmt.Sprintf("%02d", len(amountStr))
	amountTag := "54" + amountLen + amountStr

	indexCRC := strings.LastIndex(dynamicQR, "6304")
	if indexCRC != -1 {
		dynamicQR = dynamicQR[:indexCRC]
	}

	dynamicQR = dynamicQR + amountTag + "6304"

	newCRC := crc16(dynamicQR)

	return dynamicQR + newCRC
}