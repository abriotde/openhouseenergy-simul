/*
The goal of this file is to use OpenHouseEnergyModule as CLI (for client and server)
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/abriotde/openhouseenergy-simul/logger"
	"github.com/abriotde/openhouseenergy-simul/module"
	"github.com/abriotde/openhouseenergy-simul/server"
	"github.com/spf13/cobra"
)

const (
	EXIT_ARGUMENT_ERROR = 1
)

// LLDP implementation : https://pkg.go.dev/github.com/ksang/gotopo/lldp

func runClientCmd(client server.OpenHouseEnergyModuleClient, args []string) error {
	var argsLen = len(args)
	if argsLen < 1 {
		logger.Logger.Error("You must give comands to client.")
		os.Exit(EXIT_ARGUMENT_ERROR)
	} else if argsLen == 1 {
		logger.Logger.Error("Missing arguments.")
		os.Exit(EXIT_ARGUMENT_ERROR)
	}
	var clientCmd = args[0]
	if clientCmd == "send" {
		if argsLen < 3 {
			logger.Logger.Error("Missing arguments.")
			os.Exit(EXIT_ARGUMENT_ERROR)
		}
		var varName = args[1]
		varValue, err := strconv.Atoi(args[2])
		if err != nil {
			logger.Logger.Error("Argument 2 must be an integer for the value : " + args[2] + ".")
			return err
		}
		// TODO : check varname/varvalue match possible value (No injection)
		fmt.Println("Send to server : ", varName, " = ", varValue)
		_, err = client.SendModuleDescription()
		if err != nil {
			return err
		}
	} else if clientCmd == "get" {
		/*	if argsLen > 1 && args[1] == "alerts" {
				// fmt.Println("Call server GetAlertHistory.")
				alerts, err := client.GetAlertHistory()
				if err != nil {
					return err
				}
				fmt.Println("Alerts :")
				for _, alert := range alerts {
					fmt.Println(" -> ", alert.GetName(), " for value = ", alert.GetValue())
				}
			} else {
				logger.Logger.Error("Unimplemented parameter for get.")
			}*/
	} else if clientCmd == "help" {
		fmt.Println("Existing commands are : \n - 'send [type] [int_value]' : Implemented types are cpu (Should be less than 80) and battery (Should be beetween 20 and 98) : see monitorer.go for more informations. \n - 'get alerts'\n - help \n - quit : to exit on interactive mode.\n")
	} else {
		logger.Logger.Error("Unknown client command : '" + clientCmd + "' possibilities are send|get|help|quit.")
		return errors.New("Unknown client command : '" + clientCmd + "' possibilities are send|get|help|quit.")
	}
	return nil
}

var (
	// Used for flags.
	interactive   bool
	coreModuleURL string
	port          string
	confFile      string
	rootCmd       = &cobra.Command{
		Use:   "openhouseenergy-simul",
		Short: "openhouseenergy-simul is a simulation for openhouseenrgy manager.",
		Long: `openhouseenergy-simul is a simulation for openhouseenrgy manager.
			This is licenced under GPL V3`,
		Run: func(cmd *cobra.Command, args []string) {
			var myModule := module.OpenHouseEnergyModule_Init(confFile)
			if coreModuleURL != "" { // Run as simple module
				file, err := os.OpenFile("log/client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					logger.Logger.Out = file
				}
				client, err := myModule.Connect(coreModuleURL)
				if err != nil {
					logger.Logger.Error("Fail to connect to server : " + coreModuleURL + ".")
					os.Exit(EXIT_ARGUMENT_ERROR)
				}
				_, err = client.SendModuleDescription()
				if err != nil {
					logger.Logger.Error("Fail to signal the module to the server")
					os.Exit(EXIT_ARGUMENT_ERROR)
				}
				if interactive {
					fmt.Println("Interactive mode is enable.")
					reader := bufio.NewReader(os.Stdin)
					run := true
					for run {
						fmt.Print(" $ ")
						str, err := reader.ReadString('\n')
						if err != nil {
							logger.Logger.Error("Fail read input.")
							continue
						}
						last := len(str) - 1
						str = str[:last] // Remove last character : \n
						args := strings.Split(str, " ")
						if len(args) > 0 && args[0] == "quit" {
							run = false
							fmt.Println("Goodbye.")
							break
						}
						runClientCmd(client, args)
					}
					// TODO : implement
				} else {
				}
				client.Close()
			} else { // Run as core module mode

			}
			if port != "" {
				file, err := os.OpenFile("log/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					logger.Logger.Out = file
				}
				// var server OpenHouseEnergyModuleServer
				fmt.Println("Run as server mode on port ", port, ".")
				_, err = myModule.Listen(port)
				if err != nil {
					logger.Logger.Error(" : .") // TODO: + err
					os.Exit(EXIT_ARGUMENT_ERROR)
				}
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Logger.Error(err)
		os.Exit(EXIT_ARGUMENT_ERROR)
	}
}

// --interactive -i
// --server -s --client -c
// send battery|cpu XX
// get alerts
func init() {
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "It run on interactive mode.")
	rootCmd.PersistentFlags().StringVarP(&coreModuleURL, "core", "c", "", "The core module url, so it's not a core module.")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "The port to connect to listen instructions.")
	rootCmd.PersistentFlags().StringVarP(&confFile, "module", "m", "", "The file that describe the module.")

}
