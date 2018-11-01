package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func checkCommandLineArgs() {

	//networked := flag.Bool("networked", false, "get remote config")
	//netname = flag.String("netname", "chars", "Metric {chars|words|lines};.")
	debugset := flag.Bool("debug", false, "turn the noise on")
	//configuration["boardsToWatch"] = *flag.String("boards", "DevOps2020 - Board", "board 1, board 2, board n")
	boards := flag.String("boards", "DevOps2020 - Board", "board 1, board 2, board n")
	label := flag.String("label", "", "Label to look for")
	config := flag.String("configfile", "", "Path to configuration file")
	token := flag.String("tokenfile", "", "Path to API token and key file")
	flag.Parse()
	if *debugset {
		log.SetLevel(log.DebugLevel)
	}
	if *label != "" {
		configuration["toPrintedLabelName"] = *label
	}
	if *boards != "" {
		configuration["boardsToWatch"] = *boards
	}
	if *config != "" {
		configFile = *config
	}
	if *token != "" {
		tokenFile = *token
	}
	// TODO the debugger does this wrong
	// I need this for debugging...
	//tokenFile = ".token"
	//configFile = "config.cfg"

	return
}
func fetchIP() string {
	localIPAddr := getOutboundIP()
	log.Debugf("%v", localIPAddr)
	return localIPAddr.String()

}

func fetchConfiguration() {
	readConfigFromFile(configFile)
	readConfigFromFile(tokenFile)
	pdfDocDimension = getPdfDocDimensionFromString()
	pdfMargins = getPdfMarginsFromString()
	qRCodePos = getqRCodePosFromString()
	blackRectPos = getBlackRectPosFromString()
	fetchBoardListFromConfig()

}

func fetchBoardListFromConfig() {
	// try this.
	boardsToWatch = strings.Split(configuration["boardsToWatch"], ",")
	log.Debugf("board list: %v", boardsToWatch)
}

func readConfigFromFile(filename string) {
	if filename == "" {
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())

		a := strings.Split(string(scanner.Text()), "=")
		_, ok := configuration[a[0]]
		if a[0] != "" && ok {
			configuration[strings.Trim(a[0], " ")] = strings.Trim(a[1], " ")
		}

	}

}
func init() {

	checkCommandLineArgs()
	fetchConfiguration()
	configuration["ip"] = fetchIP()
	fetchBoardListFromConfig()
	log.Infof("IP is %v", configuration["ip"])
	dir, err := ioutil.TempDir(os.TempDir(), configuration["tmpDirPrefix"])
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf(dir)
	configuration["tmpDirName"] = dir
	id, err := machineid.ProtectedID("trelloknect")
	if err != nil {
		log.Fatal(err)
	}
	configuration["knechtID"] = id
}

func cleanUp(dirName string) {
	os.RemoveAll(dirName)

}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func getUUID() string {
	uuid := uuid.New()
	return uuid.String()

}

func configCardDescription() string {
	text := "Trelloprinter: " + configuration["printerName"] + "\n"
	text = text + "IP: " + configuration["ip"] + "\n"
	text = text + "Watched Boards: " + configuration["boardsToWatch"] + "\n"
	text = text + "Label to Print: " + configuration["toPrintedLabelName"] + "\n"

	return text
}
func (r *Resultset) execCommand() {
	cmd := exec.Command(r.OSCommand, r.CommandArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	r.SuccessfullExecution = true
	r.CmdStarttime = time.Now()

	err := cmd.Run()
	r.Stdout = string(stdout.Bytes())
	r.Stderr = string(stderr.Bytes())
	r.CMDStoptime = time.Now()
	r.DurationSecounds = int(r.CMDStoptime.Unix() - r.CmdStarttime.Unix())
	if err != nil {
		//todo: log.Fatalf("cmd.Run() failed with %s\n", err)
		log.Errorf("Command failed %v err: \n", err)
		r.SuccessfullExecution = false
		r.ErrorStr = err.Error()

	}
	return
}

func reportPrints() {
	for _, pdf := range printedCards {
		log.Infof("printed card %v", cardByFileName[pdf].Name)
	}

}
func sweepOut() {
	cardByFileName = nil
	printedCards = nil

}
