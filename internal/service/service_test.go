package service

import (
	customErrors "avito2/internal/errors"
	"avito2/internal/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreatePvz(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		city        = model.CityMoscow
		expectedPvz = &model.Pvz{City: city}
		dbErr       = errors.New("failed to create pvz")
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().CreatePvz(gomock.Any(), gomock.Any()).Return(expectedPvz, nil)

		pvz, err := s.svc.CreatePvz(ctx, city)

		require.NoError(t, err)
		assert.Equal(t, expectedPvz, pvz)
	})

	t.Run("db error", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().CreatePvz(gomock.Any(), gomock.Any()).Return(nil, dbErr)

		_, err := s.svc.CreatePvz(ctx, city)

		require.EqualError(t, err, dbErr.Error())
	})
}

func Test_CloseLastReception(t *testing.T) {
	t.Parallel()

	var (
		ctx               = context.Background()
		pvzId             = uuid.New()
		expectedReception = &model.Reception{PvzId: pvzId}
		dbErr             = errors.New("db error")
	)

	t.Run("succes", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().UpdateLastReceptionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedReception, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		rec, err := s.svc.CloseLastReception(ctx, pvzId)

		require.NoError(t, err)
		assert.Equal(t, expectedReception, rec)
	})

	t.Run("pvz doed not exist", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrPvzDoesNotExist)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.CloseLastReception(ctx, pvzId)

		require.EqualError(t, err, customErrors.ErrPvzDoesNotExist.Error())
	})

	t.Run("no reception in progress", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().UpdateLastReceptionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.CloseLastReception(ctx, pvzId)

		require.EqualError(t, err, customErrors.ErrReceptionInProgressDoesNotExist.Error())
	})
	t.Run("failed begin transaction", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, dbErr)

		_, err := s.svc.CloseLastReception(ctx, pvzId)

		require.Error(t, err)
	})
	t.Run("failed to update reception status", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().UpdateLastReceptionStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.CloseLastReception(ctx, pvzId)

		require.Error(t, err)
	})
}

func Test_DeleteLastProduct(t *testing.T) {
	t.Parallel()

	var (
		ctx   = context.Background()
		pvzId = uuid.New()
		rec   = &model.Reception{PvzId: pvzId}
		dbErr = errors.New("db error")
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(rec, nil)
		s.mockRepo.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.DeleteLastProduct(ctx, pvzId)

		require.NoError(t, err)
	})

	t.Run("pvz does not exist", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrPvzDoesNotExist)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.DeleteLastProduct(ctx, pvzId)

		require.EqualError(t, err, customErrors.ErrPvzDoesNotExist.Error())
	})

	t.Run("no reception in progress", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.DeleteLastProduct(ctx, pvzId)

		require.EqualError(t, err, customErrors.ErrReceptionInProgressDoesNotExist.Error())
	})

	t.Run("no product to delete", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(rec, nil)
		s.mockRepo.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any(), gomock.Any()).Return(customErrors.ErrNoProductToDelete)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.DeleteLastProduct(ctx, pvzId)

		require.EqualError(t, err, customErrors.ErrNoProductToDelete.Error())
	})
	t.Run("failed to get current reception", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.DeleteLastProduct(ctx, pvzId)

		require.Error(t, err)
	})
	t.Run("failed begin transaction", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, dbErr)

		err := s.svc.DeleteLastProduct(ctx, pvzId)

		require.Error(t, err)
	})
	t.Run("failed to delete product", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(rec, nil)
		s.mockRepo.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		err := s.svc.DeleteLastProduct(ctx, pvzId)

		require.Error(t, err)
	})
}

func Test_CreateReception(t *testing.T) {
	t.Parallel()

	var (
		ctx         = context.Background()
		pvzId       = uuid.New()
		expectedRec = &model.Reception{PvzId: pvzId}
		curRec      = &model.Reception{}
		dbErr       = errors.New("db error")
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CreateReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedRec, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		rec, err := s.svc.CreateReception(ctx, pvzId)

		require.NoError(t, err)
		assert.Equal(t, expectedRec, rec)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CreateReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedRec, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		rec, err := s.svc.CreateReception(ctx, pvzId)

		require.NoError(t, err)
		assert.Equal(t, expectedRec, rec)
	})
	t.Run("pvz does not exist", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrPvzDoesNotExist)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.CreateReception(ctx, pvzId)

		require.EqualError(t, err, customErrors.ErrPvzDoesNotExist.Error())
	})
	t.Run("reception in progress already exist", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(curRec, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.CreateReception(ctx, pvzId)

		require.EqualError(t, err, customErrors.ErrReceptionInProgressAlreadyExists.Error())
	})
	t.Run("failed to get current reception", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.CreateReception(ctx, pvzId)

		require.Error(t, err)
	})
	t.Run("failed begin transaction", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, dbErr)

		_, err := s.svc.CreateReception(ctx, pvzId)

		require.Error(t, err)
	})
	t.Run("failed to create reception", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CreateReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.CreateReception(ctx, pvzId)

		require.Error(t, err)
	})
}

func Test_AddProduct(t *testing.T) {
	t.Parallel()

	var (
		ctx             = context.Background()
		pvzId           = uuid.New()
		productType     = model.ProductTypeClothes
		expectedProduct = &model.Product{Type: productType}
		rec             = &model.Reception{PvzId: pvzId}
		dbErr           = errors.New("db error")
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(rec, nil)
		s.mockRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedProduct, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		product, err := s.svc.AddProduct(ctx, pvzId, productType)

		require.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})
	t.Run("pvz does not exist", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrPvzDoesNotExist)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.AddProduct(ctx, pvzId, productType)

		require.EqualError(t, err, customErrors.ErrPvzDoesNotExist.Error())
	})
	t.Run("no reception in progress", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.AddProduct(ctx, pvzId, productType)

		require.EqualError(t, err, customErrors.ErrReceptionInProgressDoesNotExist.Error())
	})
	t.Run("failed to get current reception", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.AddProduct(ctx, pvzId, productType)

		require.Error(t, err)
	})
	t.Run("failed begin transaction", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, dbErr)

		_, err := s.svc.AddProduct(ctx, pvzId, productType)

		require.Error(t, err)
	})
	t.Run("failed to add product", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetCurrentReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(rec, nil)
		s.mockRepo.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.AddProduct(ctx, pvzId, productType)

		require.Error(t, err)
	})
}

func Test_GetPvzInfo(t *testing.T) {
	t.Parallel()

	var (
		ctx           = context.Background()
		page          = int32(1)
		limit         = int32(10)
		startDate     = time.Time{}
		endDate       = time.Now()
		pvzId         = uuid.New()
		dbErr         = errors.New("db error")
		pvz           = model.Pvz{Id: pvzId}
		reception     = model.Reception{PvzId: pvzId}
		receptions    = []model.Reception{reception}
		product1      = model.Product{Type: model.ProductTypeClothes}
		product2      = model.Product{Type: model.ProductTypeElectronics}
		products      = []model.Product{product1, product2}
		receptionInfo = model.ReceptionInfo{
			Reception: reception,
			Products:  products,
		}
		expectedRes = &model.GetPvzInfoResponse{
			PvzList: []model.PvzInfo{
				{
					Pvz:        pvz,
					Receptions: []model.ReceptionInfo{receptionInfo},
				},
			},
		}
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetReceptionsForPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(receptions, nil)
		s.mockRepo.EXPECT().GetProductsInReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(products, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(&pvz, nil)
		s.mockRepo.EXPECT().CommitTx(gomock.Any(), gomock.Any()).Return()

		res, err := s.svc.GetPvzInfo(ctx, startDate, endDate, page, limit)

		require.NoError(t, err)
		assert.Equal(t, expectedRes, res)
	})
	t.Run("failed begin transaction", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, dbErr)

		_, err := s.svc.GetPvzInfo(ctx, startDate, endDate, page, limit)

		require.Error(t, err)
	})
	t.Run("failed to get receptions for period", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetReceptionsForPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.GetPvzInfo(ctx, startDate, endDate, page, limit)

		require.Error(t, err)
	})
	t.Run("failed to get products in reception", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetReceptionsForPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(receptions, nil)
		s.mockRepo.EXPECT().GetProductsInReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.GetPvzInfo(ctx, startDate, endDate, page, limit)

		require.Error(t, err)
	})
	t.Run("failed to get pvz", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()
		s.mockRepo.EXPECT().BeginTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
		s.mockRepo.EXPECT().GetReceptionsForPeriod(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(receptions, nil)
		s.mockRepo.EXPECT().GetProductsInReception(gomock.Any(), gomock.Any(), gomock.Any()).Return(products, nil)
		s.mockRepo.EXPECT().GetPvz(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, dbErr)
		s.mockRepo.EXPECT().RollbackTx(gomock.Any(), gomock.Any()).Return()

		_, err := s.svc.GetPvzInfo(ctx, startDate, endDate, page, limit)

		require.Error(t, err)
	})
}
