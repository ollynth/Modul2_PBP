package models

type Users struct {
	ID      int    `json: "id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Users  `json:"data"`
}

type UsersResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Users `json:"data"`
}

type Products struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    Products `json:"data"`
}

type ProductsResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Products `json:"data"`
}

type Transactions struct {
	ID        int `json:"id"`
	UserID    int `json:"userid"`
	ProductID int `json:"productid"`
	Quantity  int `json:"quantity"`
}

type TransactionResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    Transactions `json:"data"`
}

type TransactionsResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []Transactions `json:"data"`
}
