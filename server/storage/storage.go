package storage

type Calculation struct {
	Expression string
	Result string
}

type HistoryStorage interface {
	StoreCalculation(c Calculation) error
	GetHistory() ([]Calculation, error)
}