package domain

type Client struct {
	ID      uint8
	Limit   uint64
	Balance int64
}

func NewClient(id uint8, limit uint64, balance int64) (*Client, error) {

	return &Client{
		ID:      id,
		Limit:   limit,
		Balance: balance,
	}, nil
}
