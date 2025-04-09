package handler

import (
	"gofr.dev/pkg/gofr"
)

type User interface {
	Create(ctx *gofr.Context) (interface{}, error)
	GetByID(ctx *gofr.Context) (interface{}, error)
	GetAll(ctx *gofr.Context) (interface{}, error)
	Update(ctx *gofr.Context) (interface{}, error)
	Delete(ctx *gofr.Context) (interface{}, error)
}

type Account interface {
	Create(ctx *gofr.Context) (interface{}, error)
	GetByID(ctx *gofr.Context) (interface{}, error)
	GetAll(ctx *gofr.Context) (interface{}, error)
	Update(ctx *gofr.Context) (interface{}, error)
	Delete(ctx *gofr.Context) (interface{}, error)
}

type Transactions interface {
	Create(ctx *gofr.Context) (interface{}, error)
	GetAll(ctx *gofr.Context) (interface{}, error)
	GetByID(ctx *gofr.Context) (interface{}, error)
	Update(ctx *gofr.Context) (interface{}, error)
	Delete(ctx *gofr.Context) (interface{}, error)
}

type Savings interface {
	Create(ctx *gofr.Context) (interface{}, error)
	GetAll(ctx *gofr.Context) (interface{}, error)
	GetByID(ctx *gofr.Context) (interface{}, error)
	Update(ctx *gofr.Context) (interface{}, error)
	Delete(ctx *gofr.Context) (interface{}, error)
}

type Dashboard interface {
	Get(ctx *gofr.Context) (interface{}, error)
}

type Auth interface {
	CreateToken(ctx *gofr.Context) (interface{}, error)
	Login(ctx *gofr.Context) (interface{}, error)
	Refresh(ctx *gofr.Context) (interface{}, error)
}
