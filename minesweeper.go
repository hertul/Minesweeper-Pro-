// minesweeper.go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	reset  = "\033[0m"
	red    = "\033[91m"
	green  = "\033[92m"
	yellow = "\033[93m"
	blue   = "\033[94m"
	magenta= "\033[95m"
	cyan   = "\033[96m"
	white  = "\033[97m"
	gray   = "\033[90m"
	bold   = "\033[1m"
)

func colorize(text, color string) string {
	return color + text + reset
}

type Minesweeper struct {
	rows, cols, minesTotal int
	board                   [][]string   // видимое поле
	real                    [][]int      // -1 = мина, 0-8 = число
	revealed                [][]bool
	flags                   [][]bool
	gameOver                bool
	won                     bool
	firstMove               bool
	minesRemaining          int
	startTime               time.Time
	moves                   int
}

func NewMinesweeper(rows, cols, mines int) *Minesweeper {
	m := &Minesweeper{
		rows:           rows,
		cols:           cols,
		minesTotal:     mines,
		minesRemaining: mines,
		firstMove:      true,
	}
	m.board = make([][]string, rows)
	m.real = make([][]int, rows)
	m.revealed = make([][]bool, rows)
	m.flags = make([][]bool, rows)
	for i := 0; i < rows; i++ {
		m.board[i] = make([]string, cols)
		m.real[i] = make([]int, cols)
		m.revealed[i] = make([]bool, cols)
		m.flags[i] = make([]bool, cols)
	}
	return m
}

func (m *Minesweeper) generateBoard(firstRow, firstCol int) {
	mines := make(map[[2]int]bool)
	for len(mines) < m.minesTotal {
		r := rand.Intn(m.rows)
		c := rand.Intn(m.cols)
		if r == firstRow && c == firstCol {
			continue
		}
		if abs(r-firstRow) <= 1 && abs(c-firstCol) <= 1 {
			continue
		}
		key := [2]int{r, c}
		if !mines[key] {
			mines[key] = true
			m.real[r][c] = -1
		}
	}
	// Подсчёт чисел
	for r := 0; r < m.rows; r++ {
		for c := 0; c < m.cols; c++ {
			if m.real[r][c] == -1 {
				continue
			}
			count := 0
			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					nr, nc := r+dr, c+dc
					if nr >= 0 && nr < m.rows && nc >= 0 && nc < m.cols && m.real[nr][nc] == -1 {
						count++
					}
				}
			}
			m.real[r][c] = count
		}
	}
}

func (m *Minesweeper) render() {
	fmt.Print("\033[H\033[2J")
	fmt.Print(colorize("   ", bold))
	for c := 0; c < m.cols; c++ {
		fmt.Print(colorize(fmt.Sprintf("%2d ", c+1), bold))
	}
	fmt.Println()
	fmt.Print(colorize("   "+"+", gray))
	for c := 0; c < m.cols; c++ {
		fmt.Print("---")
	}
	fmt.Println(colorize("+", gray))
	for r := 0; r < m.rows; r++ {
		fmt.Print(colorize(fmt.Sprintf("%2d ", r+1), bold))
		fmt.Print(colorize("|", gray))
		for c := 0; c < m.cols; c++ {
			if m.gameOver && m.real[r][c] == -1 {
				fmt.Print(colorize(" 💣", red))
			} else if m.flags[r][c] {
				fmt.Print(colorize(" ⚑", yellow))
			} else if m.revealed[r][c] {
				val := m.real[r][c]
				if val == -1 {
					fmt.Print(colorize(" 💣", red))
				} else if val == 0 {
					fmt.Print("  ")
				} else {
					colors := []string{"", blue, green, red, magenta, cyan, yellow, white, gray}
					fmt.Print(colorize(fmt.Sprintf(" %d", val), colors[val]))
				}
			} else {
				fmt.Print(" ■")
			}
		}
		fmt.Println(colorize("|", gray))
	}
	fmt.Print(colorize("   "+"+", gray))
	for c := 0; c < m.cols; c++ {
		fmt.Print("---")
	}
	fmt.Println(colorize("+", gray))
	fmt.Println(colorize(fmt.Sprintf("Мин осталось: %d", m.minesRemaining), yellow))
	if !m.gameOver && !m.startTime.IsZero() {
		elapsed := int(time.Since(m.startTime).Seconds())
		fmt.Println(colorize(fmt.Sprintf("Время: %d сек", elapsed), blue))
	}
	if m.gameOver {
		if m.won {
			fmt.Println(colorize("🎉 ПОБЕДА!", green))
		} else {
			fmt.Println(colorize("💥 ПОРАЖЕНИЕ!", red))
		}
	}
}

func (m *Minesweeper) openCell(r, c int) {
	if m.gameOver || m.revealed[r][c] || m.flags[r][c] {
		return
	}
	if m.firstMove {
		m.generateBoard(r, c)
		m.firstMove = false
		m.startTime = time.Now()
	}
	m.revealed[r][c] = true
	m.moves++

	if m.real[r][c] == -1 {
		m.gameOver = true
		return
	}
	if m.real[r][c] == 0 {
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				nr, nc := r+dr, c+dc
				if nr >= 0 && nr < m.rows && nc >= 0 && nc < m.cols && !m.revealed[nr][nc] && !m.flags[nr][nc] {
					m.openCell(nr, nc)
				}
			}
		}
	}
	// Проверка победы
	revealedCount := 0
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			if m.revealed[i][j] {
				revealedCount++
			}
		}
	}
	if revealedCount == m.rows*m.cols-m.minesTotal {
		m.won = true
		m.gameOver = true
	}
}

func (m *Minesweeper) toggleFlag(r, c int) {
	if m.gameOver || m.revealed[r][c] {
		return
	}
	if m.flags[r][c] {
		m.flags[r][c] = false
		m.minesRemaining++
	} else {
		m.flags[r][c] = true
		m.minesRemaining--
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var rows, cols, mines int = 9, 9, 10
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "easy":
			rows, cols, mines = 9, 9, 10
		case "medium":
			rows, cols, mines = 16, 16, 40
		case "hard":
			rows, cols, mines = 30, 16, 99
		default:
			if v, err := strconv.Atoi(arg); err == nil {
				rows, cols = v, v
				if rows < 3 {
					rows, cols = 3, 3
				}
				if len(os.Args) > 2 {
					if m, err := strconv.Atoi(os.Args[2]); err == nil {
						mines = m
					}
				} else {
					mines = max(1, rows*cols/10)
				}
			}
		}
	}
	if len(os.Args) > 2 && os.Args[1] != "easy" && os.Args[1] != "medium" && os.Args[1] != "hard" {
		if m, err := strconv.Atoi(os.Args[2]); err == nil {
			mines = m
		}
	}

	game := NewMinesweeper(rows, cols, mines)
	scanner := bufio.NewScanner(os.Stdin)

	for !game.gameOver {
		game.render()
		fmt.Print("Введите команду (строка столбец для открытия, f строка столбец для флага, q для выхода): ")
		if !scanner.Scan() {
			break
		}
		cmd := strings.TrimSpace(scanner.Text())
		if cmd == "q" {
			fmt.Println("Выход.")
			return
		}
		parts := strings.Fields(cmd)
		if len(parts) < 2 {
			fmt.Println("Неверный формат.")
			continue
		}
		if parts[0] == "f" {
			if len(parts) != 3 {
				fmt.Println("Неверный формат. Используйте: f строка столбец")
				continue
			}
			r, err1 := strconv.Atoi(parts[1])
			c, err2 := strconv.Atoi(parts[2])
			if err1 == nil && err2 == nil && r >= 1 && r <= rows && c >= 1 && c <= cols {
				game.toggleFlag(r-1, c-1)
			} else {
				fmt.Println("Координаты вне поля.")
			}
		} else {
			if len(parts) != 2 {
				fmt.Println("Неверный формат. Используйте: строка столбец")
				continue
			}
			r, err1 := strconv.Atoi(parts[0])
			c, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil && r >= 1 && r <= rows && c >= 1 && c <= cols {
				game.openCell(r-1, c-1)
			} else {
				fmt.Println("Координаты вне поля.")
			}
		}
	}
	game.render()
	if game.won {
		elapsed := int(time.Since(game.startTime).Seconds())
		fmt.Printf("Поздравляем! Вы выиграли за %d ходов и %d секунд.\n", game.moves, elapsed)
	} else {
		fmt.Println("Игра окончена. Попробуйте снова!")
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
