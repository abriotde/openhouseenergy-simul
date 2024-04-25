package module

import (
	"strconv"

	"github.com/abriotde/openhouseenergy-simul/logger"
	"github.com/abriotde/openhouseenergy-simul/messages"
	"github.com/abriotde/openhouseenergy-simul/server"
)

type OpenHouseEnergyModuleSolarPannel struct {
	id       int32
	maxPower int32
}

// Launch the server and never stop (Hard kill for stop : TODO : better way: special message?)
func (module OpenHouseEnergyModuleSolarPannel) Listen(port int32) (server.OpenHouseEnergyModuleServer, error) {
	myServer, err := server.Listen(port)
	if err != nil {
		logger.Logger.Error("OpenHouseEnergyModuleSolarPannel : Impossible listen on : " + strconv.Itoa(int(port)) + ".")
		return myServer, err
	}
	return myServer, nil
}

func (module OpenHouseEnergyModuleSolarPannel) Connect(url string) (server.OpenHouseEnergyModuleClient, error) {
	myClient, err := server.Connect(url)
	if err != nil {
		logger.Logger.Error("OpenHouseEnergyModuleSolarPannel : Impossible connect to : " + url + ".")
		return myClient, err
	}
	return myClient, nil
}

func (pannel OpenHouseEnergyModuleSolarPannel) GetModuleDescription() *messages.SendModuleDescriptionRequest {
	description := messages.SendModuleDescriptionRequest{
		Id:         pannel.id,
		ModuleType: messages.SendModuleDescriptionRequest_SOLARPANEL}
	return &description
}
