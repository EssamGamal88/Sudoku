package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	size     = 9
	cellSize = 60 // px
)

var initialBoard = [9][9]byte{
	{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
	{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
	{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
	{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
	{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
	{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
	{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
	{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
	{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
}

func main() {
	a := app.New()
	w := a.NewWindow("Sudoku Solver (Animated)")

	board := initialBoard
	texts := make([][]*canvas.Text, size)
	grid := container.NewGridWrap(fyne.NewSize(float32(cellSize), float32(cellSize)))

	for i := 0; i < size; i++ {
		texts[i] = make([]*canvas.Text, size)
		for j := 0; j < size; j++ {
			val := string(board[i][j])
			if val == "." {
				val = ""
			}
			txt := canvas.NewText(val, color.Black)
			txt.TextSize = 36
			txt.Alignment = fyne.TextAlignCenter
			texts[i][j] = txt

			bg := canvas.NewRectangle(color.White)
			if board[i][j] == '.' {
				bg.FillColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}
			}
			bg.Resize(fyne.NewSize(float32(cellSize), float32(cellSize)))

			const thin = 2
			const thick = 6

			var borders []fyne.CanvasObject

			// Top border
			topThick := (i%3 == 0) || i == 0
			topBorder := float32(ifThenElseInt(topThick, thick, thin))
			top := canvas.NewRectangle(color.Black)
			top.Resize(fyne.NewSize(float32(cellSize), topBorder))
			top.Move(fyne.NewPos(0, 0))
			borders = append(borders, top)

			// Left border
			leftThick := (j%3 == 0) || j == 0
			leftBorder := float32(ifThenElseInt(leftThick, thick, thin))
			left := canvas.NewRectangle(color.Black)
			left.Resize(fyne.NewSize(leftBorder, float32(cellSize)))
			left.Move(fyne.NewPos(0, 0))
			borders = append(borders, left)

			// Right border
			rightThick := (j == size-1)
			if (j+1)%3 == 0 {
				rightThick = true
			}
			rightBorder := float32(ifThenElseInt(rightThick, thick, thin))
			right := canvas.NewRectangle(color.Black)
			right.Resize(fyne.NewSize(rightBorder, float32(cellSize)))
			right.Move(fyne.NewPos(float32(cellSize)-rightBorder, 0))
			borders = append(borders, right)

			// Bottom border
			bottomThick := (i == size-1)
			if (i+1)%3 == 0 {
				bottomThick = true
			}
			bottomBorder := float32(ifThenElseInt(bottomThick, thick, thin))
			bottom := canvas.NewRectangle(color.Black)
			bottom.Resize(fyne.NewSize(float32(cellSize), bottomBorder))
			bottom.Move(fyne.NewPos(0, float32(cellSize)-bottomBorder))
			borders = append(borders, bottom)

			cell := container.NewWithoutLayout(bg)
			for _, b := range borders {
				cell.Add(b)
			}
			// Center text in cell, horizontally and vertically
			txtWidth := txt.MinSize().Width
			txtHeight := txt.MinSize().Height
			// The '+2' helps better visually center the text in most Fyne themes/fonts
			txt.Move(fyne.NewPos(
				(float32(cellSize)-txtWidth)/2,
				(float32(cellSize)-txtHeight)/2+2,
			))
			cell.Add(txt)
			grid.Add(cell)
		}
	}

	// Title with white background
	titleBg := canvas.NewRectangle(color.White)
	titleBg.Resize(fyne.NewSize(float32(size*cellSize), 56))
	title := canvas.NewText("Animated Sudoku Solver", color.Black)
	title.TextSize = 32
	title.Alignment = fyne.TextAlignCenter
	title.Move(fyne.NewPos(0, 8))
	titleBox := container.NewMax(titleBg, title)

	solveBtn := widget.NewButton("Solve (Animated)", func() {
		go func() {
			_ = solveAnimated(&board, texts, w)
		}()
	})
	resetBtn := widget.NewButton("Reset", func() {
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				board[i][j] = initialBoard[i][j]
				val := string(board[i][j])
				if val == "." {
					val = ""
				}
				texts[i][j].Text = val
				texts[i][j].Refresh()
			}
		}
	})

	w.SetContent(container.NewVBox(
		titleBox,
		grid,
		container.NewHBox(solveBtn, resetBtn),
	))
	w.Resize(fyne.NewSize(float32(size*cellSize+60), float32(size*cellSize+180)))
	w.ShowAndRun()
}

func solveAnimated(board *[9][9]byte, texts [][]*canvas.Text, w fyne.Window) bool {
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if board[row][col] == '.' {
				for c := byte('1'); c <= '9'; c++ {
					if isValid(*board, row, col, c) {
						board[row][col] = c
						texts[row][col].Text = string(c)
						texts[row][col].Refresh()
						time.Sleep(10 * time.Millisecond)
						if solveAnimated(board, texts, w) {
							return true
						}
						board[row][col] = '.'
						texts[row][col].Text = ""
						texts[row][col].Refresh()
						time.Sleep(10 * time.Millisecond)
					}
				}
				return false
			}
		}
	}
	return true
}

func isValid(board [9][9]byte, row, col int, char byte) bool {
	for i := 0; i < size; i++ {
		if board[row][i] == char || board[i][col] == char {
			return false
		}
	}
	startRow := 3 * (row / 3)
	startCol := 3 * (col / 3)
	for r := startRow; r < startRow+3; r++ {
		for c := startCol; c < startCol+3; c++ {
			if board[r][c] == char {
				return false
			}
		}
	}
	return true
}

// Helper for border thickness
func ifThenElseInt(cond bool, a, b int) int {
	if cond {
		return a
	}
	return b
}
