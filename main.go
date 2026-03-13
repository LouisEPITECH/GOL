package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Cells [40][160]bool

func main() {
	screen := initScreen()
	drawStyle := tcell.StyleDefault.Background(tcell.ColorDarkViolet).Foreground(tcell.ColorDarkViolet)

	screen.Clear()
	cells := firstGeneration()
	drawCells(cells, screen, drawStyle)

	for {
		time.Sleep(50 * time.Millisecond)
		cells = nextGeneration(cells)

		screen.Clear()
		drawCells(cells, screen, drawStyle)
		screen.Show()

		if screen.HasPendingEvent() {
			handleEvent(screen)
		}
	}
}

func handleEvent(screen tcell.Screen) {
	event := screen.PollEvent()
	switch event := event.(type) {
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventKey:
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
			screen.Fini()
			os.Exit(0)
		}
	}
}

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorLightGreen).Foreground(tcell.ColorLightGreen)
	screen.SetStyle(defStyle)

	return screen
}

func drawCells(cells Cells, screen tcell.Screen, style tcell.Style) {
	for rowIndex := 0; rowIndex < len(cells); rowIndex++ {
		for colIndex := 0; colIndex < len(cells[rowIndex]); colIndex++ {
			if cells[rowIndex][colIndex] {
				screen.SetContent(colIndex, rowIndex, ' ', nil, style)
			}
		}
	}
}

func nextGeneration(cells Cells) Cells {
	var nextCells Cells

	for rowIndex := 0; rowIndex < len(cells); rowIndex++ {

		for colIndex := 0; colIndex < len(cells[rowIndex]); colIndex++ {
			neighbourCount := calculateNeighbourCount(cells, rowIndex, colIndex)
			alive := cells[rowIndex][colIndex]

			if alive && (neighbourCount == 2 || neighbourCount == 3) {
				nextCells[rowIndex][colIndex] = true
			} else if !alive && neighbourCount == 3 {
				nextCells[rowIndex][colIndex] = true
			} else {
				nextCells[rowIndex][colIndex] = false
			}
		}
	}

	return nextCells
}

func firstGeneration() Cells {
	rand.Seed(time.Now().UnixNano())
	var cells Cells
	rowStart := len(cells)/2 - 5
	rowEnd := len(cells)/2 + 5
	colStart := len(cells[0])/2 - 10
	colEnd := len(cells[0])/2 + 10

	for rowIndex := rowStart; rowIndex < rowEnd; rowIndex++ {
		for colIndex := colStart; colIndex < colEnd; colIndex++ {
			if rand.Intn(4) == 0 {
				cells[rowIndex][colIndex] = true
			}
		}
	}

	return cells
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func calculateNeighbourCount(cells Cells, curRow, curCol int) int {
	rowStart := max(curRow-1, 0)
	rowEnd := min(curRow+1, len(cells)-1)
	colStart := max(curCol-1, 0)
	colEnd := min(curCol+1, len(cells[0])-1)
	neighbourCount := 0

	for rowIndex := rowStart; rowIndex <= rowEnd; rowIndex++ {
		for colIndex := colStart; colIndex <= colEnd; colIndex++ {
			isRefCell := rowIndex == curRow && colIndex == curCol

			if !isRefCell && cells[rowIndex][colIndex] {
				neighbourCount++
			}
		}
	}

	return neighbourCount
}

func printCells(cells Cells) {
	for rowIndex := 0; rowIndex < len(cells); rowIndex++ {
		for colIndex := 0; colIndex < len(cells[rowIndex]); colIndex++ {
			if cells[rowIndex][colIndex] {
				fmt.Print("O")
			} else {
				fmt.Print(".")
			}
		}

		fmt.Println()
	}
}
