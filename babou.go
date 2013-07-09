// `babou` allows you to spawn an instance of the babou private tracker.
package main

import (
	flag "flag"
	fmt "fmt"

	web "github.com/drbawb/babou/app" // The babou application: composed of a server and muxer.
	tracker "github.com/drbawb/babou/tracker"

	libBabou "github.com/drbawb/babou/lib" // Core babou libraries
	//libDb "github.com/drbawb/babou/lib/db"

	os "os"
	signal "os/signal"
	syscall "syscall"
)

func main() {
	//Output welcome message:
	fmt.Println("babou fast like veyron.")

	//Parse command line flags
	appSettings := parseFlags()

	//Trap signals from the parent OS
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go trapSignals(c)

	webServerIO := make(chan int, 0)
	trackerIO := make(chan int, 0)

	if *appSettings.FullStack == true || *appSettings.WebStack == true {
		// Start web-server
		fmt.Printf("Starting web")
		server := web.NewServer(appSettings, webServerIO)
		// Receive SIGNALs from web server.

		go server.Start()
	}

	if *appSettings.FullStack == true || *appSettings.TrackerStack == true {
		// Start tracker
		fmt.Printf("Starting tracker \n")
		server := tracker.NewServer(appSettings, trackerIO)

		go server.Start()
	}

	// Block on server IOs
	for {
		select {
		case webMessage := <-webServerIO:
			if webMessage == libBabou.WEB_SERVER_STARTED {
				fmt.Println("Server has started sucessfully")
			}
		case trackerMessage := <-trackerIO:
			if trackerMessage == libBabou.TRACKER_SERVER_START {
				fmt.Println("Tracker has started successfully")
			}
		}
	}
}

func trapSignals(c chan os.Signal) {
	for sig := range c {
		if sig == syscall.SIGINT || sig == syscall.SIGQUIT || sig == syscall.SIGTERM {
			// Shutdown gracefully.
			fmt.Println("\nbabou is packing up his things ...")

			//TODO: Probably block on webserver shutdown [instant]
			// 	as well as a concurrent block on app shutdown.
			// Exit when they're both finished.
			fmt.Println("\nwaiting for webserver to shutdown...")
			fmt.Println("\nwaiting for tracker to shutdown...")

			os.Exit(0)
		} else if sig == syscall.SIGKILL {
			// Get out get out get out!!!
			os.Exit(2)
		}
	}
}

func parseFlags() *libBabou.AppSettings {
	appSettings := &libBabou.AppSettings{}

	appSettings.Debug = flag.Bool("debug", false,
		"Logs debug information to console.")

	appSettings.WebStack = flag.Bool("web-stack", false,
		"Enables the web application server.")
	appSettings.TrackerStack = flag.Bool("track-stack", false,
		"Enables the torrent tracker.")
	appSettings.FullStack = flag.Bool("full-stack", true,
		"Enables the full application stack. - Disabled if track-stack or web-stack are set.")

	appSettings.WebPort = flag.Int("web-port", 8080,
		"Sets the web application's port number.")
	appSettings.TrackerPort = flag.Int("track-port", 4200,
		"Sets the tracker's listening port number.")

	flag.Parse()

	// If the user has configured their own stack options, do not use the full stack.
	if *appSettings.WebStack == true || *appSettings.TrackerStack == true {
		*appSettings.FullStack = false
	}

	return appSettings
}
