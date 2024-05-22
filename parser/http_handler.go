package parser

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleContactRequest handles the POST request for contact information
func HandleContactRequest(service Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ContactRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Process the contact request using the service
		res, err := service.ParseIncomingRequest(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Send back the response
		c.JSON(http.StatusOK, res)
	}
}
