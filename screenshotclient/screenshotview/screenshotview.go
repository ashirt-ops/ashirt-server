// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package screenshotview

import (
	"fmt"
	"github.com/andlabs/ui"
	"github.com/sqweek/dialog"
	"github.com/theparanoids/ashirt/evidence"
	"github.com/theparanoids/ashirt/screenshotclient/database"
	"log"
	"time"
)

func loadDatabaseWarning() {
	dialog.Message("Could not load screenshot info from sql databse.").Title("Failure").Error()
}

type modelHandler struct {
	screenshots []evidence.Screenshot
}

func newModelHandler() (*modelHandler, error) {

	m := new(modelHandler)
	results, err := database.ListScreenshots()

	if err != nil {
		return nil, err
	}
	m.screenshots = results
	return m, nil
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""), // column 0 text
		ui.TableString(""), // column 1 text
		ui.TableString(""), // column 2 text
		ui.TableString(""), // column 4 text
	}
}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return len(mh.screenshots)
}

//id,FileName,FullPath,Description
func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	switch column {
	case 0:
		return ui.TableString(fmt.Sprintf("%d", mh.screenshots[row].ID))
	case 1:
		return ui.TableString(fmt.Sprintf("%s", mh.screenshots[row].FileName))
	case 2:
		return ui.TableString(fmt.Sprintf("%s", mh.screenshots[row].Description))
	case 3:
		//convert from nano seconds to seconds
		tm := time.Unix(mh.screenshots[row].OccurrenceTimestamp/1000000000, 0)
		return ui.TableString(tm.Format(time.RFC3339))
	case 4:
		return nil
	}
	return nil
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {

}

//adjusts the size of the table based on the screenshot struct
func (mh *modelHandler) RefreshTable(m *ui.TableModel) {

	screenshots, err := database.ListScreenshots()
	if err != nil {
		log.Print("Failed to load screenshot table")
		loadDatabaseWarning()
		return
	}
	var size = len(mh.screenshots)
	if len(screenshots) == size {
		mh.screenshots = screenshots
		for i := 1; i <= len(screenshots); i++ {
			m.RowChanged(i - 1)
		}
		return
	}
	if len(screenshots) > len(mh.screenshots) {

		mh.screenshots = screenshots
		for i := 1; i <= size; i++ {
			m.RowChanged(i - 1)
		}

		diff := len(screenshots) - size
		for i := 1; i <= diff; i++ {
			m.RowInserted(size - 1 + i)
		}

		return
	}
	if len(screenshots) < size {
		mh.screenshots = screenshots
		for i := 1; i <= len(screenshots); i++ {
			m.RowChanged(i - 1)
		}

		diff := size - len(screenshots)
		for i := 1; i <= diff; i++ {
			m.RowDeleted(len(screenshots) - 1 + i)
		}

		return
	}
}

//BuildScreenshotTableView creates the ui controls for the SQL table screenshosts
func BuildScreenshotTableView() (ui.Control, error) {

	mh, err := newModelHandler()
	if err != nil {
		loadDatabaseWarning()
		return nil, err
	}
	model := ui.NewTableModel(mh)

	table := ui.NewTable(&ui.TableParams{
		Model:                         model,
		RowBackgroundColorModelColumn: 3,
	})

	vbox := ui.NewVerticalBox()

	group := ui.NewGroup("Screenshot Log")
	group.SetMargined(true)
	vbox.Append(group, true)
	group.SetChild(table)

	table.AppendTextColumn("ID",
		0, ui.TableModelColumnNeverEditable, nil)

	table.AppendTextColumn("File Name",
		1, ui.TableModelColumnNeverEditable, nil)

	table.AppendTextColumn("Description",
		2, ui.TableModelColumnNeverEditable, nil)

	table.AppendTextColumn("Collected Date",
		3, ui.TableModelColumnNeverEditable, nil)

	refreshButton := ui.NewButton("Refresh")

	refreshButton.OnClicked(func(button *ui.Button) {
		mh.RefreshTable(model)
	})

	vbox.Append(refreshButton, false)

	return vbox, nil

}
