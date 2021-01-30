package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/actapio/moxspec/loglet"
	"github.com/actapio/moxspec/model"
	"github.com/actapio/moxspec/pci"
	"github.com/actapio/moxspec/smbios"
)

var (
	log      *loglet.Logger
	version  string
	revision string
)

func init() {
	log = loglet.NewLogger("main")
}

func main() {
	var err error

	runAs := filepath.Base(os.Args[0])

	switch runAs {
	case "lsraid":
		rootOrExit()
		err := lsraid()
		if err != nil {
			log.Fatal(err)
		}
		return
	case "lsdiag":
		rootOrExit()
		exitCode, err := lsdiag()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(exitCode)
	case "lssn":
		rootOrExit()
		err := lssn()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	cli := newApp(os.Args)
	if cli == nil {
		showHelp()
		os.Exit(1)
	}

	cli.appendFlag("d", false, "show verbose log")
	cli.appendFlag("noraidcli", false, "disable running RAID utilities")
	switch cli.cmd {
	case "show":
		cli.appendFlag("j", false, "print json")
	}

	err = cli.parse()
	if err != nil {
		showHelp()
		os.Exit(1)
	}

	loglet.SetLevel(loglet.INFO)
	if cli.getBool("d") {
		loglet.SetLevel(loglet.DEBUG)
	}

	switch cli.cmd {
	case "show":
		rootOrExit()
		err = show(cli)
	case "version":
		showVersion()
	default:
		showHelp()
	}

	if err != nil {
		loglet.SetOutput(os.Stderr)
		log.Error(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func decode(cli *app) (*model.Report, error) {
	spec := smbios.NewDecoder()
	pcidevs := pci.NewDecoder()

	decoders := []Decoder{
		spec,
		pcidevs,
	}
	for _, d := range decoders {
		err := d.Decode()
		if err != nil {
			return nil, err
		}
	}

	r := new(model.Report)
	shapeSystem(r, spec.GetSystem())
	shapeChassis(r, spec.GetChassis())
	shapeFirmware(r, spec.GetBIOS())
	shapeBaseboard(r, spec.GetBaseboard())
	shapeProcessor(r, spec.GetProcessor())
	shapeMemory(r, spec.GetMemoryDevice())
	shapeDisk(r, pcidevs, cli)
	shapeNetwork(r, pcidevs)
	shapeAccelerater(r, pcidevs)
	shapePowerSupply(r, spec.GetPowerSupply())
	shapeAllPCIDevices(r, pcidevs)
	shapeBMC(r)
	shapeMisc(r)

	r.Version = versionString()

	tm := time.Now()
	r.Timestamp = tm.Unix()
	r.Datetime = tm.Format(time.RFC1123Z)

	return r, nil
}

func rootOrExit() {
	uid := os.Getuid()
	if uid > 0 {
		fmt.Println("must be root")
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("NAME:")
	fmt.Printf("  mox v%s\n", versionString())
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  mox command [command options] [arguments...]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  show")
	fmt.Println("  version")
	fmt.Println("  help")
	fmt.Println()
	fmt.Println("GLOBAL OPTIONS:")
	fmt.Println("  -d,--debug   enabling debug logging")
}

func showVersion() {
	fmt.Println(versionString())
}

func versionString() string {
	if version == "" {
		version = "0.0.0"
	}
	if revision == "" {
		revision = "dev"
	}
	return fmt.Sprintf("%s-%s", version, revision)
}
