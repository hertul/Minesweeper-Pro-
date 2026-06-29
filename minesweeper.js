// minesweeper.js
#!/usr/bin/env node
'use strict';

const fs = require('fs');
const readline = require('readline');
const os = require('os');

const COLORS = {
    reset: '\x1b[0m',
    red: '\x1b[91m',
    green: '\x1b[92m',
    yellow: '\x1b[93m',
    blue: '\x1b[94m',
    magenta: '\x1b[95m',
    cyan: '\x1b[96m',
    white: '\x1b[97m',
    gray: '\x1b[90m',
    bold: '\x1b[1m'
};

function colorize(text, color) {
    return COLORS[color] + text + COLORS.reset;
}

class Minesweeper {
    constructor(rows, cols, mines) {
        this.rows = rows;
        this.cols = cols;
        this.minesTotal = mines;
        this.board = Array.from({ length: rows }, () => Array(cols).fill(' '));
        this.real = Array.from({ length: rows }, () => Array(cols).fill(0));
        this.revealed = Array.from({ length: rows }, () => Array(cols).fill(false));
        this.flags = Array.from({ length: rows }, () => Array(cols).fill(false));
        this.gameOver = false;
        this.won = false;
        this.firstMove = true;
        this.minesRemaining = mines;
        this.startTime = null;
        this.moves = 0;
    }

    generateBoard(firstRow, firstCol) {
        const mines = new Set();
        while (mines.size < this.minesTotal) {
            const r = Math.floor(Math.random() * this.rows);
            const c = Math.floor(Math.random() * this.cols);
            if (r === firstRow && c === firstCol) continue;
            if (Math.abs(r - firstRow) <= 1 && Math.abs(c - firstCol) <= 1) continue;
            const key = `${r},${c}`;
            if (!mines.has(key)) {
                mines.add(key);
                this.real[r][c] = -1;
            }
        }
        // Подсчёт чисел
        for (let r = 0; r < this.rows; r++) {
            for (let c = 0; c < this.cols; c++) {
                if (this.real[r][c] === -1) continue;
                let count = 0;
                for (let dr = -1; dr <= 1; dr++) {
                    for (let dc = -1; dc <= 1; dc++) {
                        const nr = r + dr, nc = c + dc;
                        if (nr >= 0 && nr < this.rows && nc >= 0 && nc < this.cols && this.real[nr][nc] === -1)
                            count++;
                    }
                }
                this.real[r][c] = count;
            }
        }
    }

    render() {
        console.clear();
        let output = '';
        output += colorize('   ', 'bold');
        for (let c = 0; c < this.cols; c++) {
            output += colorize(`${String(c+1).padStart(2)} `, 'bold');
        }
        output += '\n';
        output += colorize('   +', 'gray');
        for (let c = 0; c < this.cols; c++) output += '---';
        output += colorize('+\n', 'gray');
        for (let r = 0; r < this.rows; r++) {
            output += colorize(`${String(r+1).padStart(2)} `, 'bold');
            output += colorize('|', 'gray');
            for (let c = 0; c < this.cols; c++) {
                if (this.gameOver && this.real[r][c] === -1) {
                    output += colorize(' 💣', 'red');
                } else if (this.flags[r][c]) {
                    output += colorize(' ⚑', 'yellow');
                } else if (this.revealed[r][c]) {
                    const val = this.real[r][c];
                    if (val === -1) output += colorize(' 💣', 'red');
                    else if (val === 0) output += '  ';
                    else {
                        const colors = ['', 'blue', 'green', 'red', 'magenta', 'cyan', 'yellow', 'white', 'gray'];
                        output += colorize(` ${val}`, colors[val]);
                    }
                } else {
                    output += ' ■';
                }
            }
            output += colorize('|\n', 'gray');
        }
        output += colorize('   +', 'gray');
        for (let c = 0; c < this.cols; c++) output += '---';
        output += colorize('+\n', 'gray');
        output += colorize(`Мин осталось: ${this.minesRemaining}\n`, 'yellow');
        if (!this.gameOver && this.startTime) {
            const elapsed = Math.floor((Date.now() - this.startTime) / 1000);
            output += colorize(`Время: ${elapsed} сек\n`, 'blue');
        }
        if (this.gameOver) {
            output += this.won ? colorize('🎉 ПОБЕДА!\n', 'green') : colorize('💥 ПОРАЖЕНИЕ!\n', 'red');
        }
        console.log(output);
    }

    openCell(r, c) {
        if (this.gameOver || this.revealed[r][c] || this.flags[r][c]) return;
        if (this.firstMove) {
            this.generateBoard(r, c);
            this.firstMove = false;
            this.startTime = Date.now();
        }
        this.revealed[r][c] = true;
        this.moves++;
        if (this.real[r][c] === -1) {
            this.gameOver = true;
            return;
        }
        if (this.real[r][c] === 0) {
            for (let dr = -1; dr <= 1; dr++) {
                for (let dc = -1; dc <= 1; dc++) {
                    const nr = r + dr, nc = c + dc;
                    if (nr >= 0 && nr < this.rows && nc >= 0 && nc < this.cols && !this.revealed[nr][nc] && !this.flags[nr][nc]) {
                        this.openCell(nr, nc);
                    }
                }
            }
        }
        // Проверка победы
        let revealedCount = 0;
        for (let i = 0; i < this.rows; i++) {
            for (let j = 0; j < this.cols; j++) {
                if (this.revealed[i][j]) revealedCount++;
            }
        }
        if (revealedCount === this.rows * this.cols - this.minesTotal) {
            this.won = true;
            this.gameOver = true;
        }
    }

    toggleFlag(r, c) {
        if (this.gameOver || this.revealed[r][c]) return;
        if (this.flags[r][c]) {
            this.flags[r][c] = false;
            this.minesRemaining++;
        } else {
            this.flags[r][c] = true;
            this.minesRemaining--;
        }
    }
}

function main() {
    const args = process.argv.slice(2);
    let rows = 9, cols = 9, mines = 10;
    if (args.length > 0) {
        const level = args[0].toLowerCase();
        if (level === 'easy') { rows = 9; cols = 9; mines = 10; }
        else if (level === 'medium') { rows = 16; cols = 16; mines = 40; }
        else if (level === 'hard') { rows = 30; cols = 16; mines = 99; }
        else {
            const v = parseInt(level);
            if (!isNaN(v)) {
                rows = cols = v;
                if (rows < 3) rows = cols = 3;
                if (args.length > 1) {
                    const m = parseInt(args[1]);
                    if (!isNaN(m)) mines = m;
                } else {
                    mines = Math.max(1, Math.floor(rows * cols / 10));
                }
            }
        }
    }
    const game = new Minesweeper(rows, cols, mines);
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });

    function prompt() {
        game.render();
        rl.question('Введите команду (строка столбец для открытия, f строка столбец для флага, q для выхода): ', (cmd) => {
            if (cmd.trim().toLowerCase() === 'q') {
                console.log('Выход.');
                rl.close();
                process.exit(0);
            }
            const parts = cmd.trim().split(/\s+/);
            if (parts.length < 2) {
                console.log('Неверный формат.');
                prompt();
                return;
            }
            if (parts[0].toLowerCase() === 'f') {
                if (parts.length !== 3) {
                    console.log('Неверный формат. Используйте: f строка столбец');
                    prompt();
                    return;
                }
                const r = parseInt(parts[1]) - 1;
                const c = parseInt(parts[2]) - 1;
                if (r >= 0 && r < rows && c >= 0 && c < cols) {
                    game.toggleFlag(r, c);
                } else {
                    console.log('Координаты вне поля.');
                }
            } else {
                if (parts.length !== 2) {
                    console.log('Неверный формат. Используйте: строка столбец');
                    prompt();
                    return;
                }
                const r = parseInt(parts[0]) - 1;
                const c = parseInt(parts[1]) - 1;
                if (r >= 0 && r < rows && c >= 0 && c < cols) {
                    game.openCell(r, c);
                } else {
                    console.log('Координаты вне поля.');
                }
            }
            if (!game.gameOver) {
                prompt();
            } else {
                game.render();
                if (game.won) {
                    const elapsed = Math.floor((Date.now() - game.startTime) / 1000);
                    console.log(`Поздравляем! Вы выиграли за ${game.moves} ходов и ${elapsed} секунд.`);
                } else {
                    console.log('Игра окончена. Попробуйте снова!');
                }
                rl.close();
            }
        });
    }
    prompt();
}

main();
