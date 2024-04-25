package server

import (
	"context"
	"net"
	"strconv"

	"github.com/abriotde/openhouseenergy-simul/logger"
	"github.com/abriotde/openhouseenergy-simul/messages"
	"github.com/abriotde/openhouseenergy-simul/monitorer"
	"google.golang.org/grpc"
)

type OpenHouseEnergyModuleServer struct {
	listener   net.Listener
	connection net.Conn
	connected  bool
	modules    []int32
}

type server_t struct {
	messages.UnimplementedGreeterServer
	monitoring monitorer.Monitorer
}

var monitoring_server server_t

// Launch the server and never stop (Hard kill for stop : TODO : better way: special message?)
func Listen(port int32) (OpenHouseEnergyModuleServer, error) {
	var server = OpenHouseEnergyModuleServer{connected: false}
	var portStr = strconv.Itoa(int(port))
	listener, err := net.Listen("tcp", "localhost:"+portStr)
	if err != nil {
		logger.Logger.Error("OpenHouseEnergyModuleServer : Impossible listen on : " + portStr + ".")
		return server, err
	}
	defer listener.Close()
	server.listener = listener
	server.connected = true
	server.Run()
	return server, nil
}

// To treat the SendModuleDescription request. It will register value of variable at current time but for the moment it register only if it alerts.
func (s *server_t) SendModuleDescription(ctx context.Context, in *messages.SendModuleDescriptionRequest) (*messages.SendModuleDescriptionReply, error) {
	// sValue := strconv.Itoa(int(description.Id))
	logger.Logger.Info("Received: ", in.GetModuleType(), " from ", in.Id)
	// s.monitoring.Log(description.GetType(), sValue)
	return &messages.SendModuleDescriptionReply{Ok: true, ModuleId: in.Id}, nil
}

// To treat the GetAlertHistory request.
/* func (s *server_t) GetAlertHistory(ctx context.Context, in *messages.GetAlertHistoryRequest) (*messages.GetAlertHistoryReply, error) {
	logger.Logger.Info("Ask for alerts.")
	var alerts []*messages.GetAlertHistoryReply_Alert
	var nbAlerts = 0
	// On converti les alertes "Monitorer" en alertes "protobuf"
	for _, alert := range s.monitoring.GetAlertHistory() {
		a := messages.GetAlertHistoryReply_Alert{Timestamp: timestamppb.New(alert.Timestamp), Name: alert.Name, Value: alert.Value}
		alerts = append(alerts, &a)
		nbAlerts += 1
	}
	logger.Logger.Info("Have ", strconv.Itoa(nbAlerts), " alerts.")
	return &messages.GetAlertHistoryReply{AlertHistory: alerts, Ok: true}, nil
} */
func init() {
	// fmt.Println("Init().")
}

func (server OpenHouseEnergyModuleServer) Run() (OpenHouseEnergyModuleServer, error) {
	monitoring_server.monitoring.Logger = logger.Logger
	grpcServer := grpc.NewServer()
	messages.RegisterGreeterServer(grpcServer, &monitoring_server)
	logger.Logger.Info("Server listening at " + server.listener.Addr().String())
	if err := grpcServer.Serve(server.listener); err != nil {
		logger.Logger.Error("failed to serve: ") // TODO: + err
		return server, err
	}
	return server, nil
}

// Function to test client/server communication with simple ASCII (useless now?)
/* func (server OpenHouseEnergyModule) Test() (OpenHouseEnergyModule, error) {
	for {
		// Waiting connection
		conn, err := server.listener.Accept()
		if err != nil {
			logger.Logger.Error("Impossible accept.")
			return server, err
		}
		server.connection = conn

		// Have a connection, read request
		request, err := bufio.NewReader(server.connection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return server, err
		}
		if strings.TrimSpace(string(request)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return server, nil
		}
		fmt.Print("-> ", string(request))

		// Send response
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		server.connection.Write([]byte(myTime))
	}
	return server, nil
} */
