package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Provider ...
type Provider struct {
	ID         bson.ObjectId     `bson:"_id" json:"id"`
	Name       string            `bson:"name" json:"name"`
	Type       string            `bson:"type" json:"type"`
	Connection map[string]string `bson:"connection" json:"connection"`
	CreatedAt  time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time         `bson:"updated_at" json:"updated"`
}

type AWSProvider struct {
	Image  string `bson:"image" json:"image"`
	Size   string `bson:"size" json:"size"`
	VMname string `bson:"vm_name" json:"vm_name"`
}

