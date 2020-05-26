// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package main

import (
	"fmt"
	"log"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sqweek/dialog"
	"github.com/theparanoids/ashirt/campaign"
	"github.com/theparanoids/ashirt/httpclient"
	"github.com/theparanoids/ashirt/screenshotclient/config"
	"github.com/theparanoids/ashirt/screenshotclient/screenshotview"
)

var (
	currentConfig   config.Config
	currentCampaign campaign.Campaign
	client          *httpclient.Client
	mainwin         *ui.Window
)

func makeConfigPage() ui.Control {
	vbox := ui.NewVerticalBox()

	group := ui.NewGroup("ASHIRT Config")
	group.SetMargined(true)
	vbox.Append(group, true)

	campaigns, err := client.GetCampaigns()
	if err != nil {
		dialog.Message(fmt.Sprintf("Could not load campaigns: %s", err)).Title("Failure").Error()
		ui.Quit()
	}

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)

	entryWorkingFolder := ui.NewEntry()
	entryWorkingFolder.SetText(currentConfig.WorkingFolder)

	buttonChangeFolder := ui.NewButton("Change Folder")
	buttonChangeFolder.OnClicked(func(button *ui.Button) {
		directory, err := dialog.Directory().Title("Select Working Folder").Browse()
		if err != nil {
			if err.Error() != "Cancelled" {
				dialog.Message(fmt.Sprintf("Could not change folder: %s", err)).Title("Failure").Error()
			}
		} else {
			currentConfig.WorkingFolder = directory
			entryWorkingFolder.SetText(directory)
		}
	})

	workingDirRow := ui.NewHorizontalBox()
	workingDirRow.SetPadded(true)
	workingDirRow.Append(entryWorkingFolder, true)
	workingDirRow.Append(buttonChangeFolder, false)
	entryForm.Append("Working Folder", workingDirRow, true)

	entryAPIUrl := ui.NewEntry()
	entryAPIUrl.SetText(currentConfig.APIURL)
	entryForm.Append("API URL", entryAPIUrl, false)

	cbox := ui.NewCombobox()
	// var used to store the cbox id
	var cboxID int
	for i, campaign := range campaigns {
		cbox.Append(campaign.Name)

		//set default from config
		if campaign.ID == currentConfig.CampaignID {
			cboxID = i
			currentCampaign = campaign
		}
	}

	//set default selection
	cbox.SetSelected(cboxID)
	cbox.OnSelected(func(*ui.Combobox) {
		if len(campaigns) > 0 {
			selected := cbox.Selected()
			currentCampaign = campaigns[selected]
		}
	})

	entryForm.Append("Campaign", cbox, false)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	hbox := ui.NewHorizontalBox()
	vbox.Append(hbox, false)

	saveButton := ui.NewButton("Save")
	hbox.Append(saveButton, false)

	saveButton.OnClicked(func(*ui.Button) {
		currentConfig.Campaign = currentCampaign.Name
		currentConfig.CampaignID = currentCampaign.ID
		currentConfig.APIURL = entryAPIUrl.Text()
		currentConfig.WorkingFolder = entryWorkingFolder.Text()
		err := currentConfig.Save()
		if err != nil {
			saveWarning()
		} else {
			saveMessage()
		}

	})

	return vbox
}

func saveWarning() {
	dialog.Message("Could not save config").Title("Failure").Error()
}

func saveMessage() {
	dialog.Message("Config saved.").Title("Success").Info()
}

func setupUI() {
	mainwin = ui.NewWindow("A Shirt Config", 720, 320, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	// if table view loading failes quit ui
	vboxTableView, err := screenshotview.BuildScreenshotTableView()
	if err != nil {
		mainwin.Destroy()
	}
	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)
	tab.Append("Config", makeConfigPage())
	tab.Append("Log", vboxTableView)
	tab.SetMargined(0, true)
	mainwin.Show()
}

func main() {
	var err error

	currentConfig, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}

	client = httpclient.NewClient(&currentConfig)

	ui.Main(setupUI)
}
