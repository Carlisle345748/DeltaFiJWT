package types

type CreateUserInput struct {
	Email     string `form:"email" json:"email" binding:"required"`
	FirstName string `form:"firstname" json:"firstName" binding:"required"`
	LastName  string `form:"lastname" json:"lastName" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
}

type UpdateUserInput struct {
	ID        uint   `form:"id" json:"id" binding:"required"`
	FirstName string `form:"firstName" json:"firstName" binding:"required"`
	LastName  string `form:"lastName" json:"lastName" binding:"required"`
}

type DeleteUserInput struct {
	ID uint `form:"id" json:"id" binding:"required"`
}
