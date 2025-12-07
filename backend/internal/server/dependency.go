package server

import (
	"context"

	"github.com/jmoiron/sqlx"

	userHandler "devhub-backend/internal/api/http/handler/user"
	"devhub-backend/internal/api/http/middleware"
	httproute "devhub-backend/internal/api/http/route"
	dbUserRepo "devhub-backend/internal/infra/db/repository/user"
	"devhub-backend/internal/infra/logger"
	userUsecase "devhub-backend/internal/usecase/user"
)

//nolint:unparam
func (s *Server) setupRouteDependencies(ctx context.Context, appLogger logger.Logger, dbConn *sqlx.DB) (httproute.Dependency, error) {
	// Transactor factory
	// transactorFactory := infraDB.NewSqlxTransactorFactory(dbConn)

	// DB Repositories
	dbUserRepo := dbUserRepo.NewUserRepository(dbConn)
	// concertRepo := concertRepo.NewConcertRepository(dbConn)
	// zoneRepo := zonerepo.NewZoneRepository(dbConn)
	// seatRepo := seatRepo.NewSeatRepository(dbConn)
	// reservationRepo := reservationRepo.NewReservationRepository(dbConn)

	// Query retrier
	// queryBackoff, _ := retry.NewExponentialBackoffStrategy(500*time.Millisecond, 2.0, 5*time.Second)
	// queryRetrier, _ := retry.NewRetrier(retry.Config{
	// MaxAttempts: 3,
	// Backoff:     queryBackoff,
	// })

	// Usecases
	userUsecase := userUsecase.NewUserUsecase(s.cfg.App, dbUserRepo)
	// healthcheckUsecase := healthcheckUsecase.NewHealthCheckUsecase(queryRetrier, dbHealthRepo, redisHealthRepo)
	// concertUsecase := concertUsecase.NewConcertUsecase(s.cfg.App, transactorFactory, concertRepo)
	// seatUsecase := seatUsecase.NewSeatUsecase(s.cfg.App, concertRepo, zoneRepo, seatRepo, reservationRepo, transactorFactory, seatLockerRepo, seatMapRepo)

	// Application middleware
	appMiddleware := middleware.New()

	// Handlers
	userHandler := userHandler.NewUserHandler(s.cfg.App, userUsecase)
	// healthHandler := healthcheckHandler.NewHealthCheckHandler(healthcheckUsecase)
	// concertHandler := concertHandler.NewConcertHandler(s.cfg.App, concertUsecase)
	// seatHandler := seatHandler.NewSeatHandler(s.cfg.App, seatUsecase)

	return httproute.Dependency{
		Middleware:  appMiddleware,
		UserHandler: userHandler,
		// HealthCheckHandler: healthHandler,
		// ConcertHandler:     concertHandler,
		// SeatHandler:        seatHandler,
	}, nil
}
