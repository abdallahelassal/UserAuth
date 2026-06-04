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
type FindByIDOutput struct{
	Email string
	UserName string
	IsActive bool
	Roles 	[]RoleOutput
}