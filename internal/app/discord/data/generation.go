package data

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"sync"
)

// Global cache for storing color images
var (
	cache = make(map[string]*ColorGeneration)
	mutex sync.RWMutex
)

type ColorGeneration struct {
	Input     string
	ColorName string
	ColorInt  int
	ImageData []byte
	Embed     *discordgo.MessageEmbed
	TempMsgID string
}

func FindGeneration(id string) *ColorGeneration {
	mutex.RLock()

	gen, exists := cache[id]
	if exists {
		delete(cache, id)
	}

	mutex.RUnlock()

	return gen
}

// SaveGeneration stores an image in the cache with the given ID
func SaveGeneration(gen *ColorGeneration) string {
	// Generate a UUID for the image
	id := uuid.New().String()

	// Store the image in the cache
	mutex.Lock()
	cache[id] = gen
	mutex.Unlock()

	return id
}
