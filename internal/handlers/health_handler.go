package handlers

import (
	"net/http"

	"github.com/omarshah0/rest-api-with-social-auth/internal/database"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type HealthHandler struct {
	postgresDB *database.PostgresDB
	mongoDB    *database.MongoDB
	redisDB    *database.RedisDB
}

func NewHealthHandler(postgresDB *database.PostgresDB, mongoDB *database.MongoDB, redisDB *database.RedisDB) *HealthHandler {
	return &HealthHandler{
		postgresDB: postgresDB,
		mongoDB:    mongoDB,
		redisDB:    redisDB,
	}
}

// HealthCheck returns the health status of all services
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{
		"postgres": "healthy",
		"mongodb":  "healthy",
		"redis":    "healthy",
	}

	overallHealthy := true

	// Check PostgreSQL
	if err := h.postgresDB.HealthCheck(); err != nil {
		status["postgres"] = "unhealthy: " + err.Error()
		overallHealthy = false
	}

	// Check MongoDB
	if err := h.mongoDB.HealthCheck(); err != nil {
		status["mongodb"] = "unhealthy: " + err.Error()
		overallHealthy = false
	}

	// Check Redis
	if err := h.redisDB.HealthCheck(); err != nil {
		status["redis"] = "unhealthy: " + err.Error()
		overallHealthy = false
	}

	if overallHealthy {
		utils.SendSuccess(w, http.StatusOK, utils.ResponseTypeAction, status, "All services are healthy")
	} else {
		utils.SendError(w, http.StatusServiceUnavailable, utils.ErrorTypeServiceUnavailable, "One or more services are unhealthy")
	}
}

