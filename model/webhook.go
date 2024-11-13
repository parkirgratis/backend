package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Userdomyikado struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name                 string             `bson:"name,omitempty" json:"name,omitempty"`
	PhoneNumber          string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Email                string             `bson:"email,omitempty" json:"email,omitempty"`
	GithubUsername       string             `bson:"githubusername,omitempty" json:"githubusername,omitempty"`
	GitlabUsername       string             `bson:"gitlabusername,omitempty" json:"gitlabusername,omitempty"`
	GitHostUsername      string             `bson:"githostusername,omitempty" json:"githostusername,omitempty"`
	Poin                 float64            `bson:"poin,omitempty" json:"poin,omitempty"`
	GoogleProfilePicture string             `bson:"googleprofilepicture,omitempty" json:"picture,omitempty"`
	Team                 string             `json:"team,omitempty" bson:"team,omitempty"`
	Scope                string             `json:"scope,omitempty" bson:"scope,omitempty"`
	Section              string             `json:"section,omitempty" bson:"section,omitempty"`
	Chapter              string             `json:"chapter,omitempty" bson:"chapter,omitempty"`
	LinkedDevice         string             `json:"linkeddevice,omitempty" bson:"linkeddevice,omitempty"`
	JumlahAntrian        int                `json:"jumlahantrian,omitempty" bson:"jumlahantrian,omitempty"`
}