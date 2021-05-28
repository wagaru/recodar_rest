package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Email         string             `json:"email"`
	Name          string             `json:"name,omitempty"`
	Password      string             `json:"password,omitempty"`
	Picture       string             `json:"picture,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	BindingSource string             `json:"binding_source" bson:"binding_source"`
	AccessToken   string             `json:"access_token" bson:"access_token,omitempty"`
	RefreshToken  string             `json:"refresh_token" bson:"refresh_token,omitempty"`
	JWT           string             `json:"jwt,omitempty"`
	LineUserID    string             `json:"line_user_id" bson:"line_user_id,omitempty"`
	GoogleUserID  string             `json:"google_user_id" bson:"google_user_id,omitempty"`
	// Expired
}

type Accident struct {
	ID              primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Place           Place              `json:"place" binding:"required"`
	ApproxTime      time.Time          `json:"approx_time" bson:"approx_time" binding:"lt,required"`
	Description     string             `json:"description"`
	AccidentObjects []AccidentObject   `json:"accident_objects" bson:"accident_objects" binding:"required"`
	Videos          []Video            `json:"videos" binding:"required"`
	CreatedAt       *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	// CreatorID       primitive.ObjectID `json:"creator_id"`
}
type Place struct {
	Level1      string  `json:"level_1" bson:"level_1" binding:"required"`
	Level2      string  `json:"level_2" bson:"level_2" binding:"required"`
	GeoLocation GeoJSON `json:"geo_location" bson:"geo_location" binding:"required"`
}

type GeoJSON struct {
	Type        string    `json:"type" binding:"required"`
	Coordinates []float64 `json:"coordinates" binding:"required,dive,required"`
}

type Video struct {
	// ID        primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	MediaID   string `json:"media_id" bson:"media_id" binding:"required"`
	MediaURL  string `json:"media_url" bson:"media_url" binding:"required"`
	MediaType string `json:"media_type" bson:"media_type" binding:"required"`
}

type AccidentObject string

const (
	Automobile AccidentObject = "Automobile"
	Bicycle                   = "Bicycle"
	Motorcycle                = "Motorcycle"
	Pedestrian                = "Pedestrian"
	Others                    = "Others"
)
