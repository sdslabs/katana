package types

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
}
