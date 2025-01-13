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

type Difficulty string

const (
	Easy       Difficulty = "easy"
	Normal     Difficulty = "normal"
	Hard       Difficulty = "hard"
	Impossible Difficulty = "impossible"
)

var currentDifficulty Difficulty

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
		row, col := normalMove(board)
		return Move{Row: row, Col: col}
	case Hard:
		row, col := hardMove(board)
		return Move{Row: row, Col: col}
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
func normalMove(board [3][3]string) (int, int) {
	// Implement a basic strategy for normal difficulty
	// For simplicity, this example uses randomMove
	return randomMove(board).Row, randomMove(board).Col
}

// Advanced strategy for Hard difficulty
func hardMove(board [3][3]string) (int, int) {
	// Implement a more advanced strategy for hard difficulty
	// For simplicity, this example uses randomMove
	return randomMove(board).Row, randomMove(board).Col
}

// Unbeatable strategy for Impossible difficulty using Minimax with Alpha-Beta Pruning
func unbeatableMove(board [3][3]string) Move {
	var bestMove Move
	bestScore := -1000

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if board[row][col] == "" {
				board[row][col] = "O" // Agent's move
				score := minimax(board, 0, false, -1000, 1000)
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
func minimax(board [3][3]string, depth int, isMaximizing bool, alpha, beta int) int {
	if checkWinner(board, "O") {
		return 10 - depth
	}
	if checkWinner(board, "X") {
		return depth - 10
	}
	if isBoardFull(board) {
		return 0
	}

	if isMaximizing {
		maxEval := -1000
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				if board[row][col] == "" {
					board[row][col] = "O"
					eval := minimax(board, depth+1, false, alpha, beta)
					board[row][col] = ""
					maxEval = max(maxEval, eval)
					alpha = max(alpha, eval)
					if beta <= alpha {
						break
					}
				}
			}
		}
		return maxEval
	} else {
		minEval := 1000
		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				if board[row][col] == "" {
					board[row][col] = "X"
					eval := minimax(board, depth+1, true, alpha, beta)
					board[row][col] = ""
					minEval = min(minEval, eval)
					beta = min(beta, eval)
					if beta <= alpha {
						break
					}
				}
			}
		}
		return minEval
	}
}

// Helper functions
func checkWinner(board [3][3]string, player string) bool {
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

var tmpl = loadTemplate("templates/index.html")
var resultTmpl = loadTemplate("templates/result.html")
var agentTmpl = loadTemplate("templates/agent.html")
var welcomeTmpl = loadTemplate("templates/welcome.html")

func loadTemplate(filename string) *template.Template {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return template.Must(template.New(filename).Parse(string(content)))
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
