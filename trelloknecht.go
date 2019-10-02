package main

import (
	"time"

	// "github.com/adlio/trello"
	"github.com/adlio/trello"
)

var (
	configuration = map[string]string{
		// pdf setings
		"fontFamily": "Helvetica",
		"pdfUnitStr": "mm",
		"pdfDocXLen": "100.0",
		"pdfDocYLen": "62.0",
		"pdfMargin":  "3.0",

		"headLineCharsSkip": "82",

		"printQrCode":   "true",
		"qRCodeSize":    "30.0",
		"qRCodePosX":    "66.0",
		"qRCodePosY":    "25.0",
		"headFontStyle": "B",
		"headFontSize":  "16.0",
		"headTopMargin": "5.0",
		"rectX0":        "3.0",
		"rectY0":        "2.0",
		"rectX1":        "95.0",
		"rectY1":        "58.0",

		// trello settings
		"trelloAppKey":          "",
		"trelloToken":           "",
		"toPrintedLabelName":    "PRINTME_DEVOPS",
		"newLabelAfterPrint":    "PRINTED",
		"knechtID":              "",
		"trelloUserName":        "kls_drucker",
		"configTrelloBoardID":   "someIdentifier",
		"boardsToWatch":         "DevOps 2020 - Board",
		"configCardName":        "PrintBert02 Card",
		"usePrinterStatusBoard": "true",
		"printFrame":            "false",
		"ConfigListOnBoard":     "IDs",
		"printerMedia":          "Custom.62x100mm",
		"printerOrientation":    "landscape",
		"printerName":           "Brother_QL_700",
		"tmpDirName":            "",
		"tmpDirPrefix":          "trelloKnecht",
		"numberOfCopiesPrnt":    "2",
		"waitIntervalSeconds":   "60",
	}
	//utility vars

	newLabelAfterPrtIDs = make(map[string]string)
	boardNameByID       = make(map[string]string)
	listNameByID        = make(map[string]string)
	listIDByName        = make(map[string]string)
	labelIDByName       = make(map[string]string)
	cardByFileName      = make(map[string]*trello.Card)
	printedCards        = make([]string, 0)
	// composed vars

	pdfDocDimension = []float64{}
	pdfMargins      = []float64{}
	qRCodePos       = []float64{}
	blackRectPos    = []float64{}
	boardsToWatch   = []string{}
	configFile      = ""
	tokenFile       = ""

	// printer settings

)

//Resultset  Json for output
type Resultset struct {
	OSCommand            string    `json:"os.cmd"`
	CommandArgs          []string  `json:"cmd.args"`
	Stdout               string    `json:"stdout"`
	Stderr               string    `json:"stderr,omitempty"`
	CmdStarttime         time.Time `json:"cmd.starttime"`
	CMDStoptime          time.Time `json:"cmd.stoptime"`
	DurationSecounds     int       `json:"duration.seconds"`
	SuccessfullExecution bool      `json:"succesful"`
	ErrorStr             string    `json:"errorstr,omitempty"`
}

func main() {
	defer cleanUp(configuration["tmpDirName"])
	//sleeptime, err := strconv.ad(configuration["waitIntervalSeconds"])

	createIPCardOnBoard()

	for {
		cardList := getLabels()
		pdfFileList := writeLabels(cardList)
		printLabels(pdfFileList)
		swapLabel(cardList)
		reportPrints()
		sweepOut()
		time.Sleep(60 * time.Second)
	}

}
