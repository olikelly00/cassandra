package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main.go/models"
	"main.go/services"
)

// Local storage for interpretations and UUIDs

var LocalStorage map[string]string = make(map[string]string)

// GetRandomCard selects a random tarot card from the provided deck, ensuring that
// the selected card is not a duplicate of any card already present in the currentCards slice.

// deck is a slice of Card structs representing the full tarot deck.
// currentCards is a slice of Card structs representing the cards already selected.

// Returns a single Card struct from the deck that is not a duplicate of any card in currentCards.

func GetRandomCard(deck []models.Card, currentCards []models.Card) models.Card {
	randomiser := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		randomIndex := randomiser.Intn(len(deck))
		randomCard := deck[randomIndex]

		isDuplicate := false
		for _, card := range currentCards {
			if card.CardName == randomCard.CardName {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			return randomCard
		}
	}
}

// GetandInterpretThreeCards calls the FetchTarotCards function, which returns a slice of Card structs and assigns the variable name 'deck' to it.

// It then calls the drawThreeCards function, which takes the deck returned by FetchTarotCards, selects three random cards from the deck, returns them as a slice of Card structs under the variable name 'threeCards'.

// It then calls the convertCardsToJSON function, which takes the threeCards slice as input. It converts threeCards into three JSONCard structs and returns them. It also returns a slice of strings containing the names of the three cards. These strings will be interpolated into the prompt sent to OpenAI's API.

func GetandInterpretThreeCards(ctx *gin.Context) {
	deck, err := services.FetchTarotCards()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch tarot cards: %v", err),
		})
		return
	}
	requestID := uuid.New()

	threeCards := drawThreeCards(deck)

	jsonCards, cardNames := convertCardsToJSON(threeCards)

	ctx.JSON(http.StatusOK, gin.H{"cards": jsonCards, "requestID": requestID})
	userStory := ctx.Query("userstory")
	userName := ctx.Query("name")
	fmt.Print(userName, userStory)

	// This asynchronous function within GetandInterpretThreeCards sends the card names to the OpenAI API to generate an interpretation of the three cards. That interpretation is then stored in the LocalStorage map, with the requestID (UUID) as the key.

	// If the app is running in test mode, a test interpretation is generated and stored in the 'LocalStorage' map instead, to save wasting OpenAI credits.

	// Then the function calls the GetInterpretation function, responsile for finding and retrieving the interpretation from localStorage, using the requestID generated earlier.

	go func() {
		testing := os.Getenv("TESTING")
		if testing == "True" {
			interpretation := "This is a test interpretation"
			LocalStorage[requestID.String()] = interpretation
			return
		}
		apiKey := os.Getenv("API_KEY")
		interpretation, err := services.InterpretTarotCards(apiKey, cardNames, requestID, userStory, userName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to interpret tarot cards: %v", err),
			})
			return
		}
		LocalStorage[requestID.String()] = interpretation
		GetInterpretation(ctx)
		fmt.Println(interpretation)
	}()
}

func drawThreeCards(deck []models.Card) []models.Card {
	var threeCards []models.Card
	for i := 0; i < 3; i++ {
		threeCards = append(threeCards, GetRandomCard(deck, threeCards))
	}
	return threeCards
}

func convertCardsToJSON(cards []models.Card) ([]models.JSONCard, []string) {
	var jsonCards []models.JSONCard
	var cardNames []string

	for _, card := range cards {

		reversed := ReverseRandomiser()

		var FinalCardName string
		if reversed {
			FinalCardName = card.CardName + " (Reversed)"
		} else {
			FinalCardName = card.CardName
		}

		jsonCards = append(jsonCards, models.JSONCard{
			CardName:       FinalCardName,
			Type:           card.Type,
			MeaningUp:      card.MeaningUp,
			MeaningReverse: card.MeaningReverse,
			Description:    card.Description,
			ImageName:      card.ShortName + ".jpg",
			Reversed:       reversed,
		})

		if reversed {
			cardNames = append(cardNames, card.CardName, "(Reversed)")
		} else {
			cardNames = append(cardNames, card.CardName)
		}
	}

	return jsonCards, cardNames
}

func GetInterpretation(ctx *gin.Context) {
	requestID := ctx.Param("uuid")

	interpretation, ok := LocalStorage[requestID]

	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No interpretation found for this UUID"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"interpretation": interpretation})
}

func ReverseRandomiser() bool {
	randomiser := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomBool := randomiser.Intn(2)
	return randomBool == 0
}
