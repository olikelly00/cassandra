package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendInternalError handles internal server errors in a Gin web server context. It sends an appropriate error response to the client based on the application mode.

// ctx is the Gin context for the request. If ctx is nil, the error is printed to the console.
// err is the error that occurred.

// In release mode, a generic "Something went wrong" error message is sent to the client.
// In other modes, the specific error message is sent.

// Note: To enable release mode, call gin.SetMode(gin.ReleaseMode) in the main function.

func SendInternalError(ctx *gin.Context, err error) {
	if ctx == nil {
		fmt.Printf("Internal Server Error: %v\n", err)
		return
	}
	if gin.Mode() == "release" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": "Something went wrong"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}
}
