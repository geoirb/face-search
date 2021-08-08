package mongo

type faceSearch struct {
	Status    string    `bson:"status"`
	Error     string    `bson:"error"`
	UUID      string    `bson:"uuid"`
	PhotoHash string    `bson:"photo_hash"`
	Profiles  []profile `bson:"response"`
	UpdateAt  int64     `bson:"update_at"`
	CreateAt  int64     `bson:"create_at"`
}

type profile struct {
	FullName    string `bson:"full_name"`
	LinkProfile string `bson:"link_profile"`
	LinkPhoto   string `bson:"link_photo"`
	Confidence  string `bson:"confidence"`
}
