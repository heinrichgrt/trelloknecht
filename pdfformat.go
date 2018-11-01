package main

import (
	"strconv"

	"github.com/adlio/trello"
	"github.com/boombuler/barcode/qr"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/barcode"
	log "github.com/sirupsen/logrus"
)

func getPdfDocDimensionFromString() []float64 {
	r := make([]float64, 0)

	v, _ := strconv.ParseFloat(configuration["pdfDocXLen"], 64)
	r = append(r, v)
	v, _ = strconv.ParseFloat(configuration["pdfDocYLen"], 64)
	r = append(r, v)

	return r
}

func getPdfMarginsFromString() []float64 {
	r := make([]float64, 0)
	v, _ := strconv.ParseFloat(configuration["pdfMargin"], 64)
	r = append(r, v)
	v, _ = strconv.ParseFloat(configuration["pdfMargin"], 64)
	r = append(r, v)
	v, _ = strconv.ParseFloat(configuration["pdfMargin"], 64)
	r = append(r, v)
	return r

}

func getqRCodePosFromString() []float64 {
	r := make([]float64, 0)
	v, _ := strconv.ParseFloat(configuration["qRCodePosX"], 64)
	r = append(r, v)
	v, _ = strconv.ParseFloat(configuration["qRCodePosY"], 64)
	r = append(r, v)
	return r

}

func getBlackRectPosFromString() []float64 {
	r := make([]float64, 0)
	v, _ := strconv.ParseFloat(configuration["rectX0"], 64)
	r = append(r, v)
	v, _ = strconv.ParseFloat(configuration["rectY0"], 64)
	r = append(r, v)
	v, _ = strconv.ParseFloat(configuration["rectX1"], 64)
	r = append(r, v)
	v, _ = strconv.ParseFloat(configuration["rectY1"], 64)

	r = append(r, v)
	return r
}

func pdfBaseSetup() *gofpdf.Fpdf {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: configuration["pdfUnitStr"],
		Size:    gofpdf.SizeType{Wd: pdfDocDimension[0], Ht: pdfDocDimension[1]},
	})
	pdf.SetMargins(pdfMargins[0], pdfMargins[1], pdfMargins[2])
	pdf.AddPage()
	return pdf
}
func registerQR(pdf *gofpdf.Fpdf, card *trello.Card) {

	key := barcode.RegisterQR(pdf, card.Url, qr.H, qr.Unicode)
	qrSize, _ := strconv.ParseFloat(configuration["qRCodeSize"], 64)
	barcode.BarcodeUnscalable(pdf, key, qRCodePos[0], qRCodePos[1], &qrSize, &qrSize, false)

	// Output:
	// Successfully generated ../../pdf/contrib_barcode_RegisterQR.pdf
}

func writeLabel(pdf *gofpdf.Fpdf, card *trello.Card) string {
	headFontSize, _ := strconv.ParseFloat(configuration["headFontSize"], 64)
	pdf.SetFont(configuration["fontFamily"], configuration["headFontStyle"], headFontSize)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	_, lineHt := pdf.GetFontSize()

	// The option printQrCode is true by default, but can be set to false
	// via the config file.
	if configuration["printQrCode"] == "true" {
		registerQR(pdf, card)
	}
	headTopMargin, _ := strconv.ParseFloat(configuration["headTopMargin"], 64)
	pdf.SetTopMargin(headTopMargin)
	pdf.Rect(blackRectPos[0], blackRectPos[1], blackRectPos[2], blackRectPos[3], "D")
	headerString := card.Name
	htmlString := "<center><bold>" + shortenStringIfToLong(headerString) + "</bold></center>"
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, tr(htmlString))
	htmlString = "<left>" + boardNameByID[card.IDBoard] + " | " + listNameByID[card.IDList] + "</left>"
	pdf.SetFont("Times", "I", 8)
	pdf.SetAutoPageBreak(false, 0.0)
	_, lineHt = pdf.GetFontSize()
	lowerpos := lineHt + 6
	pdf.SetY(-lowerpos)
	html = pdf.HTMLBasicNew()
	html.Write(lineHt, tr(htmlString))
	lowerx := pdf.GetX()
	htmlString = "<right>" + joinedLabel(card) + "</right>"
	pdf.SetX(lowerx + 1)
	pdf.SetY(-lowerpos)
	html = pdf.HTMLBasicNew()
	html.Write(lineHt, tr(htmlString))
	fileName := configuration["tmpDirName"] + "/" + getUUID() + ".pdf"
	cardByFileName[fileName] = card

	err := pdf.OutputFileAndClose(fileName)

	if err != nil {
		log.Errorf("cannot create pdf-file %v\n", err)

	}
	return fileName
}

func printLabels(pdfList []string) {
	for _, pdf := range pdfList {
		commandResult := new(Resultset)
		commandResult.OSCommand = "/usr/bin/lp"
		commandResult.CommandArgs = []string{"-o", "fit-to-page", "-o", "media=" + configuration["printerMedia"], "-o",
			configuration["printerOrientation"], "-n", configuration["numberOfCopiesPrnt"], "-d", configuration["printerName"], pdf}
		commandResult.execCommand()
		if commandResult.SuccessfullExecution == true {
			printedCards = append(printedCards, pdf)
		}

	}

}
