package container

import (
	"5g-v2x-data-management-service/internal/config"
	controller "5g-v2x-data-management-service/internal/controllers"
	"5g-v2x-data-management-service/internal/infrastructures/database"
	"5g-v2x-data-management-service/internal/infrastructures/http"
	"5g-v2x-data-management-service/internal/repositories"
	"5g-v2x-data-management-service/internal/services"

	"go.uber.org/dig"
)

type Container struct {
	container *dig.Container
	Error     error
}

func NewContainer() *Container {
	c := new(Container)
	c.Configure()
	return c
}

func (cn *Container) Configure() {
	cn.container = dig.New()

	// infrastructures
	if err := cn.container.Provide(http.NewGRPCServer); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(database.NewMongoDatabase); err != nil {
		cn.Error = err
	}

	// config
	if err := cn.container.Provide(config.NewConfig); err != nil {
		cn.Error = err
	}

	// controllers
	if err := cn.container.Provide(controller.NewController); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(controller.NewAccidentController); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(controller.NewDrowsinessController); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(controller.NewCarController); err != nil {
		cn.Error = err
	}

	// services
	if err := cn.container.Provide(services.NewAccidentService); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(services.NewDrowsinessService); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(services.NewCarService); err != nil {
		cn.Error = err
	}

	// repositories
	if err := cn.container.Provide(repositories.NewCRUDRepository); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(repositories.NewAccidentRepository); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(repositories.NewDrowsinessRepository); err != nil {
		cn.Error = err
	}

	if err := cn.container.Provide(repositories.NewCarRepository); err != nil {
		cn.Error = err
	}

}

func (cn *Container) Run() *Container {
	if err := cn.container.Invoke(func(g *http.GRPCServer) {
		if err := g.Start(); err != nil {
			panic(err)
		}
	}); err != nil {
		panic(err)
	}
	return cn
}
