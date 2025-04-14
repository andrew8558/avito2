//go:generate mockgen -source ./service.go -destination=./mocks/service.go -package=mock_service
package service

import (
	"avito2/internal/errors"
	"avito2/internal/model"
	"avito2/internal/repository"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type Service interface {
	CreatePvz(ctx context.Context, city model.City) (*model.Pvz, error)
	CloseLastReception(ctx context.Context, pvzId uuid.UUID) (*model.Reception, error)
	DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error
	CreateReception(ctx context.Context, pvzId uuid.UUID) (*model.Reception, error)
	AddProduct(ctx context.Context, pvzId uuid.UUID, productType model.ProductType) (*model.Product, error)
	GetPvzInfo(ctx context.Context, startDate, endDate time.Time, page, limit int32) (*model.GetPvzInfoResponse, error)
}

type Svc struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Svc {
	return &Svc{
		repo: repo,
	}
}

func (s *Svc) CreatePvz(ctx context.Context, city model.City) (*model.Pvz, error) {
	pvz, err := s.repo.CreatePvz(ctx, city)
	if err != nil {
		log.Println("failed to create pvz with err:", err)
		return nil, err
	}
	return pvz, nil
}

func (s *Svc) CloseLastReception(ctx context.Context, pvzId uuid.UUID) (*model.Reception, error) {
	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return nil, err
	}

	_, err = s.repo.GetPvz(ctx, tx, pvzId)
	if err != nil {
		if err != errors.ErrPvzDoesNotExist {
			log.Println("failed to get pvz with err:", err)
		}
		s.repo.RollbackTx(ctx, tx)
		return nil, err
	}

	reception, err := s.repo.UpdateLastReceptionStatus(ctx, tx, pvzId)
	if err != nil {
		log.Println("failed to close last reception with err:", err)
		s.repo.RollbackTx(ctx, tx)
		return nil, err
	}
	s.repo.CommitTx(ctx, tx)

	if reception == nil {
		return nil, errors.ErrReceptionInProgressDoesNotExist
	}
	return reception, nil
}

func (s *Svc) DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error {
	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return err
	}

	_, err = s.repo.GetPvz(ctx, tx, pvzId)
	if err != nil {
		if err != errors.ErrPvzDoesNotExist {
			log.Println("failed to get pvz with err:", err)
		}
		s.repo.RollbackTx(ctx, tx)
		return err
	}

	curReception, err := s.repo.GetCurrentReception(ctx, tx, pvzId)
	if err != nil {
		log.Println("failed to get current reception with err:", err)
		s.repo.RollbackTx(ctx, tx)
		return err
	}

	if curReception != nil {
		err = s.repo.DeleteLastProduct(ctx, tx, curReception.Id)
		if err != nil {
			if err != errors.ErrNoProductToDelete {
				log.Println("failed to delete last product from current reception with err:", err)
				s.repo.RollbackTx(ctx, tx)
				return err
			}
			s.repo.CommitTx(ctx, tx)
			return err
		}
		s.repo.CommitTx(ctx, tx)
		return nil
	}

	s.repo.CommitTx(ctx, tx)
	return errors.ErrReceptionInProgressDoesNotExist
}

func (s *Svc) CreateReception(ctx context.Context, pvzId uuid.UUID) (*model.Reception, error) {
	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return nil, err
	}

	_, err = s.repo.GetPvz(ctx, tx, pvzId)
	if err != nil {
		if err != errors.ErrPvzDoesNotExist {
			log.Println("failed to get pvz with err:", err)
		}
		s.repo.RollbackTx(ctx, tx)
		return nil, err
	}

	curReception, err := s.repo.GetCurrentReception(ctx, tx, pvzId)
	if err != nil {
		log.Println("failed to get current reception with err:", err)
		s.repo.RollbackTx(ctx, tx)
		return nil, err
	}

	if curReception == nil {
		reception, err := s.repo.CreateReception(ctx, tx, pvzId)
		if err != nil {
			log.Println("failed to create reception with err:", err)
			s.repo.RollbackTx(ctx, tx)
			return nil, err
		}
		s.repo.CommitTx(ctx, tx)
		return reception, nil
	}

	s.repo.CommitTx(ctx, tx)
	return nil, errors.ErrReceptionInProgressAlreadyExists
}

func (s *Svc) AddProduct(ctx context.Context, pvzId uuid.UUID, productType model.ProductType) (*model.Product, error) {
	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return nil, err
	}

	_, err = s.repo.GetPvz(ctx, tx, pvzId)
	if err != nil {
		if err != errors.ErrPvzDoesNotExist {
			log.Println("failed to get pvz with err:", err)
		}
		s.repo.RollbackTx(ctx, tx)
		return nil, err
	}

	curReception, err := s.repo.GetCurrentReception(ctx, tx, pvzId)
	if err != nil {
		log.Println("failed to get current reception with err:", err)
		s.repo.RollbackTx(ctx, tx)
		return nil, err
	}

	if curReception != nil {
		product, err := s.repo.AddProduct(ctx, tx, curReception.Id, productType)
		if err != nil {
			log.Println("failed to add product to current reception with err:", err)
			s.repo.RollbackTx(ctx, tx)
			return nil, err
		}
		s.repo.CommitTx(ctx, tx)
		return product, nil
	}

	s.repo.CommitTx(ctx, tx)
	return nil, errors.ErrReceptionInProgressDoesNotExist
}

func (s *Svc) GetPvzInfo(ctx context.Context, startDate, endDate time.Time, page, limit int32) (*model.GetPvzInfoResponse, error) {
	tx, err := s.repo.BeginTransaction(ctx, &pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		log.Println("failed to begin tx with err:", err)
		return nil, err
	}

	offset := (page - 1) * limit

	receptions, err := s.repo.GetReceptionsForPeriod(ctx, tx, startDate, endDate, offset, limit)
	if err != nil {
		log.Println("failed to get receptions for period with err:", err)
		s.repo.RollbackTx(ctx, tx)
		return nil, err
	}

	pvzMap := make(map[uuid.UUID]*model.PvzInfo, limit)
	for _, reception := range receptions {
		products, err := s.repo.GetProductsInReception(ctx, tx, reception.Id)
		if err != nil {
			log.Println("failed to get products with err:", err)
			s.repo.RollbackTx(ctx, tx)
			return nil, err
		}

		pvzInfo, ok := pvzMap[reception.PvzId]
		if !ok {
			pvz, err := s.repo.GetPvz(ctx, tx, reception.PvzId)
			if err != nil {
				log.Println("failed to get pvz with err:", err)
				s.repo.RollbackTx(ctx, tx)
				return nil, err
			}

			pvzInfo = &model.PvzInfo{
				Pvz:        *pvz,
				Receptions: []model.ReceptionInfo{},
			}

			pvzMap[reception.PvzId] = pvzInfo
		}
		pvzInfo.Receptions = append(pvzInfo.Receptions, model.ReceptionInfo{Reception: reception, Products: products})
	}
	s.repo.CommitTx(ctx, tx)

	res := &model.GetPvzInfoResponse{}
	for _, v := range pvzMap {
		res.PvzList = append(res.PvzList, *v)
	}

	return res, nil
}
