package dto

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Account struct {
	Balance    float32 `json:"current"`
	TotalSpent float32 `json:"withdrawn"`
}

type ParamsWithdraw struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type OrderInfo struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

type Order struct {
	Number     string
	Status     string
	Accrual    float32
	UploadedAt string
	UserID     int64
}
