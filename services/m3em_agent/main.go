package main

import (
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/m3db/m3em/agent"
	"github.com/m3db/m3em/generated/proto/m3em"

	"github.com/m3db/m3nsch/m3nsch_server/tcp"
	"github.com/m3db/m3x/instrument"

	xlog "github.com/m3db/m3x/log"
	"github.com/pborman/getopt"
	"github.com/spf13/viper"
	"github.com/uber-go/tally"
	tallym3 "github.com/uber-go/tally/m3"
	"google.golang.org/grpc"
	validator "gopkg.in/validator.v2"
)

type configurationFile struct {
	Server  serverConfig  `yaml:"server"`
	Metrics metricsConfig `yaml:"metrics"`
	Agent   agentConfig   `yaml:"agent"`
}

type agentConfig struct {
	WorkingDir  string        `yaml:"workingDir" validate:"nonzero"`
	StartupCmds []execCommand `yaml:"startupCmds"`
	ReleaseCmds []execCommand `yaml:"releaseCmds"`
}

type execCommand struct {
	Path string   `yaml:"path" validate:"nonzero"`
	Args []string `yaml:"args"`
}

type metricsConfig struct {
	Prefix     string                `yaml:"prefix"`
	SampleRate float64               `yaml:"sampleRate" validate:"min=0.01,max=1.0"`
	M3         tallym3.Configuration `yaml:"m3"         validate:"nonzero"`
}

type serverConfig struct {
	ListenAddress string `yaml:"listenAddress" validate:"nonzero"`
	DebugAddress  string `yaml:"debugAddress"  validate:"nonzero"`
}

func readConfiguration(filename string) (configurationFile, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)

	var conf configurationFile

	if err := viper.ReadInConfig(); err != nil {
		return conf, err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return conf, err
	}

	if err := validator.Validate(conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func main() {
	var (
		configFile = getopt.StringLong("config-file", 'f', "", "Configuration file")
	)
	getopt.Parse()
	if len(*configFile) == 0 {
		getopt.Usage()
		return
	}

	logger := xlog.NewLogger(os.Stdout)
	conf, err := readConfiguration(*configFile)
	if err != nil {
		logger.Fatalf("unable to read configuration file: %v", err.Error())
	}

	// pprof server
	go func() {
		if err := http.ListenAndServe(conf.Server.DebugAddress, nil); err != nil {
			logger.Fatalf("unable to serve debug server: %v", err)
		}
	}()
	logger.Infof("serving pprof endpoints at: %v", conf.Server.DebugAddress)

	reporter, err := conf.Metrics.M3.NewReporter()
	if err != nil {
		logger.Fatalf("could not connect to metrics: %v", err)
	}
	scope, _ := tally.NewCachedRootScope(conf.Metrics.Prefix, nil, reporter, time.Second, tally.DefaultSeparator)
	listener, err := tcp.NewTCPListener(conf.Server.ListenAddress, 3*time.Minute)
	if err != nil {
		logger.Fatalf("could not create TCP Listener: %v", err)
	}

	iopts := instrument.NewOptions().
		SetLogger(logger).
		SetMetricsScope(scope).
		SetMetricsSamplingRate(conf.Metrics.SampleRate)

	agentOpts := agent.NewOptions(iopts).
		SetWorkingDirectory(conf.Agent.WorkingDir).
		SetInitHostResourcesFn(hostInitFnMaker(conf.Agent.StartupCmds, logger)).
		SetReleaseHostResourcesFn(hostReleaseFnMaker(conf.Agent.ReleaseCmds, logger)).
		SetExecGenFn(execGenFn)

	agentService, err := agent.New(agentOpts)
	if err != nil {
		logger.Fatalf("unable to create agentService: %v", err)
	}
	server := grpc.NewServer(grpc.MaxConcurrentStreams(16384))
	m3em.RegisterOperatorServer(server, agentService)
	logger.Infof("serving agent endpoints at %v", listener.Addr().String())
	if err := server.Serve(listener); err != nil {
		logger.Fatalf("could not serve: %v", err)
	}
}

func hostInitFnMaker(cmds []execCommand, logger xlog.Logger) agent.InitHostResourcesFn {
	return func() error {
		if len(cmds) == 0 {
			logger.Info("no startup commands specified, skipping.")
			return nil
		}
		for _, cmd := range cmds {
			osCmd := exec.Command(cmd.Path, cmd.Args...)
			logger.Infof("attempting to execute startup cmd: %+v", osCmd)
			output, err := osCmd.CombinedOutput()
			if err != nil {
				logger.Errorf("unable to execute cmd, err: %v", err)
				return err
			}
			logger.Infof("successfully ran cmd, output: [%v]", string(output))
		}
		return nil
	}
}

func hostReleaseFnMaker(cmds []execCommand, logger xlog.Logger) agent.ReleaseHostResourcesFn {
	return func() error {
		if len(cmds) == 0 {
			logger.Info("no release commands specified, skipping.")
			return nil
		}
		for _, cmd := range cmds {
			osCmd := exec.Command(cmd.Path, cmd.Args...)
			logger.Infof("attempting to execute release cmd: %+v", osCmd)
			output, err := osCmd.CombinedOutput()
			if err != nil {
				logger.Errorf("unable to execute cmd, err: %v", err)
				return err
			}
			logger.Infof("successfully ran cmd, output: [%v]", string(output))
		}
		return nil
	}
}

func execGenFn(binary string, config string) (string, []string) {
	return binary, []string{"-f", config}
}
