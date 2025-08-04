package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	var (
		flagMongoDB string
		flagDBName  string
	)

	var (
		ctx         = context.Background()
		compressors = []string{"snappy", "zlib", "zstd"}
		collections = []string{
			"b7s-batches",
			"b7s-batch-chunks",
			"b7s-batch-work-items",
		}
	)

	pflag.StringVarP(&flagMongoDB, "server", "s", "", "mongodb server to use")
	pflag.StringVarP(&flagDBName, "db-name", "d", "b7s-db", "db to use")

	pflag.Parse()

	slog.Info("mongodb server info",
		"server", flagMongoDB,
		"db", flagDBName)

	opts := options.Client().
		SetCompressors(compressors).
		ApplyURI(flagMongoDB)

	client, err := mongo.Connect(opts)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			slog.Error("error disconnecting from the server", "error", err)
		}
	}()

	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = client.Ping(pctx, readpref.Primary())
	if err != nil {
		return err
	}

	db := client.Database(flagDBName)

	for _, collection := range collections {

		c := db.Collection(collection)

		r, err := c.DeleteMany(ctx, bson.M{})
		if err != nil {
			slog.Error("could not delete items", "collection", collection, "error", err)
			continue
		}

		slog.Info("items deleted",
			"collection", collection,
			"deleted", r.DeletedCount,
		)
	}

	slog.Info("all done")

	return nil
}
