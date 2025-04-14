package repository

import (
	"avito2/internal/db"
	"avito2/internal/errors"
	"avito2/internal/model"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type Repo struct {
	db db.DBops
}

type Repository interface {
	BeginTransaction(ctx context.Context, options *pgx.TxOptions) (pgx.Tx, error)
	RollbackTx(ctx context.Context, tx pgx.Tx)
	CommitTx(ctx context.Context, tx pgx.Tx)
	CreatePvz(ctx context.Context, city model.City) (*model.Pvz, error)
	GetPvz(ctx context.Context, tx pgx.Tx, pvzId uuid.UUID) (*model.Pvz, error)
	UpdateLastReceptionStatus(ctx context.Context, tx pgx.Tx, pvzId uuid.UUID) (*model.Reception, error)
	GetCurrentReception(ctx context.Context, tx pgx.Tx, pvzId uuid.UUID) (*model.Reception, error)
	CreateReception(ctx context.Context, tx pgx.Tx, pvzId uuid.UUID) (*model.Reception, error)
	AddProduct(ctx context.Context, tx pgx.Tx, receptionId uuid.UUID, productType model.ProductType) (*model.Product, error)
	DeleteLastProduct(ctx context.Context, tx pgx.Tx, receptionId uuid.UUID) error
	GetReceptionsForPeriod(ctx context.Context, tx pgx.Tx, startDate, endDate time.Time, offset, limit int32) ([]model.Reception, error)
	GetProductsInReception(ctx context.Context, tx pgx.Tx, receptionId uuid.UUID) ([]model.Product, error)
}

func NewRepository(database db.DBops) *Repo {
	return &Repo{db: database}
}
func (r *Repo) BeginTransaction(ctx context.Context, options *pgx.TxOptions) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, options)
}

func (r *Repo) RollbackTx(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil {
		log.Println("failed to rollback tx wih err:", err)
	}
}

func (r *Repo) CommitTx(ctx context.Context, tx pgx.Tx) {
	if err := tx.Commit(ctx); err != nil {
		log.Println("failed to commit tx wih err:", err)
	}
}

func (r *Repo) CreatePvz(ctx context.Context, city model.City) (*model.Pvz, error) {
	regDate := time.Now()
	row := r.db.ExecQueryRow(ctx, "INSERT INTO pvz (registration_date, city) VALUES ($1, $2) RETURNING id, registration_date, city", regDate, city)

	var pvz model.Pvz
	if err := row.Scan(&pvz.Id, &pvz.RegistrationDate, &pvz.City); err != nil {
		return nil, err
	}

	return &pvz, nil
}

func (r *Repo) GetPvz(ctx context.Context, tx pgx.Tx, pvzId uuid.UUID) (*model.Pvz, error) {
	var pvz model.Pvz
	err := tx.QueryRow(ctx, "SELECT * FROM pvz WHERE id = $1", pvzId).Scan(&pvz.Id, &pvz.City, &pvz.RegistrationDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrPvzDoesNotExist
		}
		return nil, err
	}

	return &pvz, nil
}

func (r *Repo) UpdateLastReceptionStatus(ctx context.Context, tx pgx.Tx, pvzId uuid.UUID) (*model.Reception, error) {
	row := tx.QueryRow(ctx, "UPDATE receptions SET status = $1 WHERE pvz_id = $2 AND status = $3 RETURNING id, date_time, pvz_id, status",
		model.ReceptionStatusClose, pvzId, model.ReceptionStatusInProgress)
	var reception model.Reception
	if err := row.Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &reception, nil
}

func (r *Repo) GetCurrentReception(ctx context.Context, tx pgx.Tx, pvzId uuid.UUID) (*model.Reception, error) {
	var reception model.Reception
	err := tx.QueryRow(ctx, "SELECT * FROM receptions WHERE pvz_id = $1 AND status = $2 FOR UPDATE",
		pvzId, model.ReceptionStatusInProgress).Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reception, nil
}

func (r *Repo) CreateReception(ctx context.Context, tx pgx.Tx, pvzId uuid.UUID) (*model.Reception, error) {
	dateTime := time.Now()
	var reception model.Reception
	err := tx.QueryRow(ctx, "INSERT INTO receptions (date_time, pvz_id, status) VALUES ($1, $2, $3) RETURNING id, date_time, pvz_id, status",
		dateTime, pvzId, model.ReceptionStatusInProgress).Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)

	if err != nil {
		return nil, err
	}

	return &reception, nil
}

func (r *Repo) AddProduct(ctx context.Context, tx pgx.Tx, receptionId uuid.UUID, productType model.ProductType) (*model.Product, error) {
	dateTime := time.Now()
	var product model.Product
	err := tx.QueryRow(ctx, "INSERT INTO products (date_time, type, reception_id) VALUES ($1, $2, $3) RETURNING id, date_time, type, reception_id",
		dateTime, productType, receptionId).Scan(&product.Id, &product.DateTime, &product.Type, &product.ReceptionId)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *Repo) DeleteLastProduct(ctx context.Context, tx pgx.Tx, receptionId uuid.UUID) error {
	commandTag, err := tx.Exec(ctx, "DELETE FROM products WHERE id = (SELECT id FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1)", receptionId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.ErrNoProductToDelete
	}
	return nil
}

func (r *Repo) GetReceptionsForPeriod(ctx context.Context, tx pgx.Tx, startDate, endDate time.Time, offset, limit int32) ([]model.Reception, error) {
	rows, err := tx.Query(ctx, "SELECT * FROM receptions WHERE date_time BETWEEN $1 AND $2 ORDER BY date_time DESC LIMIT $3 OFFSET $4 FOR UPDATE", startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receptions := []model.Reception{}
	for rows.Next() {
		var reception model.Reception
		err := rows.Scan(&reception.Id, &reception.DateTime, &reception.PvzId, &reception.Status)
		if err != nil {
			return nil, err
		}
		receptions = append(receptions, reception)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return receptions, nil
}

func (r *Repo) GetProductsInReception(ctx context.Context, tx pgx.Tx, receptionId uuid.UUID) ([]model.Product, error) {
	rows, err := tx.Query(ctx, "SELECT * FROM products WHERE reception_id = $1 FOR UPDATE", receptionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.Id, &product.DateTime, &product.Type, &product.ReceptionId)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return products, nil
}
