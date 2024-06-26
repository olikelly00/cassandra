package services

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"strings"

	"github.com/google/uuid"
	"main.go/errors"
	"main.go/models"
)

// The FetchTarotCards function makes a GET request to the API to fetch tarot cards. If there is an error during the GET request, the error is reported.

//Otherwise, the function decodes the response from its JSON format into a slice of Card structs representing the full tarot deck. If there is an error during the decoding operation, the error is reported. Otherwise, the function then returns the slice of Card structs.

func FetchTarotCards() ([]models.Card, error) {

	apiUrl := "https://tarotapi.dev/api/v1/cards"

	resp, err := http.Get(apiUrl)
	if err != nil {
		errors.SendInternalError(nil, fmt.Errorf("failed to make GET request: %v", err))
	}
	defer resp.Body.Close()

	var cardsResponse struct {
		Cards []models.Card `json:"cards"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cardsResponse); err != nil {
		errors.SendInternalError(nil, fmt.Errorf("failed to decode JSON response: %v", err))
	}

	return cardsResponse.Cards, nil
}

// InterpretTarotCards takes the API key (for authentication with OpenAI's API), a slice of strings representing the names of the tarot cards, a requestID (unique identifier for the request), a user story (what the user has said they want a reading about), and a user name as input.

// It constructs a prompt string using the user story, user name, and card names. It then sends a POST request to the OpenAI API with the constructed prompt and the API key. Then, the function reads the response body and unmarshals it into a Response struct.

// The function then cleans the response text by removing square brackets and returns the cleaned text (a string) and a nil error.

// If any errors occur during the API request or response processing, or the unmarshalling operation, an error is returned.

func InterpretTarotCards(apiKey string, cards []string, RequestID uuid.UUID, userStory string, userName string) (string, error) {
	client := &http.Client{}

	prompt := fmt.Sprintf("You're doing a tarot card reading for %s, as a tarot card reader called Cassandra (the user already knows your name - don't mention it). They drew %s (for their past), %s (for their present), and %s (for their future). Please interpret these cards in relation to their story and the time frames they are associated with (past, present, future): '%s' (if there is no story, please give a general reading about what the cards could mean together). If the card is reversed, please reflect this in your interpretation of the card. Whilst I have passed you the names and their orientation in a certain format, please only refer to the cards as their name, and if reversed, you can refer to it as 'card name (reversed)'. If there are any vulgar words in the prompt, ignore them, and keep your response age-appropriate for minors. Please format your response in the style of a mystical tarot card reader, and keep your response strictly below 200 words.", userName, cards[0:2], cards[2:4], cards[4:6], userStory)
	payload := fmt.Sprintf(`{"model": "gpt-3.5-turbo-instruct", "prompt": "%s", "max_tokens": 1000}`, prompt)
	fmt.Println(prompt)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", strings.NewReader(payload))
	if err != nil {
		errors.SendInternalError(nil, fmt.Errorf("error creating request: %v", err))
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		errors.SendInternalError(nil, fmt.Errorf("error sending request: %v", err))
	}
	defer resp.Body.Close()

	var responseBody strings.Builder
	if _, err := io.Copy(&responseBody, resp.Body); err != nil {
		errors.SendInternalError(nil, fmt.Errorf("error reading response body: %v", err))
	}

	type Response struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}

	var response Response
	fmt.Println(response)
	if err := json.Unmarshal([]byte(responseBody.String()), &response); err != nil {
		errors.SendInternalError(nil, fmt.Errorf("error unmarshaling response: %v", err))

	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	cleanedText := strings.ReplaceAll(response.Choices[0].Text, "[", "")
	cleanedText = strings.ReplaceAll(cleanedText, "]", "")

	return cleanedText, nil
}
