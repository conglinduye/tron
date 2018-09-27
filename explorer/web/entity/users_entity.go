package entity

type ApiUsers struct {
	Id 			uint64		`json:"id"`
	Username	string 		`json:"username"`
	Password	string		`json:"password"`
}

type ApiToken struct {
	Token 		string 		`json:"token"`
}
