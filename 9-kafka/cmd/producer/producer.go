package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"kafka/pkg/models"
	"log"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

var ErrUserNotFoundInProducer = errors.New("user not found")

func main() {
	users := []models.User{
		{
			ID:   1,
			Name: "Tùng",
		},
		{
			ID:   2,
			Name: "Mino",
		},
		{
			ID:   3,
			Name: "LeoNI",
		},

		{
			ID:   4,
			Name: "Puppo",
		},
		{
			ID:   5,
			Name: "Jerrsy",
		},
	}

	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/send", sendMessageHandler(producer, users))

	fmt.Printf("Kafka PRODUCER 📨 started at http://localhost%s\n",
		ProducerPort)

	err = router.Run(ProducerPort)
	if err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}

// Kafka functions
func sendKafkaMessage(
	producer sarama.SyncProducer,
	users []models.User,
	ctx *gin.Context,
	fromID int,
	toID int,
) error {
	message := ctx.PostForm("message")

	fromUser, err := findUserByID(fromID, users)
	if err != nil {
		return err
	}

	toUser, err := findUserByID(toID, users)
	if err != nil {
		return err
	}

	notification := models.Notification{
		From:    fromUser,
		To:      toUser,
		Message: message,
	}

	notificationJson, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notifcation: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(strconv.Itoa(toUser.ID)),
		Value: sarama.StringEncoder(notificationJson),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func sendMessageHandler(
	producer sarama.SyncProducer,
	users []models.User,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fromID, err := getIDFromRequest("fromID", ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		toID, err := getIDFromRequest("toID", ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		err = sendKafkaMessage(producer, users, ctx, fromID, toID)
		if errors.Is(err, ErrUserNotFoundInProducer) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "user not found",
			})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "notification sent successfully!",
		})
	}
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(
		[]string{KafkaServerAddress},
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}

// Helper Functionos
func findUserByID(id int, users []models.User) (models.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFoundInProducer
}

func getIDFromRequest(formValue string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.PostForm(formValue))
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID from form value %s, %w", formValue, err)
	}

	return id, nil
}
