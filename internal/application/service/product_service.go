package service

import (
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
)

type ProductService struct {
	s *service
}

func newProductService(service *service) contract.ProductService {
	return &ProductService{s: service}
}

func (ps *ProductService) SaveProduct(p entity.Product) error {
	ps.s.log.Info("Process started")
	defer ps.s.log.Info("Process finished")

	// TODO: verify the availability(name)
	// TODO: generate the SKU
	// TODO: generate the slug
	// TODO: save in db

	return nil
}

func (ps *ProductService) GetByID(id uint) (*entity.Product, error) {
	ps.s.log.Info("Process started")
	defer ps.s.log.Info("Process finished")

	return nil, nil
}
