package profile

type CreateProfileRequest struct {
	UserID   string `form:"user_id" bson:"user_id"`
	Gender   string `form:"gender" bson:"gender"`
	BirthDay string `form:"birth_day" bson:"birth_day"`
	Address  string `form:"address" bson:"address"`
	Bio      string `form:"bio" bson:"bio"`
}

type UpdateProfileRequest struct {
	UserID   string `form:"user_id" bson:"user_id"`
	Gender   string `form:"gender" bson:"gender"`
	BirthDay string `form:"birth_day" bson:"birth_day"`
	PublicID string `form:"public_id" bson:"public_id"`
	Address  string `form:"address" bson:"address"`
	Bio      string `form:"bio" bson:"bio"`
}