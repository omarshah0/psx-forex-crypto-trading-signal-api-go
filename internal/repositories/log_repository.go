package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository struct {
	collection *mongo.Collection
}

func NewLogRepository(db *mongo.Database) *LogRepository {
	return &LogRepository{
		collection: db.Collection("request_logs"),
	}
}

// RequestLog represents the structure of a request log
type RequestLog struct {
	Timestamp   time.Time              `bson:"timestamp"`
	Method      string                 `bson:"method"`
	Path        string                 `bson:"path"`
	Headers     map[string]interface{} `bson:"headers"`
	Body        interface{}            `bson:"body,omitempty"`
	QueryParams map[string]interface{} `bson:"query_params,omitempty"`
	Response    ResponseLog            `bson:"response"`
	UserID      *int64                 `bson:"user_id,omitempty"`
	IPAddress   string                 `bson:"ip_address"`
	Duration    int64                  `bson:"duration_ms"`
}

// ResponseLog represents the structure of a response log
type ResponseLog struct {
	Status  int                    `bson:"status"`
	Headers map[string]interface{} `bson:"headers"`
	Body    interface{}            `bson:"body,omitempty"`
}

// Create creates a new request log
func (r *LogRepository) Create(ctx context.Context, log *RequestLog) error {
	_, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create request log: %w", err)
	}
	return nil
}

// GetLogs retrieves logs with optional filters
func (r *LogRepository) GetLogs(ctx context.Context, filter bson.M, limit int64, skip int64) ([]RequestLog, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []RequestLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("failed to decode logs: %w", err)
	}

	return logs, nil
}
