package usecase



type CreateUserInput struct {
	UserName string
	Email	 string
	Password string 
}
type LoginUserInput struct{
	Email string
	Password string
}