package storage

import (
	"context"
	"fmt"

	guid "github.com/google/uuid"
	"github.com/rlongo/ambrosia/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbname            = "ambrosia"
	RecipesCollection = "recipes"
)

// AmbrosiaStorageMongo is a MongoDB driver implementation
type AmbrosiaStorageMongo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// Open creates a new DB connection
func OpenMongo(storageConnectionString string, collectionName string) (*AmbrosiaStorageMongo, error) {
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
		filter = bson.M{"tags": bson.M{"$in": []string{filterTags[0]}}}

		for i := 1; i < len(filterTags); i++ {
			tagFilter := bson.M{"tags": bson.M{"$in": []string{filterTags[i]}}}
			filter = bson.M{"$and": bson.A{filter, tagFilter}}
		}
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
			return result, nil
		}
	}

	return nil, err
}

func (db *AmbrosiaStorageMongo) GetRecipe(id api.RecipeID) (api.Recipe, error) {
	var result api.Recipe

	if err := db.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result); err != nil {
		return api.Recipe{}, fmt.Errorf("No matches found")
	}

	return result, nil
}

func (db *AmbrosiaStorageMongo) AddRecipe(recipe *api.Recipe) error {
	uuid, err := guid.NewUUID()
	if err != nil {
		return err
	}
	recipe.ID = api.RecipeID(uuid)
	_, err = db.collection.InsertOne(context.Background(), *recipe)
	return err
}

func (db *AmbrosiaStorageMongo) UpdateRecipe(recipe *api.Recipe) error {
	result, err := db.collection.UpdateOne(context.Background(), bson.M{"_id": recipe.ID}, bson.M{"$set": recipe})
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("No Documents Updated")
	}

	return nil
}
