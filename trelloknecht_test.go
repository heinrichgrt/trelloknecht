package main

import (
	"os"
	"testing"

	"github.com/adlio/trello"
	log "github.com/sirupsen/logrus"
)

var client *trello.Client

func TestMain(m *testing.M) {
	// we need this for trello board interaction for the testfunctions
	readConfigFromFile("config.cfg")
	readConfigFromFile(".token")
	log.Infof("token: %v\n", configuration["trelloToken"])
	client = trello.NewClient(configuration["trelloAppKey"], configuration["trelloToken"])
	board, err := client.GetBoard("5bceb330ba13f689ee477774", trello.Defaults())
	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatalf("Can not get cards from config board")
	}
	for _, card := range cards {
		log.Infof("card: %v id: %v", card.Name, card.ID)
	}

	lists, err := board.GetLists(trello.Defaults())
	for _, list := range lists {

		log.Infof("list %v id: %v\n", list.Name, list.ID)
	}
	log.Infof("started\n")
	os.Exit(m.Run())

}
