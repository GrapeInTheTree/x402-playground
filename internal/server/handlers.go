package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/GrapeInTheTree/x402-playground/pkg/health"
)

func WeatherHandler(c *gin.Context) {
	cities := []string{"Seoul", "Tokyo", "New York", "London", "Berlin", "Sydney"}
	conditions := []string{"Sunny", "Cloudy", "Rainy", "Snowy", "Windy", "Foggy"}

	city := cities[rand.Intn(len(cities))]
	c.JSON(http.StatusOK, gin.H{
		"city":        city,
		"temperature": 15 + rand.Intn(20),
		"condition":   conditions[rand.Intn(len(conditions))],
		"humidity":    30 + rand.Intn(60),
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	})
}

func JokeHandler(c *gin.Context) {
	jokes := []struct {
		Setup     string `json:"setup"`
		Punchline string `json:"punchline"`
	}{
		{"Why do programmers prefer dark mode?", "Because light attracts bugs."},
		{"What's a blockchain developer's favorite food?", "Hash browns."},
		{"Why did the smart contract fail?", "It had too many conditions."},
		{"How do you comfort a JavaScript bug?", "You console it."},
		{"Why do Go developers never get lost?", "They always know the goroutine."},
	}

	joke := jokes[rand.Intn(len(jokes))]
	c.JSON(http.StatusOK, joke)
}

func PremiumDataHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"report": "Premium Analytics Report",
		"metrics": gin.H{
			"totalTransactions": 1_234_567 + rand.Intn(100_000),
			"activeUsers":       89_432 + rand.Intn(10_000),
			"revenue":           fmt.Sprintf("$%d.%02d", 500_000+rand.Intn(100_000), rand.Intn(100)),
			"growthRate":        fmt.Sprintf("%.1f%%", 5.0+float64(rand.Intn(100))/10.0),
		},
		"generatedAt": time.Now().UTC().Format(time.RFC3339),
	})
}

func HealthHandler(service, network string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, health.Response{
			Status:  "ok",
			Service: service,
			Network: network,
		})
	}
}
