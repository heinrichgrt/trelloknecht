package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/adlio/trello"
	log "github.com/sirupsen/logrus"
)

func isPrintedLabelOnBoard(card *trello.Card) bool {
	res := true
	for _, label := range card.Labels {
		if label.Name == configuration["newLabelAfterPrint"] {
			res = false
		}

	}
	return res
}

func swapLabel(cards []*trello.Card) {
	for _, card := range cards {
		r := new(trello.Label)
		//var l card.Labels
		err := card.RemoveIDLabel(labelIDByName[configuration["toPrintedLabelName"]], r)
		if err != nil {
			log.Fatalf("removing  Label : %v with %v \n", configuration["toPrintedLabelName"], err)
		}
		if isPrintedLabelOnBoard(card) {
			err = card.AddIDLabel(newLabelAfterPrtIDs[card.IDBoard])
			if err != nil {
				log.Fatalf("adding Label: %v  with %v\n", configuration["newLabelAfterPrint"], err)
			}
		}
	}
}

func writeLabels(cardList []*trello.Card) []string {
	pdfFileList := make([]string, 0)
	for _, card := range cardList {
		pdf := pdfBaseSetup()
		pdfFileName := writeLabel(pdf, card)
		pdfFileList = append(pdfFileList, pdfFileName)
	}

	return pdfFileList
}

func getBoards(client *trello.Client) []*trello.Board {
	member, err := client.GetMember(configuration["trelloUserName"], trello.Defaults())
	if err != nil {
		log.Fatal("Cannot  get member info from trello")
	}

	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		log.Fatal("Cannot get board lists from trello")
	}
	return boards
}

func joinedLabel(card *trello.Card) string {
	labelString := ""
	labelList := make([]string, 0)
	for _, label := range card.Labels {
		if matched, _ := regexp.MatchString("PRINT.*", label.Name); matched == false {
			labelList = append(labelList, label.Name)
		}
	}
	labelString = strings.Join(labelList, ", ")
	return labelString
}

func filterBoards(userBoards []*trello.Board) []*trello.Board {
	boardList := make([]*trello.Board, 0)
	for boardID := range userBoards {
		for watchID := range boardsToWatch {
			if userBoards[boardID].Name == boardsToWatch[watchID] {
				boardList = append(boardList, userBoards[boardID])
			}
		}

	}
	return boardList
}
func boarListIDsToNames(board *trello.Board) {
	lists, _ := board.GetLists(trello.Defaults())
	for _, list := range lists {
		// todo add this to cleanup
		listNameByID[list.ID] = list.Name
		listIDByName[list.Name] = list.ID

	}

}
func getOwnCardFromPrinterBoard(c *trello.Client) *trello.Card {
	board, err := c.GetBoard(configuration["configTrelloBoardID"], trello.Defaults())
	if err != nil {
		log.Fatalf("Can not get config board data")
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatalf("Can not get cards from config board")
	}
	for _, card := range cards {
		if card.Name == configuration["configCardName"] {
			return card
		}
	}
	return nil
}

func createOwnCard(client *trello.Client) {
	list, err := client.GetList(listIDByName["IPs"], trello.Defaults())
	if err != nil {
		log.Fatal("cannot get the ip list from the board\n")
	}

	dueDate := time.Now().AddDate(0, 0, 1)

	card := trello.Card{
		Name:     configuration["configCardName"],
		Desc:     configCardDescription(),
		Due:      &dueDate,
		IDLabels: []string{"label1", "label2"},
	}

	err = list.AddCard(&card, trello.Arguments{"pos": "bottom"})
	if err != nil {
		log.Fatalf("cannot add card %v to the board with %v", card, err)
	}

}
func createIPCardOnBoard() {
	if !strings.EqualFold(configuration["usePrinterStatusBoard"], "true") {
		return
	}
	// does the board exist?
	client := trello.NewClient(configuration["trelloAppKey"], configuration["trelloToken"])

	//board, err := client.GetBoard("5bceb330ba13f689ee477774", trello.Defaults())
	board, err := client.GetBoard(configuration["configTrelloBoardID"], trello.Defaults())
	if err != nil {
		log.Fatalf("The configuration board: %v can not be reached. Check if it exist and this user can access it\n", configuration["configTrelloBoardID"])
		return
	}
	boarListIDsToNames(board)
	onwCard := getOwnCardFromPrinterBoard(client)
	if onwCard != nil {
		updateOwnCard(onwCard)
	} else {
		createOwnCard(client)
	}
	log.Debugf("ListID: %v", listIDByName["IPs"])
	log.Infof("Trello printer: %v stored info on config board", configuration["knechtID"])
}

func cardAddAndRemove() {

}
func createTestCard(client *trello.Client) {
	list, err := client.GetList(listIDByName["IPs"], trello.Defaults())
	if err != nil {
		log.Fatal("cannot get the ip list from the board\n")
	}

	dueDate := time.Now().AddDate(0, 0, 1)

	card := trello.Card{
		Name:     "TestCard",
		Desc:     "TestDescription",
		Due:      &dueDate,
		IDLabels: []string{"label1", "label2"},
	}

	err = list.AddCard(&card, trello.Arguments{"pos": "bottom"})
	if err != nil {
		log.Fatalf("cannot add card %v to the board with %v", card, err)
	}

}
func getLabels() []*trello.Card {
	cardList := make([]*trello.Card, 0)
	client := trello.NewClient(configuration["trelloAppKey"], configuration["trelloToken"])
	boards := getBoards(client)

	for _, board := range filterBoards(boards) {
		boarListIDsToNames(board)
		getPrintedLabelID(board)
		boardNameByID[board.ID] = board.Name
		cardList = append(cardList, getMatchingCardsFromBoard(board)...)
	}

	return cardList

}

func getMatchingCardsFromBoard(board *trello.Board) []*trello.Card {
	cardList := make([]*trello.Card, 0)

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal("cannot get cards from board")
	}
	for _, card := range cards {

		for _, label := range card.Labels {
			log.Debugf("label %v on %v", label, card)
			if label.Name == configuration["toPrintedLabelName"] {
				cardList = append(cardList, card)
				labelIDByName[label.Name] = label.ID
			}
		}
	}
	return cardList
}

func updateOwnCard(card *trello.Card) {
	card.Desc = configCardDescription()
	args := trello.Arguments{
		"desc": card.Desc,
	}
	err := card.Update(args)
	if err != nil {
		log.Fatalf("updating card on printer board failed")

	}

}

func getPrintedLabelID(board *trello.Board) {
	labels, err := board.GetLabels(trello.Defaults())
	if err != nil {
		log.Fatalf("cannot get labels from board: %v\n", err)
	}
	for _, label := range labels {
		if label.Name == configuration["newLabelAfterPrint"] {
			newLabelAfterPrtIDs[board.ID] = label.ID
		}
	}
}
func shortenStringIfToLong(instring string) string {
	wordList := strings.Split(instring, " ")
	shortendString := ""
	iterator := 0
	headLineLength, err := strconv.Atoi(configuration["headLineCharsSkip"])
	if err != nil {
		log.Fatal("configvalue headLineCharsSkip is nan")
	}

	if err != nil {
		log.Fatal("configvalue headLineMaxChars is nan")
	}
	for len(shortendString) < headLineLength && iterator < len(wordList) {

		shortendString += " " + wordList[iterator]
		iterator++
	}
	if iterator < len(wordList) {
		shortendString += "..."
	}

	return strings.Trim(shortendString, " ")
}
