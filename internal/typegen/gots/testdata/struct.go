// User is an account holder.
type User struct {
	ID      int64             `json:"id"`
	Email   *string           `json:"email"`
	Tags    []string          `json:"tags,omitempty"`
	Meta    map[string]string `json:"meta"`
	Created time.Time         `json:"created_at"`
	Address Address           `json:"address"`
	secret  string
}

type Address struct {
	City string `json:"city"`
}
