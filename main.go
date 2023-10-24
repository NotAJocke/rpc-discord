package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/joho/godotenv"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apps := []string{"xcode", "figma"}
	app := flag.String("app", "", "The app you want to start, accepted values: "+strings.Join(apps, ", "))
	workspace := flag.String("workspace", "", "The workspace you're working in")
	ws := flag.String("ws", "", "The workspace you're working in")
	flag.Parse()

	if *app == "" {
		fmt.Println("App not set, use 'app' flag.\nAccepted values: " + strings.Join(apps, ", "))
		return
	}
	if *workspace == "" && *ws == "" {
		fmt.Println("Workspace not set, use 'workspace' or 'ws' flag")
		return
	}

	var appId string
	var activity client.Activity
	switch *app {
	case "xcode":
		appId = os.Getenv("XCODE_APP_ID")
		activity = getXcodeActivity(*workspace + *ws)

	case "figma":
		appId = os.Getenv("FIGMA_APP_ID")
		activity = getFigmaActivity(*workspace + *ws)
	}

	err = client.Login(appId)
	if err != nil {
		fmt.Printf("Couldn't log client in: %s", err)
		return
	}

	err = client.SetActivity(activity)
	if err != nil {
		fmt.Printf("Couln't set client activity: %s", err)
		return
	}

	fmt.Println("Started activity")

	<-sigs
}

func getXcodeActivity(workspace string) client.Activity {
	startTime := time.Now()

	return client.Activity{
		Details:    "Editing Swift files",
		State:      "Workspace: " + workspace,
		LargeImage: "xcode",
		Timestamps: &client.Timestamps{
			Start: &startTime,
		},
	}
}

func getFigmaActivity(workspace string) client.Activity {
	startTime := time.Now()

	return client.Activity{
		Details:    "Designing user interfaces",
		State:      "Workspace: " + workspace,
		LargeImage: "figma",
		Timestamps: &client.Timestamps{
			Start: &startTime,
		},
	}
}
