package logging

import (
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	appName = "urn:entrust:obac-server"
	version = "1.0"
)

const (
	ServiceStarted            string = "Start"
	EnforceHandled            string = "Enforce"
	GetTokenIdsHandled        string = "GetAssets"
	AuthorizeHandled          string = "Authorize"
	GetTokenHandled           string = "GetToken"
	GetPolicyContractsHandled string = "GetPolicyContracts"
	CancelHandled             string = "Cancel"
	PostMappingHandled        string = "PostMapping"
	GetMappingHandled         string = "GetMapping"
	GetAssociationsHandled    string = "GetAssociations"
	EditHandled               string = "PolicyEdit"
	GetPolictHandled          string = "GetPolicy"
)

func InitLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func InfoEvent(event string, fields log.Fields) {
	entry := log.WithFields(getStaticFields(event))
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.Info(event + " operation finished")
}

func InfoErrorEvent(event string, fields log.Fields, err error) {

	entry := log.WithFields(getStaticFields(event))
	entry = entry.WithFields(log.Fields{
		"error": err,
	})
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.Info(event + " operation finished")
}

func ErrorEvent(event string, fields log.Fields, err error) {
	entry := log.WithFields(getStaticFields(event))
	if fields != nil {
		entry = entry.WithFields(fields)
	}
	entry.WithError(err).Error(event + " finished with error")
}

func Info(fields log.Fields) {
	InfoEvent(GetFunctionName(), fields)
}

func GetFunctionName() string {
	pc := make([]uintptr, 10)
	runtime.Callers(3, pc)
	funcName := runtime.FuncForPC(pc[0]).Name()

	split := strings.Split(funcName, ".")
	return split[len(split)-1]
}

func getStaticFields(event string) log.Fields {
	return log.Fields{
		"appName":  appName,
		"version":  version,
		"event":    event,
		"hostname": getHostname(),
	}
}

func getHostname() string {
	name, _ := os.Hostname()
	return name
}
