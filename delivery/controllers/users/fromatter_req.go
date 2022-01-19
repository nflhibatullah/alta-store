package user

type RegisterUserRequestFormat struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type PutUserRequestFormat struct {
	Email    string `json:"email" form:"email"`
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

type DeleteRequestFormat struct {
	Password string `json:"password" form:"password"`
}
