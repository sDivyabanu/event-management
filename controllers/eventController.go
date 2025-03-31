package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"backend-challenge/config"
	"backend-challenge/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Date        string `json:"date" binding:"required"`
}

type EventResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	OrganizerID string `json:"organizer_id"`
}

func CreateEvent(c *gin.Context) {
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	organizerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(organizerID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organizer ID"})
		return
	}

	event := models.Event{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		Date:        req.Date,
		OrganizerID: objID,
	}

	collection := config.GetCollection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Event creation failed"})
		return
	}

	response := EventResponse{
		ID:          event.ID.Hex(),
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		OrganizerID: event.OrganizerID.Hex(),
	}

	c.JSON(http.StatusCreated, response)
}

func GetAllEvents(c *gin.Context) {
	collection := config.GetCollection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}
	defer cursor.Close(ctx)

	var wg sync.WaitGroup
	eventChan := make(chan models.Event)
	var events []models.Event

	
	go func() {
		for event := range eventChan {
			events = append(events, event)
		}
	}()

	for cursor.Next(ctx) {
		var event models.Event
		if err := cursor.Decode(&event); err != nil {
			continue
		}
		
		wg.Add(1)
		go func(e models.Event) {
			defer wg.Done()
			eventChan <- e
		}(event)
	}

	wg.Wait()
	close(eventChan)

	var response []EventResponse
	for _, event := range events {
		response = append(response, EventResponse{
			ID:          event.ID.Hex(),
			Name:        event.Name,
			Description: event.Description,
			Date:        event.Date,
			OrganizerID: event.OrganizerID.Hex(),
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetEventByID(c *gin.Context) {
	eventID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	collection := config.GetCollection("events")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var event models.Event
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&event)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	response := EventResponse{
		ID:          event.ID.Hex(),
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		OrganizerID: event.OrganizerID.Hex(),
	}

	c.JSON(http.StatusOK, response)
}