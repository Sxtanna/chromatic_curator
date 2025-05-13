package app

import (
	"emperror.dev/errors"
	emp "emperror.dev/handler/logur"
	"github.com/Sxtanna/chromatic_curator/internal/app/discord"
	"github.com/Sxtanna/chromatic_curator/internal/common"
	log "logur.dev/logur"
)

var (
	services = make([]Service, 0)
)

func InitializeApp(abort <-chan struct{}, logger log.Logger, handler *emp.Handler, config common.Configuration) common.Group {
	group := make(common.Group, 0)

	addServiceToGroup := func(service Service) {
		group = group.Act(
			func() error {
				err := service.Start()

				if err == nil || !errors.Is(err, common.ServiceStartedNormallyButDoesNotBlock) {
					return err
				}

				group.Await(abort)

				return nil
			},
			func(err error) {
				handler.Handle(service.Close(err))
			},
		)
	}

	services = append(services, &discord.BotService{Logger: logger})

	for _, service := range services {
		if inits, ok := service.(InitializedService); ok {
			if err := inits.Init(config); err == nil {
				addServiceToGroup(service)
			} else {
				handler.Handle(inits.Init(config))
			}
		}
	}

	return group
}
