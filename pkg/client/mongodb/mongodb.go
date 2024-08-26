package mongodb

import (
	"context"
	"fmt"
	"github.com/Henus321/boney-james-go-backend/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, sc config.StorageConfig) (db *mongo.Database, err error) {
	var mongoDBURL string
	var isAuth bool
	if sc.Username == "" && sc.Password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", sc.Host, sc.Port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", sc.Username, sc.Password, sc.Host, sc.Port)
	}

	// это тоже может быть в connection string выше
	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if isAuth {
		clientOptions.SetAuth(options.Credential{
			AuthSource: sc.Database,
			Username:   sc.Username,
			Password:   sc.Password,
		})
	}

	// Connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB due to error: %v", err)
	}

	// Ping
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB due to error: %v", err)
	}

	return client.Database(sc.Database), nil
}
