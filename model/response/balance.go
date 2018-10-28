package response

type Balance struct {
	TotalBalance     float64 `bson:"totalBalance" json:"totalBalance"`
	AvailableBalance float64 `bson:"availableBalance" json:"availableBalance"`
}
