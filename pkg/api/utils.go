package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"dojo-api/db"
	"dojo-api/pkg/auth"
	"dojo-api/pkg/event"
	"dojo-api/pkg/metric"
	"dojo-api/pkg/miner"
	"dojo-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Define a common response structure
type ApiResponse struct {
	Success bool        `json:"success" swaggertype:"boolean"`
	Body    interface{} `json:"body"`
	Error   interface{} `json:"error"`
}

func defaultErrorResponse(errorMsg interface{}) ApiResponse {
	return ApiResponse{Success: false, Body: nil, Error: errorMsg}
}

func defaultSuccessResponse(body interface{}) ApiResponse {
	return ApiResponse{Success: true, Body: body, Error: nil}
}

func handleCurrentSession(c *gin.Context) (*auth.SecureCookieSession, error) {
	session, exists := c.Get("session")
	if !exists {
		return nil, errors.New("no session found")
	}

	currSession, ok := session.(auth.SecureCookieSession)
	if !ok {
		return nil, errors.New("invalid session")
	}
	return &currSession, nil
}

func buildApiKeyResponse(apiKeys []db.APIKeyModel) miner.MinerApiKeysResponse {
	keys := make([]string, 0)
	for _, apiKey := range apiKeys {
		keys = append(keys, apiKey.Key)
	}
	return miner.MinerApiKeysResponse{
		ApiKeys: keys,
	}
}

func buildSubscriptionKeyResponse(subScriptionKeys []db.SubscriptionKeyModel) miner.MinerSubscriptionKeysResponse {
	keys := make([]string, 0)
	for _, subScriptionKey := range subScriptionKeys {
		keys = append(keys, subScriptionKey.Key)
	}
	return miner.MinerSubscriptionKeysResponse{
		SubscriptionKeys: keys,
	}
}

func handleMetricData(currentTask *db.TaskModel, updatedTask *db.TaskModel) {
	metricService := metric.NewMetricService()
	eventService := event.NewEventService()
	ctx := context.Background()

	// Always update total task results count
	go func() {
		if err := metricService.UpdateTotalTaskResultsCount(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to update total tasks results count")
		} else {
			log.Info().Msg("Updated total task results count")
		}
	}()

	// Only update completed task count when task gets its first result
	// TODO: need to consider race condition
	if updatedTask.NumResults == 1 {
		go func() {
			if err := metricService.UpdateCompletedTaskCount(ctx); err != nil {
				log.Error().Err(err).Msg("Failed to update completed task count")
			} else {
				log.Info().Msg("Updated completed task count")
			}
		}()
	}

	// Handle task completion events and metrics
	// TODO: reconsider this logic for task completion events, and avg task completion time
	if (currentTask.Status != db.TaskStatusCompleted) && updatedTask.Status == db.TaskStatusCompleted {
		go func() {
			// Update the task completion event
			if err := eventService.CreateTaskCompletionEvent(ctx, *updatedTask); err != nil {
				log.Error().Err(err).Msg("Failed to create task completion event")
			} else {
				log.Info().Msg("Created task completion event")
			}
		}()

		go func() {
			// Update the avg task completion
			if err := metricService.UpdateAvgTaskCompletionTime(ctx); err != nil {
				log.Error().Err(err).Msg("Failed to update average task completion time")
			} else {
				log.Info().Msg("Updated average task completion time")
			}
		}()
	}
}

// Get the user's IP address from the gin request headers
func getCallerIP(c *gin.Context) string {
	if runtimeEnv := utils.LoadDotEnv("RUNTIME_ENV"); runtimeEnv == "aws" {
		forwardedFor := c.Request.Header.Get("X-Original-Forwarded-For")
		if forwardedFor != "" {
			// Split the string by comma and get the last IP
			ips := strings.Split(forwardedFor, ",")
			if len(ips) > 0 {
				// Trim any whitespace from the last IP
				lastIP := strings.TrimSpace(ips[len(ips)-1])
				log.Debug().Msgf("Got last caller IP from X-Original-Forwarded-For header: %s", lastIP)
				return lastIP
			}
		}
	}
	callerIp := c.ClientIP()
	log.Debug().Msgf("Got caller IP from ClientIP: %s", callerIp)
	return callerIp
}

func CustomGinLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			c.AbortWithStatus(500)
			return
		}

		requestSize := len(body)
		log.Printf("Request size: %d bytes", requestSize)

		// Restore the request body to the context
		c.Request.Body = io.NopCloser(io.NopCloser(bytes.NewBuffer(body)))

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = getCallerIP(c)
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		// param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		consoleWriter := zerolog.ConsoleWriter{
			Out: os.Stderr,
		}
		consoleWriter.FormatLevel = func(i interface{}) string {
			return "GIN"
		}

		logger := log.With().Logger().Output(consoleWriter)

		// Log main request info
		logger.Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status_code", param.StatusCode).
			Str("latency", param.Latency.String()).
			Int("req_size", len(body)).
			Int("resp_size", param.BodySize).
			Str("ip", param.ClientIP).
			Msg("")

		// Log headers separately
		logger.Info().
			Interface("headers", c.Request.Header).
			Msg("")
	}
}
