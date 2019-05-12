package storage

import (
	"context"
	"fmt"

	"github.com/rlongo/ambrosia/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbname         = "ambrosia"
	collectionName = "recipes"
)

// AmbrosiaStorageMongo is a MongoDB driver implementation
type AmbrosiaStorageMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// Open creates a new DB connection
func OpenMongo(storageConnectionString string) (*AmbrosiaStorageMongo, error) {
	clientOptions := options.Client().ApplyURI(storageConnectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbname).Collection(collectionName)
	return &AmbrosiaStorageMongo{client: client, collection: collection}, err
}

// Close cleans up the DB driver
func (db *AmbrosiaStorageMongo) Close() error {
	return db.client.Disconnect(context.Background())
}

func (db *AmbrosiaStorageMongo) GetRecipes(filterTags []string, filterAuthor string) (api.Recipes, error) {
	var (
		result api.Recipes
		cur    *mongo.Cursor
		filter bson.M
		err    error
	)

	hasFilterTags := filterTags != nil && len(filterTags) > 0
	hasFilterAuthor := len(filterAuthor) > 0

	if hasFilterTags {
		filter = bson.M{"tags": bson.M{"$in": filterTags}}
	}

	if hasFilterAuthor {
		filterAuthor := bson.M{"author": filterAuthor}

		if hasFilterTags {
			filter = bson.M{"$and": bson.A{filter, filterAuthor}}
		} else {
			filter = filterAuthor
		}
	}

	if cur, err = db.collection.Find(context.Background(), filter); err == nil {
		if err = cur.All(context.Background(), &result); err == nil {
			fmt.Println("Got all recipes!")
			return result, nil
		}
	}

	return nil, err
}

func (db *AmbrosiaStorageMongo) GetRecipe(id api.RecipeID) (api.Recipe, error) {
	return api.Recipe{}, fmt.Errorf("Result not found")
}

func (db *AmbrosiaStorageMongo) AddRecipe(recipe *api.Recipe) error {
	_, err := db.collection.InsertOne(context.Background(), recipe)
	return err
}

func (db *AmbrosiaStorageMongo) UpdateRecipe(recipe *api.Recipe) error {
	return fmt.Errorf("Result not found")
}
