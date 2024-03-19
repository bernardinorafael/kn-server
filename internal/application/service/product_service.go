package service

import (
	"time"

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

	_, err := ps.s.prodRepo.GetByName(p.Name)
	if err != nil {
		ps.s.log.Error("cannot find product by name", err)
		return err
	}

	product := entity.Product{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ps.s.prodRepo.Save(product)
	if err != nil {
		ps.s.log.Error("error saving product in DB", err)
		return err
	}

	return nil
}

func (ps *ProductService) GetByID(id uint) (*entity.Product, error) {
	ps.s.log.Info("Process started")
	defer ps.s.log.Info("Process finished")

	return nil, nil
}
