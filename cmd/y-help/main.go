package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Set up CLI
	parser := argparse.NewParser("y-help", "Help queue for Electrical and Computer Engineering department at BYU")

	name := parser.String("n",
		"name",
		&argparse.Options{Required: true, Help: "Your name"})

	course := parser.Selector("c",
		"course",
		[]string{"224", "320"},
		&argparse.Options{Required: true, Help: "The course you need help with"})

	isTA := parser.Flag("",
		"ta",
		&argparse.Options{Required: false, Help: "Specifiy if you are a TA"})

	// Processes CLI arguments
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		fmt.Println(parser.Usage(nil))
		return
	}

	var app *tview.Application

	if *isTA {
		app = setupTaApplication(*name, *course)
	} else {
		app = setupApplication(*name, *course)
	}

	// Starts the TUI application
	if err := app.Run(); err != nil {
		panic(err)
	}
}

var app *tview.Application
var cur_focus tview.Primitive

func setupTaApplication(name, course string) *tview.Application {
	app = tview.NewApplication()

	return app
}

func createHeader(course string) *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("ECEn %v", course))
}

func createQueue(queue_name string) *tview.List {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true)
	list.SetTitle(queue_name)
	list.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorDefault).Background(tcell.ColorDefault))
	list.SetFocusFunc(func() {
		list.SetBorderColor(tcell.ColorGreen)
		list.SetTitleColor(tcell.ColorGreen)
	})
	list.SetBlurFunc(func() {
		list.SetBorderColor(tcell.ColorDefault)
		list.SetTitleColor(tcell.ColorDefault)
	})

	return list
}

func createFooter() *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText("q: quit, ←→: switch queues, a: add, d or r: remove")
}

func globalHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == 'q' {
		app.Stop()
	}
	return event
}

func inQueue(name string, list *tview.List) bool {
	return len(list.FindItems(name, "", false, true)) > 0
}

func removeItem(name string, list *tview.List) {
	items := list.FindItems(name, "", false, true)

	if len(items) > 0 {
		list.RemoveItem(items[0])
	}
}

func queueKeyHandler(name string, focusQueue *tview.List, otherQueue *tview.List, modal *tview.Modal, grid *tview.Grid) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' || event.Rune() == 'l' || event.Key() == tcell.KeyTab || event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			app.SetFocus(otherQueue)
			return nil // Don't pass the event
		} else if event.Rune() == 'a' {
			if !inQueue(name, focusQueue) {
				focusQueue.AddItem(name, "", 0, nil)
			}
		} else if event.Rune() == 'r' || event.Rune() == 'd' {
			if inQueue(name, focusQueue) {
				cur_focus = focusQueue
				grid.AddItem(modal, 0, 0, 2, 2, 0, 0, false)
				app.SetFocus(modal)
			}
		}

		return event
	}
}

func setupApplication(name, course string) *tview.Application {
	app = tview.NewApplication()

	header := createHeader(course)
	help_queue := createQueue(" Help Queue ")
	pass_off_queue := createQueue(" Pass-off Queue ")
	footer := createFooter()
	modal := tview.NewModal()
	grid := tview.NewGrid()

	help_queue.SetInputCapture(queueKeyHandler(name, help_queue, pass_off_queue, modal, grid))
	pass_off_queue.SetInputCapture(queueKeyHandler(name, pass_off_queue, help_queue, modal, grid))

	// Set up status bar
	status_bar := tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true).
		SetText("[red]Disconnected[white]")

	// Set up grid
	grid.SetRows(1, 0, 1).
		SetColumns(0, 0).
		AddItem(header, 0, 0, 1, 2, 0, 0, false).
		AddItem(help_queue, 1, 0, 1, 1, 0, 0, true).
		AddItem(pass_off_queue, 1, 1, 1, 1, 0, 0, false).
		AddItem(footer, 2, 0, 1, 1, 0, 0, false).
		AddItem(status_bar, 2, 1, 1, 1, 0, 0, false)

	// Set up alert modal
	modal.SetText("Do you want to remove yourself from the queue?").
		AddButtons([]string{"Remove", "Cancel"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Remove" {
			list := cur_focus.(*tview.List)
			if inQueue(name, list) {
				removeItem(name, list)
			}
		}
		grid.RemoveItem(modal)
		app.SetFocus(cur_focus)
	})
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'y' {
			list := cur_focus.(*tview.List)
			if inQueue(name, list) {
				removeItem(name, list)
			}
			grid.RemoveItem(modal)
			app.SetFocus(cur_focus)
		} else if event.Rune() == 'n' {
			grid.RemoveItem(modal)
			app.SetFocus(cur_focus)
		}
		return event
	})

	// Set up application
	app.SetRoot(grid, true)
	app.EnableMouse(true)

	app.SetInputCapture(globalHandler)

	return app
}
