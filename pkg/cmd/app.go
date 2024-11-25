package cmd

import (
	"fmt"

	"github.com/mchmarny/airthings-go/pkg/client"
)

func Execute() {
	c := client.NewClient()
	r, err := c.GetDevices()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Devices: %v\n", r)
}
