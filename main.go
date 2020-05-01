package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adlio/trello"

	log "github.com/sirupsen/logrus"
)

func getPreviousTime(t time.Duration) time.Time {
	now := time.Now()
	return now.Add(-t)
}

func main() {
	appKey := os.Getenv("TRELLO_APP_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	boardId := os.Getenv("TRELLO_BOARD_ID")
	username := os.Getenv("TRELLO_USERNAME")

	client := trello.NewClient(appKey, token)

	trelloBoard, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {
		log.WithField("err", err).Fatalf("client.GetBoard(%v) has failed", boardId)
	}

	lists, err := trelloBoard.GetLists(trello.Defaults())
	if err != nil {
		log.WithField("err", err).Fatal("trelloBoard.GetLists() has failed.")
	}

	interestedCards := make(map[string][]*trello.Card)

	for _, list := range lists {
		cards, err := list.GetCards(trello.Defaults())
		if err != nil {
			log.WithField("err", err).Warn("list.GetCards() failed")
			continue
		}
		for _, card := range cards {
			if card.DateLastActivity.After(getPreviousTime(time.Hour * 24)) {
				member, err := card.CreatorMember()
				if err != nil {
					log.WithError(err).Warn("card.CreatorMember() failed")
					continue
				}

				if member.Username != username {
					log.WithFields(log.Fields{
						"member.Username": member.Username,
						"username":        username,
					}).Info("Ignoring the card since the card is not created by the user")
					continue
				}
				interestedCards[list.Name] = append(interestedCards[list.Name], card)
			}
		}
	}

	for listName, cards := range interestedCards {
		fmt.Println(listName)
		fmt.Println(strings.Repeat("=", len(listName)))
		printCards(cards)
		fmt.Println()
	}
}

func printCards(cards []*trello.Card) {
	for _, card := range cards {
		fmt.Printf("- %v\n", card.Name)
	}
}