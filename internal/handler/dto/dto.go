package dto

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Account struct {
	Balance    float64 `json:"current"`
	TotalSpent float64 `json:"withdrawn"`
}
