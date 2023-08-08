package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Processes CLI arguments

	// Starts the TUI application
	app := setupApplication()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

var cur_focus tview.Primitive

func setupApplication() *tview.Application {
	app := tview.NewApplication()

	// Set up header
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("ECEn 224")

	help_queue := tview.NewList()
	pass_off_queue := tview.NewList()
	modal := tview.NewModal()
	grid := tview.NewGrid()

	// Set up the help queue
	help_queue.ShowSecondaryText(false)
	help_queue.SetBorder(true)
	help_queue.SetTitle(" Help Queue ")
	help_queue.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' || event.Rune() == 'l' || event.Rune() == 9 /* Tab */ {
			app.SetFocus(pass_off_queue)
			return nil // Don't pass the event
		} else if event.Rune() == 'a' {
			help_queue.AddItem("New Person", "", 0, nil)
		} else if event.Rune() == 'r' || event.Rune() == 'd' {
			cur_focus = help_queue
			grid.AddItem(modal, 0, 0, 2, 2, 0, 0, false)
			app.SetFocus(modal)
		}
		return event
	})
	help_queue.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		help_queue.RemoveItem(index)
	})

	// Set up pass-off queue
	pass_off_queue.ShowSecondaryText(false)
	pass_off_queue.SetBorder(true)
	pass_off_queue.SetTitle(" Pass-off Queue ")
	pass_off_queue.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' || event.Rune() == 'l' || event.Rune() == 9 /* Tab */ {
			app.SetFocus(help_queue)
			return nil // Don't pass the event
		} else if event.Rune() == 'a' {
			pass_off_queue.AddItem("New Person", "", 0, nil)
		} else if event.Rune() == 'r' || event.Rune() == 'd' {
			cur_focus = pass_off_queue
			grid.AddItem(modal, 0, 0, 2, 2, 0, 0, false)
			app.SetFocus(modal)
		}
		return event
	})

	// Set up footer
	footer := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText("q: quit, tab: switch queues, a: add, d or r: remove")

	// Set up grid
	grid.SetRows(1, 0, 1).
		SetColumns(0, 0)

	grid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(help_queue, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(pass_off_queue, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(footer, 2, 0, 1, 2, 0, 0, false)

	// Set up alert modal

	modal.SetText("Do you want to remove yourself from the queue?").
		AddButtons([]string{"Remove", "Cancel"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Remove" {
			cur_focus.(*tview.List).RemoveItem(cur_focus.(*tview.List).GetCurrentItem())
		}
		grid.RemoveItem(modal)
		app.SetFocus(cur_focus)
	})

	// Set up application
	app.SetRoot(grid, true)
	app.EnableMouse(true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	// Check for new people in the help queue every second
	go func() {
		for {
			time.Sleep(1 * time.Second)
			app.QueueUpdateDraw(func() {
				help_queue.AddItem("Philip Lundrigan", "", 0, nil)
				pass_off_queue.AddItem("Philip Lundrigan", "", 0, nil)
			})
		}
	}()

	return app
}