package user

type RegisterRequest struct {
	FristName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Phone     string `json:"phone" bson:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" bson:"old_password"`
	NewPassword string `json:"new_password" bson:"new_password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" bson:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" bson:"token"`
	NewPassword string `json:"new_password" bson:"new_password"`
}