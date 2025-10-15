package dto

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Account struct {
	Balance    float64 `json:"current"`
	TotalSpent float64 `json:"withdrawn"`
}

type ParamsWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type OrderInfo struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}
