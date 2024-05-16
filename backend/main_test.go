package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"main.go/controllers"
	"main.go/env"
	"main.go/models"
)

type TestSuiteEnv struct {
	suite.Suite
	app *gin.Engine
	res *httptest.ResponseRecorder
}

func RequestSetup(app *gin.Engine, suite *TestSuiteEnv, reqType string, path string) []byte {
	req, _ := http.NewRequest(reqType, path, nil)
	app.ServeHTTP(suite.res, req)
	responseData, _ := io.ReadAll(suite.res.Body)
	return responseData
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite() {
	env.LoadEnv(".test.env")

	suite.app = setupApp()

}

func (suite *TestSuiteEnv) SetupTest() {
	suite.res = httptest.NewRecorder()

}

// This function actually runs our test suite when we run 'go test' in terminal.
func TestSuite(t *testing.T) {
	suite.Run(t, new(TestSuiteEnv))
}

// This test sets up a GET request to the /cards path and unmarshals the response data from JSON into Go structs. The test then uses the assert package to check the response code of this request was 200/OK.
func (suite *TestSuiteEnv) Test_GetThreeCards_ResponseCode() {
	app := suite.app
	responseData := RequestSetup(app, suite, "GET", "/cards")

	var jsonCards struct {
		Cards []models.JSONCard
	}

	_ = json.Unmarshal(responseData, &jsonCards)

	assert.Equal(suite.T(), 200, suite.res.Code)
}

// This test sets up a GET request to the /cards path and unmarshals the response data from JSON into Go structs. The test then uses the assert package to check the length of the slice of Card structs is 3.
func (suite *TestSuiteEnv) Test_GetThreeCards_ExpectedFormat() {
	app := suite.app
	responseData := RequestSetup(app, suite, "GET", "/cards")

	var jsonCards struct {
		Cards []models.JSONCard
	}

	_ = json.Unmarshal(responseData, &jsonCards)

	assert.Len(suite.T(), jsonCards.Cards, 3)
}

// This function sets up a GET request to fetch three cards from the tarot deck. It then sets up a second request and compares the results of both to check they are different. This tests that the 'random' generation is working properly.
func (suite *TestSuiteEnv) Test_GetThreeCardsIsRandom() {
	app := suite.app

	//Response 1
	responseData := RequestSetup(app, suite, "GET", "/cards")
	var jsonCards struct {
		Cards []models.JSONCard
	}

	_ = json.Unmarshal(responseData, &jsonCards)
	//Response 2
	responseData2 := RequestSetup(app, suite, "GET", "/cards")
	var jsonCards2 struct {
		Cards []models.JSONCard
	}

	_ = json.Unmarshal(responseData2, &jsonCards2)
	assert.NotEqual(suite.T(), jsonCards.Cards[0].CardName, jsonCards2.Cards[0].CardName) //1 in 24336 of this failing!!!
}

// This test sets up a GET request to the /cards path (to fetch three cards) and then to the /cards/interpret/requestID path (to get a reading of those cards).
// It then unmarshals the response data from JSON into Go structs. The test then uses the assert package to check the response code of this request was 200/OK.

func (suite *TestSuiteEnv) Test_GetAndInterpretCards_ResponseCode() {
	app := suite.app

	responseData := RequestSetup(app, suite, "GET", "/cards")
	var jsonCards struct {
		Cards     []models.JSONCard
		RequestID uuid.UUID
	}

	_ = json.Unmarshal(responseData, &jsonCards)

	_ = RequestSetup(app, suite, "GET", "/cards/interpret/"+jsonCards.RequestID.String())

	assert.Equal(suite.T(), 200, suite.res.Code)
}

// This test sets up a GET request to the /cards path (to fetch three cards) and then to the /cards/interpret/requestID path (to get a reading of those cards).
// It then unmarshals the response data from JSON into Go structs. The test then uses the assert package to check that, while in test mode, the interpretation of the cards is "This is a test interpretation".

func (suite *TestSuiteEnv) Test_GetAndInterpretCards_ExpectedFormat() {
	app := suite.app

	responseData := RequestSetup(app, suite, "GET", "/cards")
	var jsonCards struct {
		Cards     []models.JSONCard
		RequestID uuid.UUID
	}

	_ = json.Unmarshal(responseData, &jsonCards)

	responseData2 := RequestSetup(app, suite, "GET", "/cards/interpret/"+jsonCards.RequestID.String())
	var interpretationResponse struct {
		Interpretation string
	}

	_ = json.Unmarshal(responseData2, &interpretationResponse)

	assert.Equal(suite.T(), controllers.LocalStorage[jsonCards.RequestID.String()], "This is a test interpretation")
}
