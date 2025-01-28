package base

type Track struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	Password       string `json:"password"`
	FirstName      string `json:"firstName"`
	SecondName       string `json:"secondName"`
	LastName        string `json:"lastName"`
	CompanyID       string `json:"companyID"`
}