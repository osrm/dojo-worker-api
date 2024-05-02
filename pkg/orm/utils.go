package orm

import (
	"dojo-api/db"
	"sync"

	"github.com/rs/zerolog/log"
)

type SimpleConnHandler struct {
	client      *db.PrismaClient
	isConnected bool
}

var connHandler *SimpleConnHandler
var once sync.Once

func GetConnHandler() *SimpleConnHandler {
	once.Do(func() {
		connHandler = &SimpleConnHandler{
			client:      db.NewClient(),
			isConnected: false,
		}
	})
	return connHandler
}

func GetPrismaClient() *db.PrismaClient {
	handler := GetConnHandler()
	if !handler.isConnected {
		if err := handler.client.Prisma.Connect(); err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to Prisma client")
			return nil
		}
		handler.isConnected = true
	}
	return handler.client
}

func (h *SimpleConnHandler) OnShutdown() error {
	if h.client == nil {
		log.Warn().Msg("Prisma client not initialised")
		return nil
	}

	if err := h.client.Prisma.Disconnect(); err != nil {
		log.Error().Err(err).Msg("Failed to disconnect from Prisma client")
		return err
	} else {
		h.isConnected = false
	}
	return nil
}
