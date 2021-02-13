package main

import (
	"5g-v2x-data-management-service/internal/container"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	if err := container.NewContainer().Run().Error; err != nil {
		panic(err)
	}

}
