package services

import (
	"fmt"

	h "app.ddcli.datnn/application/lib"
)

type StorageInformation struct {
	StoreID          string
	StoreName        string
	RootFolder       string
	TotalStorage     float64
	UsedStorage      float64
	IsLimitedStorage bool
}

var BAR_STYLE = []rune{'|', '▌', '░'}

func (s *StorageInformation) String() string {
	return fmt.Sprintf("[%s]\n", s.StoreID) +
		fmt.Sprintf("Used: %.2f MB\nTotal: %.2f MB\n", s.UsedStorage, s.TotalStorage)
}

func (s *StorageInformation) GetRemainingStorage() float64 {
	return s.TotalStorage - s.UsedStorage
}

type ServiceResource interface {
	New(*h.Config) error
	GetInformation() (*StorageInformation, error)
}
