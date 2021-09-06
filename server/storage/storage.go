package storage

type Calculation struct {
	Expression string
	Result string
}

type HistoryStorage interface {
	StoreCalculation(c Calculation) error
	GetHistory() ([]Calculation, error)
}

type InMemoryHistoryStorage struct {
	Calculations []Calculation
}

func (s *InMemoryHistoryStorage) StoreCalculation(c Calculation) error {
	s.Calculations = append(s.Calculations, c)
	return nil
}

func (s *InMemoryHistoryStorage) GetHistory() ([]Calculation, error) {
	return s.Calculations, nil
}