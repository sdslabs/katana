package types

import "time"

type AdminUser struct {
	Username string `json:"username" bson:"username" binding:"required" `
	Password string `json:"password" bson:"password" binding:"required"`
}

type CTFTeam struct {
	Index     int    `json:"id" bson:"id" binding:"required"`
	Name      string `json:"name" bson:"username" binding:"required"`
	PodName   string `json:"podname" bson:"podname" binding:"required"`
	Password  string `json:"password" bson:"password" binding:"required"`
	PublicKey string `json:"publicKey" bson:"publicKey" binding:"required"` // TODO : initialize
	Score     int    `json:"score" bson:"score" binding:"required"`
}

type Challenge struct {
	ID     int    `json:"id" bson:"id" binding:"required"`
	TeamID int    `json:"teamid" bson:"teamid" binding:"required"`
	Name   string `json:"name" bson:"name" binding:"required"`
	Points int    `json:"points" bson:"points" binding:"required"`
}

type Flag struct {
	Value       string    `json:"value" bson:"value" binding:"required"`
	ChallengeID int       `json:"challengeid" bson:"challengeid" binding:"required"`
	TeamID      int       `json:"teamid" bson:"teamid" binding:"required"`
	CreatedAt   time.Time `json:"createtime" bson:"createtime" binding:"required"`
}

type Submission struct {
	Submitter   int    `json:"submitter" bson:"submitter" binding:"required"`
	ChallengeID int    `json:"challengeid" bson:"challengeid" binding:"required"`
	Flag        string `json:"flag" bson:"flag" binding:"required"`
}
