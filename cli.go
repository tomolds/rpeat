package rpeat

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func ParseCommandLine() ServerConfig {
	var c ServerConfig

	var HOST, PORT, jobs, config, auth, timezone string
	var useHttps, clean, nostate, nohistory, demo bool
	var cert, key string

	flag.Usage = func() {
		fmt.Printf(`
rpeat-server

rpeat is designed as a drop-in replacement for the ubiquitous cron. What it provides
is a lightweight 'rpeat' server that consume "Jobs", which not only have a start
time, but myriad options such as permissioning, logging, alerts, retry, and more!

All tasks (a.k.a. jobs) are defined in a json file, passed in as -jobs or defined in JobsFiles field of
the configuration json file (which itself is specified via -config)

--config

--jobs

--auth

`)
		fmt.Printf("\n\n")
		fmt.Printf("Usage: rpeat-server [options]>\n\n")
		flag.PrintDefaults()
		fmt.Printf("...\n")
	}

	HOME, _ := os.LookupEnv("RPEAT_HOME")
	flag.StringVar(&HOME, "HOME", HOME, "HOME")

	flag.BoolVar(&demo, "demo", false, "create demo config, jobs and auth file and run")
	flag.StringVar(&HOST, "HOST", "localhost", "address of rpeat server instance")
	flag.StringVar(&PORT, "PORT", "4334", "port of rpeat server instance")
	flag.BoolVar(&useHttps, "https", true, "serve via TLS/https. Requires SSL cert and key to be specified (true)")
	flag.BoolVar(&clean, "clean", false, "ignore previous jobs state and history. Same as specifying both --nostate --noclean (false)")
	flag.BoolVar(&nostate, "nostate", false, "ignore previous jobs state if any (false)")
	flag.BoolVar(&nohistory, "nohistory", false, "ignore previous jobs history if any (false)")
	flag.StringVar(&cert, "cert", "server.crt", "TLS cert")
	flag.StringVar(&key, "key", "server.key", "TLS key")
	flag.StringVar(&jobs, "jobs", "", "comma seperated array of jobs files")
	flag.StringVar(&config, "config", "", "rpeat configuration file")
	flag.StringVar(&auth, "auth", "", "rpeat auth file")
	flag.StringVar(&timezone, "timezone", "", "timezone of server (defaults to system tz)")
	flag.Parse()

	if demo {
		CreateDemo()
		u, _ := user.Current()
		//log.Printf("user: %s:", u.Username)
		//ServerLogger.Printf("running demo for: %s", user)
		config = filepath.Join(u.HomeDir, ".rpeat", "demo", "config.json")
	}

	jobsfiles := strings.Split(jobs, ",")

	if config != "" {
		var err error
		c, err = LoadServerConfig(config, true)
		if err != nil {
			ServerLogger.Fatal("problem processing config file. Terminating")
		}
		// replace any configurations with passed variables
		flaggedArgs := FlagArgs(flag.CommandLine)
		for _, flag := range flaggedArgs {
			switch flag {
			case "HOST":
				c.HOST = HOST
			case "PORT":
				c.PORT = PORT
			case "https":
				c.Https = useHttps
			case "clean":
				c.Clean = true
				c.KeepHistory = false
			case "nostate":
				c.Clean = nostate
			case "nohistory":
				c.KeepHistory = !nohistory
			case "jobs":
				c.JobsFiles = jobsfiles
			case "cert":
				c.TLS.Cert, _ = filepath.Abs(cert)
			case "key":
				c.TLS.Key, _ = filepath.Abs(key)
			case "auth":
				c.AuthFile, _ = filepath.Abs(auth)
			case "timezone":
				c.Timezone = timezone
			case "HOME":
				c.HOME, _ = filepath.Abs(HOME)
			}
		}

		c.GITCOMMIT = GITCOMMIT
		c.BUILDDATE = BUILDDATE
	}

	for i, jobsfile := range c.JobsFiles {
		if exists, err := FileExists(jobsfile); !exists {
			ServerLogger.Fatal(fmt.Sprintf("problem reading jobs file %s: [%s]", jobsfile, err))
		}
		c.JobsFiles[i], _ = filepath.Abs(jobsfile)
	}

	return c

}
