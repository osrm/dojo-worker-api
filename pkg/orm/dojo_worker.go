package orm

import (
	"context"
	"errors"

	"dojo-api/db"

	"github.com/rs/zerolog/log"
)

type DojoWorkerORM struct {
	dbClient      *db.PrismaClient
	clientWrapper *PrismaClientWrapper
}

func NewDojoWorkerORM() *DojoWorkerORM {
	clientWrapper := GetPrismaClient()
	return &DojoWorkerORM{dbClient: clientWrapper.Client, clientWrapper: clientWrapper}
}

func (s *DojoWorkerORM) CreateDojoWorker(walletAddress string, chainId string) (*db.DojoWorkerModel, error) {
	s.clientWrapper.BeforeQuery()
	defer s.clientWrapper.AfterQuery()

	ctx := context.Background()
	worker, err := s.dbClient.DojoWorker.CreateOne(
		db.DojoWorker.WalletAddress.Set(walletAddress),
		db.DojoWorker.ChainID.Set(chainId),
	).Exec(ctx)
	return worker, err
}

func (s *DojoWorkerORM) GetById(ctx context.Context, workerId string) (*db.DojoWorkerModel, error) {
	s.clientWrapper.BeforeQuery()
	defer s.clientWrapper.AfterQuery()

	worker, err := s.dbClient.DojoWorker.FindUnique(
		db.DojoWorker.ID.Equals(workerId),
	).Exec(ctx)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			log.Error().Err(err).Msg("Worker not found")
			return nil, err
		}
		return nil, err
	}
	return worker, nil
}

func (s *DojoWorkerORM) GetDojoWorkerByWalletAddress(walletAddress string) (*db.DojoWorkerModel, error) {
	s.clientWrapper.BeforeQuery()
	defer s.clientWrapper.AfterQuery()

	ctx := context.Background()
	worker, err := s.dbClient.DojoWorker.FindFirst(
		db.DojoWorker.WalletAddress.Equals(walletAddress),
	).Exec(ctx)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			log.Error().Err(err).Msg("Worker not found")
			return nil, err
		}
		return nil, err
	}
	return worker, nil
}
