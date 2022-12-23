package main

import (
	"flag"
	"fmt"
	zd "github.com/easysoft/zendata"
	"github.com/easysoft/zendata/internal/agent"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/server"
	serverConfig "github.com/easysoft/zendata/internal/server/config"
	"github.com/easysoft/zendata/internal/server/core/web"
	serverConst "github.com/easysoft/zendata/internal/server/utils/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	flagSet *flag.FlagSet
	uuid    = ""
)

func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	flagSet = flag.NewFlagSet("zd", flag.ContinueOnError)

	flagSet.StringVar(&uuid, "uuid", "", "区分服务进程的唯一ID")

	flagSet.StringVar(&vari.Ip, "i", "", "")
	flagSet.StringVar(&vari.Ip, "ip", "", "")
	flagSet.IntVar(&vari.DataServicePort, "p", 0, "")
	flagSet.IntVar(&vari.DataServicePort, "port", 0, "")
	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	configUtils.InitConfig("")
	vari.DB, _ = serverConfig.NewGormDB()

	vari.AgentLogDir = vari.ZdPath + serverConst.AgentLogDir + consts.PthSep
	err := fileUtils.MkDirIfNeeded(vari.AgentLogDir)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("perm_deny", vari.AgentLogDir), color.FgRed)
		os.Exit(1)
	}

	if vari.Ip == "" {
		vari.Ip = commonUtils.GetIp()
	}
	if vari.DataServicePort == 0 {
		vari.DataServicePort = consts.DefaultDataServicePort
	}

	go func() {
		startDataServer()
	}()

	startAdminServer()
}

func startDataServer() {
	port := strconv.Itoa(vari.DataServicePort)
	logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server",
		vari.Ip, port, vari.Ip, port, vari.Ip, port), color.FgCyan)

	config := serverConfig.NewConfig()
	server, err := server.InitServer(config)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server_fail", port), color.FgRed)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", server.Config.ServerPort),
		Handler: dataHandler(server),
	}

	httpServer.ListenAndServe()
}

func startAdminServer() {
	webServer := web.Init()
	if webServer == nil {
		return
	}

	webServer.Run()
}

func dataHandler(server *server.Server) http.Handler {
	mux := http.NewServeMux()

	uiFs, err := zd.GetUiFileSys()
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(uiFs)))

	//mux.HandleFunc("/admin", server.AdminHandler)
	mux.HandleFunc("/data", agent.DataHandler)

	return mux
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
