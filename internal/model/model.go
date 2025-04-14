package model

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleEmployee, RoleModerator:
		return true
	}
	return false
}

type City string

const (
	CityMoscow          City = "Москва"
	CityKazan           City = "Казань"
	CitySaintPetersburg City = "Санкт-Петербург"
)

func (c City) IsValid() bool {
	switch c {
	case CityMoscow, CityKazan, CitySaintPetersburg:
		return true
	}
	return false
}

type ReceptionStatus string

const (
	ReceptionStatusInProgress ReceptionStatus = "in_progress"
	ReceptionStatusClose      ReceptionStatus = "close"
)

type ProductType string

const (
	ProductTypeElectronics ProductType = "электроника"
	ProductTypeClothes     ProductType = "одежда"
	ProductTypeShoes       ProductType = "обувь"
)

func (p ProductType) IsValid() bool {
	switch p {
	case ProductTypeClothes, ProductTypeElectronics, ProductTypeShoes:
		return true
	}
	return false
}

type DummyLoginRequest struct {
	Role Role `json:"role"`
}

type CreatePvzRequest struct {
	City City `json:"city"`
}

type Pvz struct {
	Id               uuid.UUID `json:"id" db:"id"`
	RegistrationDate time.Time `json:"registration_date" db:"registration_date"`
	City             City      `json:"city" db:"city"`
}

type Reception struct {
	Id       uuid.UUID       `json:"id" db:"id"`
	DateTime time.Time       `json:"date_time" db:"date_time"`
	PvzId    uuid.UUID       `json:"pvz_id" db:"pvz_id"`
	Status   ReceptionStatus `json:"status" db:"status"`
}

type CreateReceptionRequest struct {
	PvzId string `json:"pvz_id"`
}

type AddProductRequest struct {
	Type  ProductType `json:"type"`
	PvzId string      `json:"pvz_id"`
}

type Product struct {
	Id          uuid.UUID   `json:"id" db:"id"`
	DateTime    time.Time   `json:"date_time" db:"date_time"`
	Type        ProductType `json:"type" db:"type"`
	ReceptionId string      `json:"reception_id" db:"reception_id"`
}

type ReceptionInfo struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products"`
}

type PvzInfo struct {
	Pvz        Pvz             `json:"pvz"`
	Receptions []ReceptionInfo `json:"receptions"`
}

type GetPvzInfoResponse struct {
	PvzList []PvzInfo `json:"pvz_list"`
}
