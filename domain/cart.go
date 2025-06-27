package domain

type CartItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

type Source struct {
	Username string     `json:"username"`
	Cart     []CartItem `json:"cart"`
}

type Hit struct {
	Source Source `json:"_source"`
}

type Hits struct {
	Hits []Hit `json:"hits"`
}

type ESResponse struct {
	Hits Hits `json:"hits"`
}
