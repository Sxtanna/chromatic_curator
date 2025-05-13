package app

import "github.com/Sxtanna/chromatic_curator/internal/common"

type Service interface {
	// Start should either block or return common.ServiceStartedNormallyButDoesNotBlock
	Start() error

	Close(_ error) error
}

type InitializedService interface {
	Init(config common.Configuration) error

	Service
}
