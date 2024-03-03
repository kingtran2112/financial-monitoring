package importing

type Spending struct {
	Id       string `csv:"Id"`
	Date     string `csv:"Date"`
	Group    string `csv:"Group"`
	Amount   int32  `csv:"Amount"`
	Currency string `csv:"Currency"`
	Note     string `csv:"Note"`
	Wallet   string `csv:"Wallet"`
	Type     SpendingType
}

type SpendingType string

const (
	INCOME  SpendingType = "INCOME"
	EXPENSE SpendingType = "EXPENSE"
)

func (t SpendingType) IsValid() bool {
	switch t {
	case INCOME, EXPENSE:
		return true
	default:
		return false
	}
}

func (t SpendingType) String() string {
	return string(t)
}
