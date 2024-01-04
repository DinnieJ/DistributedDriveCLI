package services

type IStoreResource interface {
	GetRemainingStorage() uint32
}

type StoreInformation struct {
	StoreSource  string
	StoreID      string
	StoreName    string
	TotalStorage uint32
	UsedStorage  uint32
}

func (s *StoreInformation) RemainingStorage() uint32 {
	return s.TotalStorage - s.UsedStorage
}
