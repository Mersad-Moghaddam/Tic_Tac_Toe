package main

import (
	"fmt"
	"html/template"
	"net/http"
)

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
        let aiAgent = null;

        playerForm.addEventListener('submit', function(event) {
            event.preventDefault();
            player1 = document.getElementById('player1').value;
            player2 = document.getElementById('player2').value;
            const firstPlayer = document.getElementById('firstPlayer').value;
            currentPlayer = firstPlayer === 'player1' ? player1 : player2;
            playerForm.style.display = 'none';
            gameBoard.style.display = 'grid';
        });

        cells.forEach(cell => {
            cell.addEventListener('click', handleClick, { once: true });
        });

        function handleClick(e) {
            const cell = e.target;
            cell.textContent = currentPlayer === player1 ? 'X' : 'O';
            cell.classList.add(currentPlayer === player1 ? 'x' : 'o');
            if (checkWin()) {
                setTimeout(() => {
                    alert(currentPlayer + ' wins!');
                    window.location.href = '/result?winner=' + currentPlayer + '&player1=' + player1 + '&player2=' + player2;
                }, 100);
            } else {
                currentPlayer = currentPlayer === player1 ? player2 : player1;
                if (aiAgent && currentPlayer === player2) {
                    const [row, col] = aiAgent.GetMove(getBoardState());
                    const aiCell = cells[row * 3 + col];
                    aiCell.click();
                }
            }
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

        // Initialize AI agent if difficulty is set
        const urlParams = new URLSearchParams(window.location.search);
        const difficulty = urlParams.get('difficulty');
        if (difficulty) {
            aiAgent = new AIAgent(parseInt(difficulty));
            player2 = 'AI';
            playerForm.style.display = 'none';
            gameBoard.style.display = 'grid';
            currentPlayer = player1;
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
        <p>Winner: {{.Winner}}</p>
        <p>Player 1: {{.Player1}}</p>
        <p>Player 2: {{.Player2}}</p>
        <button onclick="window.location.href='/'">Play Again</button>
        <footer>Created by Mersad</footer>
    </div>
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
            const difficulty = document.getElementById('difficulty').value;
            window.location.href = '/game?difficulty=' + difficulty;
        });
    </script>
</body>
</html>
`))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		welcomeTmpl.Execute(w, nil)
	})
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})
	http.HandleFunc("/agent", func(w http.ResponseWriter, r *http.Request) {
		agentTmpl.Execute(w, nil)
	})
	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("Server running at http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
