// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"time"

	"github.com/theparanoids/ashirt/campaign"
	"github.com/theparanoids/ashirt/evidence"
	"github.com/theparanoids/ashirt/httpclient"
	"github.com/theparanoids/ashirt/screenshotclient/config"
	"github.com/theparanoids/ashirt/screenshotclient/database"
	"github.com/theparanoids/ashirt/screenshotclient/fuzzyfilefinder"
	"github.com/theparanoids/ashirt/screenshotclient/screencapture"
	"github.com/theparanoids/ashirt/screenshotclient/settings"

	_ "github.com/andlabs/ui/winmanifest"
	"github.com/getlantern/systray"
	"github.com/marcsauter/single"
	_ "github.com/mattn/go-sqlite3"
	"github.com/radovskyb/watcher"
)

var (
	// NOTE(joe): This whole thing seems pretty racey. We need to get some encapsulation and locking up in this
	s               *settings.Settings
	camp            *campaign.Campaign
	client          *httpclient.Client
	ScreenshotRegex = regexp.MustCompile(`\AScreen Shot`)
	running         = false
)

func main() {
	//only run one instance of ashirt
	s := single.New("ashirt")
	if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
		log.Fatal("Another instance of the app is already running, exiting")
	} else if err != nil {
		// Another error occurred, might be worth handling it as well
		log.Fatalf("failed to acquire exclusive app lock: %v", err)
	}

	defer s.TryUnlock()
	systray.Run(onReady, onExit)
}

func onReady() {
	var err error

	s, err = settings.GetSettings()
	if err != nil {
		log.Fatal("unable to get settings: ", err)
	}

	err = database.SetupDB()
	if err != nil {
		log.Fatalln(err)
	}

	client = httpclient.NewClient(&s.Config)

	icon, err := getIcon("assets/ashirtoff.ico")
	if err != nil {
		log.Fatalln("unable to get icon: ", err)
	}

	systray.SetIcon(icon)

	runCampaign := systray.AddMenuItem("Run Campaign: "+s.Config.Campaign, "runCampaign")
	stopCampaign := systray.AddMenuItem("Stop Campaign: "+s.Config.Campaign, "stopCampaign")
	changeCampaign := systray.AddMenuItem("Change Campaign", "changeCampaign")

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quits this app")
	fmt.Println("Ashirt is running")
	fmt.Println("You should see a white shirt icon in your tool bar.")
	fmt.Println("Project Screenshot Directory: " + s.ProjectDefault)
	fmt.Println("OS Default Screenshot Directory: " + s.HomeDefault)

	//used for async publishing background thread reads sqlite db then pushes screenshots to ashirt api
	go func() {
		for {
			time.Sleep(5 * time.Second)

			screenshots, err := database.GetUnsubmitted()
			if err != nil {
				log.Println("failed to retrieve unsubmitted screenshots: ", err)
			}

			for _, screenshot := range screenshots {
				//check to see if screenshot exists before pushing it to the api.
				if !evidence.Exists(screenshot.FullPath) {
					log.Println("unsubmitted file does not exist")
					err = database.MarkError(screenshot.ID)
					if err != nil {
						log.Println("failed to record error for screenshot: ", err)
					}

					continue
				}

				//use the fuzzy file finder if the file name is not found from the watcher event and the file fixer fails
				//the fuzzy file finder finds the first closest match to the file name
				if !fuzzyfilefinder.Exists(screenshot.FullPath) {
					if newPath, err := fuzzyfilefinder.FuzzyFinder(screenshot.FullPath, 8); err != nil {
						log.Println("Using fuzzy file finder failed: ", err)
					} else {
						screenshot.FullPath = newPath
					}
				}

				//get file hash
				if fileHash, err := evidence.GetFileHash(screenshot.FullPath); err != nil {
					log.Println("failed to submit screenshot: ", err)
					screenshot.FileHash = fileHash
				}

				err = evidence.ProcessFile(&screenshot, s)
				if err != nil {
					log.Println("failed to process file: ", err)
					continue
				}

				err := database.UpdateScreenshots(screenshot)
				if err != nil {
					log.Println("failed to update database with new file location: ", err)
					err = database.MarkError(screenshot.ID)
					if err != nil {
						log.Println("failed to record error for screenshot: ", err)
					}

					continue
				}

				//upload the file to the api
				err = client.UploadScreenshot(&screenshot)
				if err != nil {
					log.Println("unable to upload screenshot: ", err)
					err = database.MarkError(screenshot.ID)
					if err != nil {
						log.Println("failed to mark error occurring: ", err)
					}
					continue
				}

				err = database.MarkSubmitted(screenshot.ID)
				if err != nil {
					log.Println("failed to mark screenshot as submitted: ", err)
				}
			}
		}
	}()

	//Ashirt System Tray
	go func() {
		stopCampaign.Hide()
		w := watcher.New()
		for {
			select {
			case <-runCampaign.ClickedCh:
				w.Close()
				icon, err := getIcon("assets/ashirt.ico")
				if err != nil {
					log.Println("unable to get icon: ", err)
					continue
				}
				systray.SetIcon(icon)

				go func() {
					camp, err := getCampaign()
					if err != nil {
						log.Println("unable to get campaign: ", err)
					}
					w = watcher.New()
					err = screencapture.SetScreenshotDirectory(s.ProjectDefault)
					if err != nil {
						log.Println("unable to set screenshot directory: ", err)
					}
					shirtWatcher(w, s.ProjectDefault, *camp)
				}()

				runCampaign.Hide()
				stopCampaign.Show()
				changeCampaign.Hide()
				running = true
			case <-stopCampaign.ClickedCh:
				icon, err := getIcon("assets/ashirtoff.ico")
				if err != nil {
					log.Println("unable to get icon: ", err)
					continue
				}
				systray.SetIcon(icon)
				err = screencapture.ResetScreenshotDirectory(s.HomeDefault)
				if err != nil {
					log.Println("unable to set screenshot directory: ", err)
					continue
				}
				w.Close()
				stopCampaign.Hide()
				runCampaign.Show()
				changeCampaign.Show()
				running = false
				log.Println("watcher stopped")
			case <-changeCampaign.ClickedCh:
				_, err = spawnUI()
				if err != nil {
					log.Println("failed to launch config program", err)
					continue
				}
				camp, err = getCampaign()
				if err != nil {
					log.Println("unable to get campaign: ", err)
				}

				s, err = settings.GetSettings()
				if err != nil {
					log.Println("unable to get settings: ", err)
				}

				runCampaign.SetTitle("Run Campaign: " + camp.Name)
				stopCampaign.SetTitle("Stop Campaign: " + camp.Name)
			case <-mQuit.ClickedCh:
				// only clean up if we're running to avoid removing/changing back twice
				if running {
					err = screencapture.ResetScreenshotDirectory(s.HomeDefault)
					if err != nil {
						log.Println("unable to reset screenshot directory: ", err)
					}
				}
				systray.Quit()
				return
			}
		}
	}()
}

//loads Campign
func getCampaign() (*campaign.Campaign, error) {
	c, err := config.Read()
	if err != nil {
		return nil, err
	}

	return campaign.New(c.CampaignID, c.Campaign), nil
}

// used to watch a directory for screenshots from the os
func shirtWatcher(w *watcher.Watcher, projectDir string, camp campaign.Campaign) {
	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.

	// Only notify rename and move events.
	w.FilterOps(watcher.Create)
	w.AddFilterHook(watcher.RegexFilterHook(ScreenshotRegex, false))
	//start a thread for the listener
	go func() {
		for {
			select {
			case event := <-w.Event:
				var fileName = event.Name()
				screenshot, err := evidence.CollectScreenshot(projectDir, fileName)
				if err != nil {
					if err != evidence.ErrCancelled {
						log.Println("unable to collect screenshot: ", err)
					}

					continue
				}

				screenshot.Campaign = camp
				//log the screenshot info to the database to be later picked up by the
				// async thread that will push the file to the api
				err = database.InsertScreenshot(screenshot)
				if err != nil {
					log.Println("failed to insert screenshot into database: ", err)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	log.Println("watcher started on: " + projectDir)
	// Watch this folder for changes.
	if err := w.Add(projectDir); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func onExit() {
	// Cleaning stuff here.
}

func getIcon(s string) ([]byte, error) {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// spawnUI runs the program for changing the current operation. This is required
// to be build as a separate application instead of spawning a new window
// because the Go systray package and ui package both want control of the main
// thread. We should look into switching out UI toolkits for something like Qt
// or GTK that will provide the behavior for all our desired UI functionality so
// that it can all work together. Without doing so this can introduce racey
// behavior
func spawnUI() (string, error) {
	var (
		cmdOut []byte
		err    error
	)
	// runs this command defaults read com.apple.screencapture location
	cmdName := "./dropdown"
	cmdArgs := []string{""}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return "", err
	}

	sha := string(cmdOut)
	return sha, nil
}
