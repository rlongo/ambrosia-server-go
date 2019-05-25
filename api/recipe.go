package api

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"

	"go.mongodb.org/mongo-driver/bson/bsontype"

	guid "github.com/google/uuid"
)

// RecipeID is the id field for a recipe. This is used by the DB for indexing
type RecipeID [16]byte

// The Recipe struct is a wrapper for the entire recipe
type Recipe struct {
	ID     RecipeID `json:"_id" bson:"_id"`
	Name   string   `json:"name" bson:"name"`
	Author string   `json:"author" bson:"author"`
	Rating uint8    `json:"rating" bson:"rating"`
	Notes  string   `json:"notes" bson:"notes"`
	Tags   []string `json:"tags" bson:"tags"`
	Staves Stages   `json:"stages" bson:"stages"`
}

// The Recipes is a convenience wrapper for all recipe types
type Recipes []Recipe

// The StorageServiceRecipes is an interface for the backend storage service
type StorageServiceRecipes interface {
	GetRecipes(filterTags []string, filterAuthor string) (Recipes, error)
	GetRecipe(id RecipeID) (Recipe, error)
	AddRecipe(recipe *Recipe) error
	UpdateRecipe(recipe *Recipe) error
}

// MarshalJSON implements json.Marshaller.
func (r RecipeID) MarshalJSON() ([]byte, error) {
	f := guid.UUID(r).String()
	return json.Marshal(f)
}

// UnmarshalJSON implements json.Unmarshaller.
func (r *RecipeID) UnmarshalJSON(b []byte) error {
	var decoded string

	jsonErr := json.Unmarshal(b, &decoded)

	if jsonErr == nil {
		temp, parseErr := guid.Parse(decoded)

		if parseErr != nil {
			return parseErr
		}
		*r = RecipeID(temp)
	}
	return jsonErr
}

// MarshalBSON implements bson.Getter.
func (r RecipeID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.String, bsoncore.AppendString(nil, guid.UUID(r).String()), nil
}

// UnmarshalBSON implements bson.Setter.
func (r *RecipeID) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	var s string

	if t != bsontype.String {
		return fmt.Errorf("Can only decode strings")
	}

	s, _, ok := bsoncore.ReadString(raw)
	if !ok {
		return fmt.Errorf("Failed to unmarshal recipeID")
	}

	uuid, err := guid.Parse(s)
	if err == nil {
		*r = RecipeID(uuid)
	}

	return err
}
