package models

type Frequency string

const ( // this is the frequency of the subscription
	Hourly Frequency = "hourly"
	Daily  Frequency = "daily"
)

type Subscription struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	City      string    `json:"city"`
	Frequency Frequency `json:"frequency"`
	Confirmed bool      `json:"confirmed"`
}
