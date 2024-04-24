/*
	The goal of this file is to use MiniserverAisprid as CLI (for client and server)
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/abriotde/minialertAisprid/server"
	"github.com/abriotde/minialertAisprid/logger"
	"bufio"
	"strings"
	"errors"
)

const (
	EXIT_ARGUMENT_ERROR = 1
)

func runClientCmd(client server.MiniserverAispridClient, args []string) error {
	var argsLen = len(args)
	if argsLen<1 {
		logger.Logger.Error("You must give comands to client.")
		os.Exit(EXIT_ARGUMENT_ERROR)
	} else if argsLen<1 {
		logger.Logger.Error("Missing arguments.")
		os.Exit(EXIT_ARGUMENT_ERROR)
	}
	var clientCmd = args[0]
	if clientCmd=="send" {
		if argsLen<3 {
			logger.Logger.Error("Missing arguments.")
			os.Exit(EXIT_ARGUMENT_ERROR)
		}
		var varName = args[1]
		varValue,err := strconv.Atoi(args[2])
		if err != nil {
			logger.Logger.Error("Argument 2 must be an integer for the value : "+args[2]+".")
			return err
		}
		// TODO : check varname/varvalue match possible value (No injection)
        	fmt.Println("Send to server : ", varName, " = ", varValue)
        	_,err = client.Set(varName, int32(varValue))
        	if err!=nil {
        		return err
        	}
	} else if clientCmd=="get" {
		if argsLen>1 && args[1]=="alerts" {
			// fmt.Println("Call server GetAlertHistory.")
			alerts,err := client.GetAlertHistory()
			if err!=nil {
				return err
			}
			fmt.Println("Alerts :")
			for _, alert := range alerts {
				fmt.Println(" -> ", alert.GetName(), " for value = ", alert.GetValue())
			}
		} else {
			logger.Logger.Error("Unimplemented parameter for get.")
		}
	} else if clientCmd=="help" {
		fmt.Println("Existing commands are : \n - 'send [type] [int_value]' : Implemented types are cpu (Should be less than 80) and battery (Should be beetween 20 and 98) : see monitorer.go for more informations. \n - 'get alerts'\n - help \n - quit : to exit on interactive mode.\n")
	} else {
		logger.Logger.Error("Unknown client command : '"+clientCmd+"' possibilities are send|get|help|quit.")
		return errors.New("Unknown client command : '"+clientCmd+"' possibilities are send|get|help|quit.")
	}
        return nil
}

var (
	// Used for flags.
	interactive     bool
	serverURL		string
	port		string
	rootCmd = &cobra.Command{
		Use:   "minialertAisprid",
		Short: "MinialertAisprid is a minimalistic chalenge to send messages and receive alerts.",
		Long: `Minialert is a minimalistic chalenge consisting in a client/server which send messages and receive alerts lists
			This is licenced under GPL V3`,
		Run: func(cmd *cobra.Command, args []string) {
			if port!="" {
				file, err := os.OpenFile("log/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					logger.Logger.Out = file
				}
				// var server MiniserverAisprid
				fmt.Println("Run as server mode on port ", port, ".")
				_, err = server.Listen("localhost:"+port)
				if err != nil {
					logger.Logger.Error(" : .") // TODO: + err
					os.Exit(EXIT_ARGUMENT_ERROR)
				}
			} else {
				file, err := os.OpenFile("log/client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					logger.Logger.Out = file
				}
				client, err := server.Connect(serverURL)
				if err != nil {
					logger.Logger.Error("Fail to connect to server : "+serverURL+".")
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
						if len(args)>0 && args[0]=="quit" {
							run = false
			    				fmt.Println("Goodbye.")
							break
						}
		    				runClientCmd(client, args)
			    		}
			    		// TODO : implement
			    	} else {
			    		runClientCmd(client, args)
			   	}
			   	client.Close()
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
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "For client, it run on interactive mode.")
	rootCmd.PersistentFlags().StringVarP(&serverURL, "server", "s", "localhost:8080", "The server to connect when in client mode (default). If no port specified, it connect on 8080 port." )
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "The port to connect so run it as server.")
}

