package main

import "log"

func RunStrategy() {
	CCI := int(Desc01m.CCI[499])
	RSI := int(Desc01m.RSI[499])
	MFI := int(Desc01m.MFI[499])

	if (CCI < -150) && (RSI < 30) && (MFI < 20) && (P != "L") {
		log.Println("LONG", CCI, RSI, MFI, L, S)
		P = "L"
		L++
	} else if (CCI > 150) && (RSI > 70) && (MFI > 80) && (P != "S") {
		log.Println("SHORT", CCI, RSI, MFI, L, S)
		P = "S"
		S++
	}
}
