package handler_manager

import (
	"avito2/internal/service"
	"avito2/internal/utils"
)

type HandlerManager struct {
	jwtGen utils.JWTGenerator
	svc    service.Service
}

func NewHandlerManager(svc service.Service, jwtGen utils.JWTGenerator) *HandlerManager {
	return &HandlerManager{
		svc:    svc,
		jwtGen: jwtGen,
	}
}
