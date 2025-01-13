package main

import (
	"math/rand"
	"time"
)

// Difficulty levels
const (
	Easy = iota
	Normal
	Hard
	Impossible
)

// AI agent structure
type AIAgent struct {
	difficulty int
}

// NewAIAgent creates a new AI agent with the specified difficulty
func NewAIAgent(difficulty int) *AIAgent {
	return &AIAgent{difficulty: difficulty}
}

// GetMove returns the AI agent's move based on the difficulty level
func (agent *AIAgent) GetMove(board [3][3]string) (int, int) {
	switch agent.difficulty {
	case Easy:
		return agent.getRandomMove(board)
	case Normal:
		return agent.getNormalMove(board)
	case Hard:
		return agent.getHardMove(board)
	case Impossible:
		return agent.getImpossibleMove(board)
	default:
		return agent.getRandomMove(board)
	}
}

// getRandomMove returns a random valid move
func (agent *AIAgent) getRandomMove(board [3][3]string) (int, int) {
	rand.Seed(time.Now().UnixNano())
	for {
		row := rand.Intn(3)
		col := rand.Intn(3)
		if board[row][col] == "" {
			return row, col
		}
	}
}

// getNormalMove returns a move with basic strategy
func (agent *AIAgent) getNormalMove(board [3][3]string) (int, int) {
	// Try to win or block opponent from winning
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				board[i][j] = "O"
				if checkWin(board, "O") {
					return i, j
				}
				board[i][j] = "X"
				if checkWin(board, "X") {
					board[i][j] = ""
					return i, j
				}
				board[i][j] = ""
			}
		}
	}
	return agent.getRandomMove(board)
}

// getHardMove returns a move with advanced strategy
func (agent *AIAgent) getHardMove(board [3][3]string) (int, int) {
	// Implement advanced strategy here
	// Prioritize center, corners, and sides
	if board[1][1] == "" {
		return 1, 1
	}
	moves := [][2]int{{0, 0}, {0, 2}, {2, 0}, {2, 2}, {0, 1}, {1, 0}, {1, 2}, {2, 1}}
	for _, move := range moves {
		if board[move[0]][move[1]] == "" {
			return move[0], move[1]
		}
	}
	return agent.getRandomMove(board)
}

// getImpossibleMove returns the best possible move using minimax algorithm
func (agent *AIAgent) getImpossibleMove(board [3][3]string) (int, int) {
	bestScore := -1000
	var move [2]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				board[i][j] = "O"
				score := minimax(board, 0, false)
				board[i][j] = ""
				if score > bestScore {
					bestScore = score
					move = [2]int{i, j}
				}
			}
		}
	}
	return move[0], move[1]
}

// minimax algorithm to evaluate the best move
func minimax(board [3][3]string, depth int, isMaximizing bool) int {
	if checkWin(board, "O") {
		return 10 - depth
	}
	if checkWin(board, "X") {
		return depth - 10
	}
	if isBoardFull(board) {
		return 0
	}

	if isMaximizing {
		bestScore := -1000
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == "" {
					board[i][j] = "O"
					score := minimax(board, depth+1, false)
					board[i][j] = ""
					if score > bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	} else {
		bestScore := 1000
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == "" {
					board[i][j] = "X"
					score := minimax(board, depth+1, true)
					board[i][j] = ""
					if score < bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	}
}

// checkWin checks if a player has won
func checkWin(board [3][3]string, player string) bool {
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

// isBoardFull checks if the board is full
func isBoardFull(board [3][3]string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

// ...existing code...
