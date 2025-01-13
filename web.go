package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
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

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tic Tac Toe</title>
    <style>
        body {
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background-color: #f0f0f0;
            font-family: Arial, sans-serif;
            margin: 0;
        }
        .container {
            text-align: center;
            background: #fff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .game-board {
            display: grid;
            grid-template-columns: repeat(3, 100px);
            grid-template-rows: repeat(3, 100px);
            gap: 5px;
            margin-top: 20px;
        }
        .cell {
            width: 100px;
            height: 100px;
            background-color: #fff;
            border: 1px solid #ccc;
            display: flex;
            justify-content: center;
            align-items: center;
            font-size: 24px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        .cell:hover {
            background-color: #e0e0e0;
        }
        .cell.x {
            color: red;
        }
        .cell.o {
            color: green;
        }
        .form-group {
            margin-bottom: 15px;
        }
        button {
            padding: 10px 20px;
            background-color: #007BFF;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #0056b3;
        }
        footer {
            margin-top: 20px;
            font-size: 14px;
            color: #888;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Tic Tac Toe</h1>
        <form id="playerForm">
            <div class="form-group">
                <label for="player1">Player 1 Name:</label>
                <input type="text" id="player1" name="player1" required>
            </div>
            <div class="form-group">
                <label for="player2">Player 2 Name:</label>
                <input type="text" id="player2" name="player2" required>
            </div>
            <div class="form-group">
                <label for="firstPlayer">Who plays first:</label>
                <select id="firstPlayer" name="firstPlayer" required>
                    <option value="player1">Player 1</option>
                    <option value="player2">Player 2</option>
                </select>
            </div>
            <button type="submit">Start Game</button>
        </form>
        <div class="game-board" id="gameBoard" style="display: none;">
            <div class="cell" data-cell></div>
            <div class="cell" data-cell></div>
            <div class="cell" data-cell></div>
            <div class="cell" data-cell></div>
            <div class="cell" data-cell></div>
            <div class="cell" data-cell></div>
            <div class="cell" data-cell></div>
            <div class="cell" data-cell></div>
            <div class="cell" data-cell></div>
        </div>
        <footer>Created by Mersad</footer>
    </div>
    <script>
        const playerForm = document.getElementById('playerForm');
        const gameBoard = document.getElementById('gameBoard');
        const cells = document.querySelectorAll('[data-cell]');
        let currentPlayer;
        let player1, player2;

        const urlParams = new URLSearchParams(window.location.search);
        const difficulty = urlParams.get('difficulty');
        const playerName = urlParams.get('playerName');
        const firstPlayer = urlParams.get('firstPlayer');
        let isAgentGame = difficulty !== null;

        if (isAgentGame) {
            player1 = playerName || "Player";
            player2 = "Agent";
            currentPlayer = firstPlayer === 'player' ? player1 : player2;
            playerForm.style.display = 'none';
            gameBoard.style.display = 'grid';
            if (currentPlayer === player2) {
                setTimeout(agentMove, 500);
            }
        } else {
            playerForm.addEventListener('submit', function(event) {
                event.preventDefault();
                player1 = document.getElementById('player1').value;
                player2 = document.getElementById('player2').value;
                const firstPlayer = document.getElementById('firstPlayer').value;
                currentPlayer = firstPlayer === 'player1' ? player1 : player2;
                playerForm.style.display = 'none';
                gameBoard.style.display = 'grid';
            });
        }

        cells.forEach(cell => {
            cell.addEventListener('click', handleClick, { once: true });
        });

        function handleClick(e) {
            const cell = e.target;
            if (cell.textContent !== '') return; // Prevent choosing an already chosen cell
            cell.textContent = currentPlayer === player1 ? 'X' : 'O';
            cell.classList.add(currentPlayer === player1 ? 'x' : 'o');
            if (checkWin()) {
                setTimeout(() => {
                    alert(currentPlayer + ' wins!');
                    window.location.href = '/result?winner=' + currentPlayer + '&player1=' + player1 + '&player2=' + player2;
                }, 100);
            } else if (isBoardFull()) {
                setTimeout(() => {
                    alert('It\'s a tie!');
                    window.location.href = '/result?winner=Tie&player1=' + player1 + '&player2=' + player2;
                }, 100);
            } else {
                currentPlayer = currentPlayer === player1 ? player2 : player1;
                if (isAgentGame && currentPlayer === player2) {
                    setTimeout(agentMove, 500);
                }
            }
        }

        function agentMove() {
            fetch('/agent-move?difficulty=' + difficulty + '&board=' + JSON.stringify(getBoardState()))
                .then(response => response.json())
                .then(data => {
                    const cellIndex = data.row * 3 + data.col;
                    const cell = cells[cellIndex];
                    cell.textContent = 'O';
                    cell.classList.add('o');
                    if (checkWin()) {
                        setTimeout(() => {
                            alert(player2 + ' wins!');
                            window.location.href = '/result?winner=' + player2 + '&player1=' + player1 + '&player2=' + player2;
                        }, 100);
                    } else if (isBoardFull()) {
                        setTimeout(() => {
                            alert('It\'s a tie!');
                            window.location.href = '/result?winner=Tie&player1=' + player1 + '&player2=' + player2;
                        }, 100);
                    } else {
                        currentPlayer = player1;
                    }
                });
        }

        function getBoardState() {
            const board = [['', '', ''], ['', '', ''], ['', '', '']];
            cells.forEach((cell, index) => {
                const row = Math.floor(index / 3);
                const col = index % 3;
                board[row][col] = cell.textContent;
            });
            return board;
        }

        function checkWin() {
            const winPatterns = [
                [0, 1, 2], [3, 4, 5], [6, 7, 8], // rows
                [0, 3, 6], [1, 4, 7], [2, 5, 8], // columns
                [0, 4, 8], [2, 4, 6] // diagonals
            ];
            return winPatterns.some(pattern => {
                return pattern.every(index => {
                    return cells[index].textContent === (currentPlayer === player1 ? 'X' : 'O');
                });
            });
        }

        function isBoardFull() {
            return [...cells].every(cell => cell.textContent !== '');
        }
    </script>
</body>
</html>
`))

var resultTmpl = template.Must(template.New("result").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Game Result</title>
    <style>
        body {
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background-color: #f0f0f0;
            font-family: Arial, sans-serif;
            margin: 0;
        }
        .container {
            text-align: center;
            background: #fff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        button {
            padding: 10px 20px;
            background-color: #007BFF;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #0056b3;
        }
        footer {
            margin-top: 20px;
            font-size: 14px;
            color: #888;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Game Result</h1>
        {{if eq .Winner "Tie"}}
        <p>It's a tie!</p>
        {{else}}
        <p>Winner: {{.Winner}}</p>
        {{end}}
        <p>Player 1: {{.Player1}}</p>
        <p>Player 2: {{.Player2}}</p>
        <button onclick="window.location.href='/'">Play Again</button>
        <footer>Created by Mersad</footer>
    </div>
</body>
</html>
`))

var agentTmpl = template.Must(template.New("agent").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose Difficulty</title>
    <style>
        body {
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background-color: #f0f0f0;
            font-family: Arial, sans-serif;
            margin: 0;
        }
        .container {
            text-align: center;
            background: #fff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .form-group {
            margin-bottom: 15px;
        }
        button {
            padding: 10px 20px;
            background-color: #007BFF;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #0056b3;
        }
        footer {
            margin-top: 20px;
            font-size: 14px;
            color: #888;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Choose Difficulty</h1>
        <form id="difficultyForm">
            <div class="form-group">
                <label for="playerName">Your Name:</label>
                <input type="text" id="playerName" name="playerName" required>
            </div>
            <div class="form-group">
                <label for="firstPlayer">Who plays first:</label>
                <select id="firstPlayer" name="firstPlayer" required>
                    <option value="player">You</option>
                    <option value="agent">Agent</option>
                </select>
            </div>
            <div class="form-group">
                <label for="difficulty">Select Difficulty:</label>
                <select id="difficulty" name="difficulty" required>
                    <option value="easy">Easy</option>
                    <option value="normal">Normal</option>
                    <option value="hard">Hard</option>
                    <option value="impossible">Impossible</option>
                </select>
            </div>
            <button type="submit">Let's Go</button>
        </form>
        <footer>Created by Mersad</footer>
    </div>
    <script>
        const difficultyForm = document.getElementById('difficultyForm');
        difficultyForm.addEventListener('submit', function(event) {
            event.preventDefault();
            const playerName = document.getElementById('playerName').value;
            const firstPlayer = document.getElementById('firstPlayer').value;
            const difficulty = document.getElementById('difficulty').value;
            window.location.href = '/game?difficulty=' + difficulty + '&playerName=' + playerName + '&firstPlayer=' + firstPlayer;
        });
    </script>
</body>
</html>
`))

var welcomeTmpl = template.Must(template.New("welcome").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Tic Tac Toe</title>
    <style>
        body {
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            height: 100vh;
            background-color: #f0f0f0;
            font-family: Arial, sans-serif;
            margin: 0;
            animation: fadeIn 2s;
        }
        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }
        .container {
            text-align: center;
            background: #fff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            animation: slideIn 1s;
        }
        @keyframes slideIn {
            from { transform: translateY(-50px); }
            to { transform: translateY(0); }
        }
        button {
            padding: 10px 20px;
            background-color: #007BFF;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        button:hover {
            background-color: #0056b3;
        }
        .shiny-button {
            background: linear-gradient(45deg, #ff0066, #ffcc00);
            background-size: 200% 200%;
            animation: shiny 2s linear infinite;
            color: #fff;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            padding: 10px 20px;
            transition: background-color 0.3s;
        }
        @keyframes shiny {
            0% { background-position: 0% 50%; }
            50% { background-position: 100% 50%; }
            100% { background-position: 0% 50%; }
        }
        footer {
            margin-top: 20px;
            font-size: 14px;
            color: #888;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome to Tic Tac Toe</h1>
        <p>This project is a simple web-based Tic Tac Toe game built with Go and HTML/CSS.</p>
        <button onclick="window.location.href='/game'">Start Game</button>
        <button class="shiny-button" onclick="window.location.href='/agent'">Play with Agent</button>
        <footer>Created by Mersad</footer>
    </div>
</body>
</html>
`))

func aiMove(board [3][3]string, difficulty Difficulty) (int, int) {
	switch difficulty {
	case Easy:
		return randomMove(board)
	case Normal:
		return normalMove(board)
	case Hard:
		return hardMove(board)
	case Impossible:
		return impossibleMove(board)
	default:
		return randomMove(board)
	}
}

func randomMove(board [3][3]string) (int, int) {
	for {
		row := rand.Intn(3)
		col := rand.Intn(3)
		if board[row][col] == "" {
			return row, col
		}
	}
}

func normalMove(board [3][3]string) (int, int) {
	// Implement a basic strategy for normal difficulty
	// For simplicity, this example uses randomMove
	return randomMove(board)
}

func hardMove(board [3][3]string) (int, int) {
	// Implement a more advanced strategy for hard difficulty
	// For simplicity, this example uses randomMove
	return randomMove(board)
}

func impossibleMove(board [3][3]string) (int, int) {
	// Implement an unbeatable strategy for impossible difficulty
	// For simplicity, this example uses a basic minimax algorithm
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

func minimax(board [3][3]string, depth int, isMaximizing bool) int {
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

func checkWinner(board [3][3]string, player string) bool {
	winPatterns := [][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}}, {{1, 0}, {1, 1}, {1, 2}}, {{2, 0}, {2, 1}, {2, 2}}, // rows
		{{0, 0}, {1, 0}, {2, 0}}, {{0, 1}, {1, 1}, {2, 1}}, {{0, 2}, {1, 2}, {2, 2}}, // columns
		{{0, 0}, {1, 1}, {2, 2}}, {{0, 2}, {1, 1}, {2, 0}}, // diagonals
	}
	for _, pattern := range winPatterns {
		if board[pattern[0][0]][pattern[0][1]] == player &&
			board[pattern[1][0]][pattern[1][1]] == player &&
			board[pattern[2][0]][pattern[2][1]] == player {
			return true
		}
	}
	return false
}

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
	mux.HandleFunc("/agent-move", func(w http.ResponseWriter, r *http.Request) {
		board := getBoardStateFromRequest(r)
		difficulty := Difficulty(r.URL.Query().Get("difficulty"))
		row, col := aiMove(board, difficulty)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"row": %d, "col": %d}`, row, col)
	})

	fmt.Println("Server running at http://localhost:8000")
	http.ListenAndServe(":8000", mux)
}

func getBoardStateFromRequest(r *http.Request) [3][3]string {
	board := [3][3]string{}
	boardParam := r.URL.Query().Get("board")
	if boardParam != "" {
		for i, row := range board {
			for j := range row {
				board[i][j] = string(boardParam[i*3+j])
			}
		}
	}
	return board
}
