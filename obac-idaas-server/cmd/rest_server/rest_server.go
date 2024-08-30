package main

import (
	"entrust.com/iat/golang-server/pkg/server"
	"entrust.com/iat/obac-idaas-server/pkg/config"
	"entrust.com/iat/obac-idaas-server/pkg/controller"
	"entrust.com/iat/obac-idaas-server/pkg/logging"
	"entrust.com/iat/obac-idaas-server/pkg/repository"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	logging.InitLogger()

	mainConfig, err := config.LoadConfig("config")
	if err != nil {
		log.Fatal(err)
	}

	//authKeys := oidc2.NewAuthKeys()
	//err = authKeys.Load(mainConfig.AuthConfig)
	//if err != nil {
	//	log.Fatal(err)
	//}

	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	bcRepo, err := repository.NewBlockchainRepository(&mainConfig.BcProvider)
	if err != nil {
		logging.InfoErrorEvent("Blockchain repo initialization", logrus.Fields{"config": mainConfig.BcProvider}, err)
		log.Fatal(err)
	}

	cadRepo, err := repository.NewContentAddressableRepository()
	if err != nil {
		log.Fatal(err)
	}

	controllerManager := controller.NewControllerManager(mainConfig.App, mainConfig.Idp, repo, bcRepo, cadRepo)
	routes := make(map[string]bool)

	routes[controller.AuthorizePath] = false
	routes[controller.TokenPath] = false
	routes[controller.EnforcePath] = false
	routes[controller.GetAssetsPath] = false
	routes[controller.BookingInfoPath] = false
	routes[controller.EditPolicyPath] = false
	routes2 := []string{controller.AuthorizePath, controller.TokenPath, controller.EnforcePath, controller.GetAssetsPath, controller.BookingInfoPath, controller.EditPolicyPath, controller.GetPolicyContractsPath, controller.GetPolicyPath, controller.GetMappingPath, controller.PostMappingPath, controller.GetAssociationsPath}

	restServer := server.NewPrivateServer(mainConfig.ServerConfig, routes2)

	err = restServer.RegisterEndpoint(controller.AuthorizePath, "GET", "", controllerManager.OIDCHandlerController.Authorize)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.TokenPath, "POST", "", controllerManager.OIDCHandlerController.GetToken)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.EnforcePath, "POST", "", controllerManager.OIDCHandlerController.Enforce)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.GetAssetsPath, "GET", "", controllerManager.OIDCHandlerController.GetAssets)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.GetPolicyContractsPath, "GET", "", controllerManager.PolicyHandlerController.GetPolicyContracts)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.GetMappingPath, "GET", "", controllerManager.OIDCHandlerController.GetMapping)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.PostMappingPath, "POST", "", controllerManager.OIDCHandlerController.PostMapping)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.GetAssociationsPath, "GET", "", controllerManager.OIDCHandlerController.GetAssociations)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.BookingInfoPath, "GET", "", controllerManager.FIFAHandlerController.GetBookingInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.BookingInfoPath, "POST", "", controllerManager.FIFAHandlerController.InsertBookingInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.EditPolicyPath, "POST", "", controllerManager.PolicyHandlerController.Edit)
	if err != nil {
		log.Fatal(err)
	}

	err = restServer.RegisterEndpoint(controller.GetPolicyPath, "GET", "", controllerManager.PolicyHandlerController.GetPolicy)
	if err != nil {
		log.Fatal(err)
	}

	logFields := logrus.Fields{
		"serverPort": mainConfig.ServerConfig.Port,
		"withTls":    mainConfig.ServerConfig.TLS.CertPath != "",
	}

	logging.InfoEvent(logging.ServiceStarted, logFields)
	err = restServer.Serve()
	if err != nil {
		logging.ErrorEvent(logging.ServiceStarted, logFields, err)
	}
}
