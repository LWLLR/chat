package service

//Login 登录
func Login(data []rune) {
	name := ""
	for _, i := range data {
		name += string(i)
	}
}
