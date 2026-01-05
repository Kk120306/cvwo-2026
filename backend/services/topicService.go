package services

import (
	"errors"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/helpers"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"gorm.io/gorm"
)

// TopicService handles topic business logic
type TopicService struct{}

// NewTopicService creates a new instance of TopicService
func NewTopicService() *TopicService {
	return &TopicService{}
}

// CreateTopicInput represents the data needed to create a topic
type CreateTopicInput struct {
	Name string
}

// UpdateTopicInput represents the data needed to update a topic
type UpdateTopicInput struct {
	Name string
}

// GetAllTopics retrieves all topics sorted by creation date
func (s *TopicService) GetAllTopics() ([]models.Topic, error) {
	// Create a slice for the topics
	var topics []models.Topic

	// query all the topics from database ordered - appends it to topic slice
	result := database.DB.Order("created_at desc").Find(&topics)

	// If there was a error in retrieving from the database
	if result.Error != nil {
		return nil, errors.New("failed to retrieve topics")
	}

	return topics, nil
}

// CreateTopic creates a new topic with auto-generated slug
func (s *TopicService) CreateTopic(input CreateTopicInput) (*models.Topic, error) {
	// Validate input
	if input.Name == "" {
		return nil, errors.New("topic name cannot be empty")
	}

	topic := models.Topic{
		Name: input.Name,
		Slug: helpers.GenerateSlug(input.Name),
	}

	// Creating the Topic
	result := database.DB.Create(&topic)

	// if there was an error during creation
	if result.Error != nil {
		return nil, errors.New("failed to create topic")
	}

	return &topic, nil
}

// FindTopicBySlug finds a topic by its slug
func (s *TopicService) FindTopicBySlug(slug string) (*models.Topic, error) {
	// Find the topic
	// https://gorm.io/docs/update.html
	var topic models.Topic
	result := database.DB.First(&topic, "slug = ?", slug)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("topic not found")
		}
		return nil, errors.New("failed to retrieve topic")
	}

	return &topic, nil
}

// UpdateTopic updates a topic's name and regenerates slug
func (s *TopicService) UpdateTopic(topic *models.Topic, input UpdateTopicInput) error {
	// Validate input
	if input.Name == "" {
		return errors.New("topic name cannot be empty")
	}

	// Updating the topic name and slug
	newSlug := helpers.GenerateSlug(input.Name)
	// Since topic is already found, it calls WHERE = loaded ID
	res := database.DB.Model(topic).Updates(models.Topic{
		Name: input.Name,
		Slug: newSlug,
	})

	// If there was an error during update
	if res.Error != nil {
		return errors.New("failed to update topic")
	}

	// Update the topic object with new values
	topic.Name = input.Name
	topic.Slug = newSlug

	return nil
}

// DeleteTopic deletes a topic by slug
func (s *TopicService) DeleteTopic(slug string) error {
	// Delete through Gorm - passing in the slug
	// https://gorm.io/docs/delete.html
	var topic models.Topic
	result := database.DB.Delete(&topic, "slug = ?", slug)
	if result.Error != nil {
		return errors.New("failed to delete topic")
	}

	// When there were no changes to the database, means topic was not found
	if result.RowsAffected == 0 {
		return errors.New("topic not found")
	}

	return nil
}