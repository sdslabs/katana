package types

type AdminUser struct {
	Username string `json:"username" bson:"username" binding:"required" `
	Password string `json:"password" bson:"password" binding:"required"`
}

type CTFTeam struct {
	Index      int         `json:"id" bson:"id" binding:"required"`
	Name       string      `json:"name" bson:"username" binding:"required"`
	PodName    string      `json:"podname" bson:"podname" binding:"required"`
	Password   string      `json:"password" bson:"password" binding:"required"`
	Challenges []Challenge `json:"challenges"`
	Score	   int 		   `json:"score"`
	PublicKey  string      `json:"publicKey" bson:"publicKey" binding:"required"`
}

type Challenge struct {
	ChallengeName string  `json:"ChallengeName"`
	Uptime        float64 `json:"uptime"`
	Attacks       int     `json:"attacks"`
	Defenses      int     `json:"defenses"`
	Flag          string  `json:"flag"`
}
