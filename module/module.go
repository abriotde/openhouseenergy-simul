package module

import (
	"fmt"
	"log"
	"os"

	"github.com/abriotde/openhouseenergy-simul/messages"
	"github.com/abriotde/openhouseenergy-simul/server"
	"gopkg.in/yaml.v3"
)

type OpenHouseEnergyModule interface {
	GetModuleDescription() *messages.SendModuleDescriptionRequest
	Listen(port int32) (server.OpenHouseEnergyModuleServer, error)
	Connect(url string) (server.OpenHouseEnergyModuleClient, error)
}

type moduleConfiguration struct {
	Id           int32
	ModuleType   string `yaml:"type"`
	ModuleTypeId int32  `yaml:"typeId"`
	MaxValue     int32  `yaml:"max"`
}

func (c *moduleConfiguration) getConfiguration(confFile string) *moduleConfiguration {
	yamlFile, err := os.ReadFile(confFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func New(confFile string) OpenHouseEnergyModule {
	var conf moduleConfiguration
	conf.getConfiguration(confFile)
	fmt.Println("Load conf : ", conf)
	if conf.ModuleType == "solarPannel" {
		pannel := OpenHouseEnergyModuleSolarPannel{
			id:       conf.Id,
			maxPower: conf.MaxValue}
		return pannel
	}
	return nil
}
