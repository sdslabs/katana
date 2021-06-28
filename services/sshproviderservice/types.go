package sshproviderservice

type CTFTeam struct {
	Index    int    `json:"id" bson:"password" binding:"required"`
	Name     string `json:"name" bson:"username" binding:"required"`
	PodName  string `json:"podname" bson:"podname" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}
