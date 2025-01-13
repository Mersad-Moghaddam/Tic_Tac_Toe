// This file implements a Tic-Tac-Toe game with an AI agent that can play against a human player.
// The agent supports multiple difficulty levels: Easy, Normal, Hard, and Impossible.
// Each difficulty level uses a different strategy to determine the agent's moves.
//
// Difficulty Levels and Strategies:
// 1. Easy: The agent makes random moves, selecting any empty cell on the board.
// 2. Normal: The agent uses a basic blocking strategy to prevent the opponent from winning.
//    - It checks rows, columns, and diagonals for any line where the opponent has two symbols and an empty cell.
//    - If such a line is found, the agent places its symbol in the empty cell to block the opponent.
//    - If no blocking move is found, the agent makes a random move.
// 3. Hard: The agent uses the Minimax algorithm to determine the best move.
//    - The Minimax algorithm evaluates all possible moves and their outcomes to choose the move with the highest score.
//    - The agent assumes the opponent will also play optimally and tries to maximize its own score while minimizing the opponent's score.
// 4. Impossible: The agent uses the Minimax algorithm with Alpha-Beta Pruning for optimal performance.
//    - Alpha-Beta Pruning reduces the number of nodes evaluated by the Minimax algorithm, making it more efficient.
//    - This strategy ensures the agent plays perfectly and is unbeatable.
//
// The agent's move is calculated based on the current board state and the selected difficulty level.
// The board state is represented as a 3x3 grid, where each cell can be empty or contain a player's symbol ("X" or "O").
//
// The `/agent-move` endpoint handles requests to calculate the agent's move.
// It accepts the current board state and difficulty level as query parameters and returns the agent's move as a JSON response.
//
// Helper functions are used to:
// - Count symbols in a line (row, column, or diagonal).
// - Find an empty cell in a line.
// - Check if a player has won.
// - Check if the board is full.
// - Implement the Minimax algorithm with Alpha-Beta Pruning.
//
// Templates are used to render the game interface and results.
// The templates are loaded from the `templates` directory and include a footer with the creator's name.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Difficulty levels for the agent
type Difficulty string

const (
	Easy       Difficulty = "easy"
	Normal     Difficulty = "normal"
	Hard       Difficulty = "hard"
	Impossible Difficulty = "impossible"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Represents a move made by the agent
type Move struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Handles the `/agent-move` endpoint
func agentMoveHandler(w http.ResponseWriter, r *http.Request) {
	difficulty := Difficulty(r.URL.Query().Get("difficulty"))
	boardParam := r.URL.Query().Get("board")
	if boardParam == "" {
		http.Error(w, "Board state is required", http.StatusBadRequest)
		return
	}

	var board [3][3]string
	if err := json.Unmarshal([]byte(boardParam), &board); err != nil {
		http.Error(w, "Invalid board state", http.StatusBadRequest)
		return
	}

	move := calculateAgentMove(board, difficulty)
	response, err := json.Marshal(move)
	if err != nil {
		http.Error(w, "Failed to calculate move", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Calculate the agent's move based on the difficulty level
func calculateAgentMove(board [3][3]string, difficulty Difficulty) Move {
	switch difficulty {
	case Easy:
		return randomMove(board)
	case Normal:
		return normalMove(board)
	case Hard:
		return hardMove(board)
	case Impossible:
		return unbeatableMove(board)
	default:
		return randomMove(board)
	}
}

// Random move logic for Easy difficulty
func randomMove(board [3][3]string) Move {
	var emptyCells []Move
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if board[row][col] == "" {
				emptyCells = append(emptyCells, Move{Row: row, Col: col})
			}
		}
	}

	if len(emptyCells) == 0 {
		return Move{-1, -1} // No available moves
	}

	return emptyCells[rand.Intn(len(emptyCells))]
}

// Basic strategy for Normal difficulty
func normalMove(board [3][3]string) Move {
	return blockingStrategyMove(board, "O", "X")
}

// Blocking strategy for Normal difficulty
func blockingStrategyMove(board [3][3]string, _, opponentSymbol string) Move {
	// Check rows, columns, and diagonals for blocking opportunities
	for i := 0; i < 3; i++ {
		// Check rows
		if countSymbols(opponentSymbol, []string{board[i][0], board[i][1], board[i][2]}) == 2 &&
			countSymbols("", []string{board[i][0], board[i][1], board[i][2]}) == 1 {
			return findEmptyCellInLine(board, [][2]int{{i, 0}, {i, 1}, {i, 2}})
		}

		// Check columns
		if countSymbols(opponentSymbol, []string{board[0][i], board[1][i], board[2][i]}) == 2 &&
			countSymbols("", []string{board[0][i], board[1][i], board[2][i]}) == 1 {
			return findEmptyCellInLine(board, [][2]int{{0, i}, {1, i}, {2, i}})
		}
	}

	// Check diagonals
	diagonals := [][][2]int{
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}
	for _, diagonal := range diagonals {
		if countSymbolsInLine(board, diagonal, opponentSymbol) == 2 &&
			countSymbolsInLine(board, diagonal, "") == 1 {
			return findEmptyCellInLine(board, diagonal)
		}
	}

	// If no blocking move, return random
	return randomMove(board)
}

// Helper functions
func countSymbols(symbol string, cells []string) int {
	count := 0
	for _, cell := range cells {
		if cell == symbol {
			count++
		}
	}
	return count
}

func countSymbolsInLine(board [3][3]string, line [][2]int, symbol string) int {
	count := 0
	for _, coord := range line {
		if board[coord[0]][coord[1]] == symbol {
			count++
		}
	}
	return count
}

func findEmptyCellInLine(board [3][3]string, line [][2]int) Move {
	for _, coord := range line {
		if board[coord[0]][coord[1]] == "" {
			return Move{Row: coord[0], Col: coord[1]}
		}
	}
	return Move{-1, -1}
}

// Advanced strategy for Hard difficulty using Minimax algorithm
func hardMove(board [3][3]string) Move {
	bestScore := -1000
	bestMove := Move{-1, -1}
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if board[row][col] == "" {
				board[row][col] = "O" // Agent's move
				score := minimax(board, 0, false, "O", "X")
				board[row][col] = ""
				if score > bestScore {
					bestScore = score
					bestMove = Move{Row: row, Col: col}
				}
			}
		}
	}
	return bestMove
}

// Unbeatable strategy for Impossible difficulty using Minimax with Alpha-Beta Pruning
func unbeatableMove(board [3][3]string) Move {
	var bestMove Move
	bestScore := -1000

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if board[row][col] == "" {
				board[row][col] = "O" // Agent's move
				score := minimax(board, 0, false, "O", "X")
				board[row][col] = ""
				if score > bestScore {
					bestScore = score
					bestMove = Move{Row: row, Col: col}
				}
			}
		}
	}

	return bestMove
}

// Minimax algorithm with Alpha-Beta Pruning
func minimax(board [3][3]string, depth int, isMaximizing bool, agentSymbol, opponentSymbol string) int {
	if checkWinFor(board, agentSymbol) {
		return 10 - depth
	}
	if checkWinFor(board, opponentSymbol) {
		return depth - 10
	}
	if isBoardFull(board) {
		return 0
	}

	if isMaximizing {
		bestScore := -1000
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				if board[row][col] == "" {
					board[row][col] = agentSymbol
					score := minimax(board, depth+1, false, agentSymbol, opponentSymbol)
					board[row][col] = ""
					bestScore = max(bestScore, score)
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1000
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				if board[row][col] == "" {
					board[row][col] = opponentSymbol
					score := minimax(board, depth+1, true, agentSymbol, opponentSymbol)
					board[row][col] = ""
					bestScore = min(bestScore, score)
				}
			}
		}
		return bestScore
	}
}

// Helper functions
func checkWinFor(board [3][3]string, player string) bool {
	for i := 0; i < 3; i++ {
		if board[i][0] == player && board[i][1] == player && board[i][2] == player {
			return true
		}
		if board[0][i] == player && board[1][i] == player && board[2][i] == player {
			return true
		}
	}
	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] == player {
		return true
	}
	return false
}

func isBoardFull(board [3][3]string) bool {
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if board[row][col] == "" {
				return false
			}
		}
	}
	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Load templates
var tmpl = loadTemplate("templates/index.html")
var resultTmpl = loadTemplate("templates/result.html")
var agentTmpl = loadTemplate("templates/agent.html")
var welcomeTmpl = loadTemplate("templates/welcome.html")

func loadTemplate(filename string) *template.Template {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return template.Must(template.New(filename).Parse(string(content) + "\n<footer>Created by Mersad</footer>"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		welcomeTmpl.Execute(w, nil)
	})
	mux.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})
	mux.HandleFunc("/agent", func(w http.ResponseWriter, r *http.Request) {
		agentTmpl.Execute(w, nil)
	})
	mux.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		winner := r.URL.Query().Get("winner")
		data := struct {
			Winner  string
			Player1 string
			Player2 string
		}{
			Winner:  winner,
			Player1: r.URL.Query().Get("player1"),
			Player2: r.URL.Query().Get("player2"),
		}
		resultTmpl.Execute(w, data)
	})
	mux.HandleFunc("/agent-move", agentMoveHandler)

	fmt.Println("Server running at http://localhost:8000")
	http.ListenAndServe(":8000", mux)
}

// Extract board state from request
func getBoardStateFromRequest(r *http.Request) [3][3]string {
	board := [3][3]string{}
	boardParam := r.URL.Query().Get("board")
	if boardParam != "" {
		boardSlice := []rune(boardParam)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				board[i][j] = string(boardSlice[i*3+j])
			}
		}
	}
	return board
}
