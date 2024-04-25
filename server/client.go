package server

import (
	"context"
	"time"

	"github.com/abriotde/openhouseenergy-simul/logger"
	"github.com/abriotde/openhouseenergy-simul/messages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// It is client for connecting to "OpenHouseEnergyModule" server.
type OpenHouseEnergyModuleClient struct {
	connection     *grpc.ClientConn
	grpcConnection messages.GreeterClient
	connected      bool
}

// Fonction to connect to the server : prerequisites for all communications : call to GetAlertHistory and Set
func Connect(url string) (OpenHouseEnergyModuleClient, error) {
	var client = OpenHouseEnergyModuleClient{connected: false}
	// conn, err := net.Dial("tcp", url)
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Logger.Error("Impossible to connect to : " + url + ".")
		return client, err
	}
	client.connection = conn
	client.connected = true
	client.grpcConnection = messages.NewGreeterClient(client.connection)
	logger.Logger.Debug("Connected to  : ", url)
	return client, nil
}

// Fonction to close connection to the server
func (client OpenHouseEnergyModuleClient) Close() {
	client.connection.Close()
}

// Return alert history based on previous Set() call where name/value exceed threashold defined on server side.
/* func (client OpenHouseEnergyModuleClient) GetCoreQuery() ([]*messages.CoreQueryRequest, error) {
	// fmt.Println("GetAlertHistory from server : ")
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.grpcConnection.GetCoreQuery(ctx, &messages.CoreQueryRequest{})
	if err != nil {
		logger.Logger.Error("could not GetAlertHistory") // TODO: + err
		return make([]*messages.CoreQueryRequest, 0), err
	}
	if r.GetOk() != true {
		logger.Logger.Error("could not GetAlertHistory : Server refuse.")
		return make([]*messages.GetAlertHistoryReply_Alert, 0), errors.New("could not GetAlertHistory : Server refuse.")
	}
	return r.GetAlertHistory(), nil
} */

// Set the current value of the variable name registred on the server.
func (client OpenHouseEnergyModuleClient) SendModuleDescription() (string, error) {
	// fmt.Println("Set to server : ", varName, " = ", varValue)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.grpcConnection.SendModuleDescription(ctx, &messages.SendModuleDescriptionRequest{})
	if err != nil {
		logger.Logger.Error("could not SendModuleDescription : : ") // TODO: + err
		return "KO", err
	}
	if r.GetOk() != true {
		logger.Logger.Error("could not SendModuleDescription : Server refuse.")
		return "KO", nil
	}
	return "OK", nil
}

// Function to test client/server communication with simple ASCII (useless now?)
func (client OpenHouseEnergyModuleClient) Test() (string, error) {
	/* fmt.Println("Test mode : client")
	        reader := bufio.NewReader(os.Stdin)
	        fmt.Print(">> ")
	        text, _ := reader.ReadString('\n')
	        fmt.Fprintf(client.connection, text+"\n")

	        message, _ := bufio.NewReader(client.connection).ReadString('\n')
	        fmt.Print("->: " + message)
	        if strings.TrimSpace(string(text)) == "STOP" {
	                fmt.Println("TCP client exiting...")
			return "OK", nil
	        } */
	return "OK", nil
}
