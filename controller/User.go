package controller

type UserController struct {
	BasicController
}

func (u *UserController) Login() {
	u.ServeJSON()
}
