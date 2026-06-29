// minesweeper.cpp
#include <iostream>
#include <vector>
#include <string>
#include <random>
#include <chrono>
#include <thread>
#include <cstdlib>
#include <ctime>
#include <algorithm>
#include <termios.h>
#include <unistd.h>
#include <fcntl.h>

using namespace std;

const string RESET = "\033[0m";
const string RED = "\033[91m";
const string GREEN = "\033[92m";
const string YELLOW = "\033[93m";
const string BLUE = "\033[94m";
const string MAGENTA = "\033[95m";
const string CYAN = "\033[96m";
const string WHITE = "\033[97m";
const string GRAY = "\033[90m";
const string BOLD = "\033[1m";

string colorize(const string& text, const string& color) {
    return color + text + RESET;
}

class Minesweeper {
public:
    Minesweeper(int rows, int cols, int mines) 
        : rows(rows), cols(cols), minesTotal(mines), minesRemaining(mines), firstMove(true) {
        board.resize(rows, vector<char>(cols, ' '));
        real.resize(rows, vector<int>(cols, 0));
        revealed.resize(rows, vector<bool>(cols, false));
        flags.resize(rows, vector<bool>(cols, false));
    }

    void generateBoard(int firstRow, int firstCol) {
        random_device rd;
        mt19937 gen(rd());
        uniform_int_distribution<> distRow(0, rows-1);
        uniform_int_distribution<> distCol(0, cols-1);
        int placed = 0;
        while (placed < minesTotal) {
            int r = distRow(gen);
            int c = distCol(gen);
            if (r == firstRow && c == firstCol) continue;
            if (abs(r - firstRow) <= 1 && abs(c - firstCol) <= 1) continue;
            if (real[r][c] == -1) continue;
            real[r][c] = -1;
            placed++;
        }
        // Подсчёт чисел
        for (int r = 0; r < rows; ++r) {
            for (int c = 0; c < cols; ++c) {
                if (real[r][c] == -1) continue;
                int cnt = 0;
                for (int dr = -1; dr <= 1; ++dr) {
                    for (int dc = -1; dc <= 1; ++dc) {
                        int nr = r + dr, nc = c + dc;
                        if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && real[nr][nc] == -1)
                            cnt++;
                    }
                }
                real[r][c] = cnt;
            }
        }
    }

    void render() {
        cout << "\033[2J\033[1;1H";
        cout << colorize("   ", BOLD);
        for (int c = 0; c < cols; ++c) {
            cout << colorize(to_string(c+1) + (c+1 < 10 ? "  " : " "), BOLD);
        }
        cout << endl;
        cout << colorize("   +", GRAY);
        for (int c = 0; c < cols; ++c) cout << "---";
        cout << colorize("+\n", GRAY);
        for (int r = 0; r < rows; ++r) {
            cout << colorize(to_string(r+1) + (r+1 < 10 ? " " : ""), BOLD);
            cout << colorize("|", GRAY);
            for (int c = 0; c < cols; ++c) {
                if (gameOver && real[r][c] == -1) {
                    cout << colorize(" 💣", RED);
                } else if (flags[r][c]) {
                    cout << colorize(" ⚑", YELLOW);
                } else if (revealed[r][c]) {
                    int val = real[r][c];
                    if (val == -1) cout << colorize(" 💣", RED);
                    else if (val == 0) cout << "  ";
                    else {
                        vector<string> cols2 = {"", BLUE, GREEN, RED, MAGENTA, CYAN, YELLOW, WHITE, GRAY};
                        cout << colorize(" " + to_string(val), cols2[val]);
                    }
                } else {
                    cout << " ■";
                }
            }
            cout << colorize("|\n", GRAY);
        }
        cout << colorize("   +", GRAY);
        for (int c = 0; c < cols; ++c) cout << "---";
        cout << colorize("+\n", GRAY);
        cout << colorize("Мин осталось: " + to_string(minesRemaining), YELLOW) << endl;
        if (!gameOver && startTime != 0) {
            int elapsed = (int)(time(nullptr) - startTime);
            cout << colorize("Время: " + to_string(elapsed) + " сек", BLUE) << endl;
        }
        if (gameOver) {
            cout << (won ? colorize("🎉 ПОБЕДА!", GREEN) : colorize("💥 ПОРАЖЕНИЕ!", RED)) << endl;
        }
    }

    void openCell(int r, int c) {
        if (gameOver || revealed[r][c] || flags[r][c]) return;
        if (firstMove) {
            generateBoard(r, c);
            firstMove = false;
            startTime = time(nullptr);
        }
        revealed[r][c] = true;
        moves++;
        if (real[r][c] == -1) {
            gameOver = true;
            return;
        }
        if (real[r][c] == 0) {
            for (int dr = -1; dr <= 1; ++dr) {
                for (int dc = -1; dc <= 1; ++dc) {
                    int nr = r + dr, nc = c + dc;
                    if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && !revealed[nr][nc] && !flags[nr][nc]) {
                        openCell(nr, nc);
                    }
                }
            }
        }
        // Проверка победы
        int revealedCount = 0;
        for (int i = 0; i < rows; ++i)
            for (int j = 0; j < cols; ++j)
                if (revealed[i][j]) revealedCount++;
        if (revealedCount == rows * cols - minesTotal) {
            won = true;
            gameOver = true;
        }
    }

    void toggleFlag(int r, int c) {
        if (gameOver || revealed[r][c]) return;
        if (flags[r][c]) {
            flags[r][c] = false;
            minesRemaining++;
        } else {
            flags[r][c] = true;
            minesRemaining--;
        }
    }

    bool isGameOver() const { return gameOver; }
    bool isWon() const { return won; }
    int getMoves() const { return moves; }
    time_t getStartTime() const { return startTime; }

private:
    int rows, cols, minesTotal, minesRemaining;
    vector<vector<char>> board;
    vector<vector<int>> real;
    vector<vector<bool>> revealed, flags;
    bool gameOver = false, won = false, firstMove;
    int moves = 0;
    time_t startTime = 0;
};

int main(int argc, char* argv[]) {
    srand(time(nullptr));
    int rows = 9, cols = 9, mines = 10;
    if (argc > 1) {
        string arg = argv[1];
        if (arg == "easy") { rows=9; cols=9; mines=10; }
        else if (arg == "medium") { rows=16; cols=16; mines=40; }
        else if (arg == "hard") { rows=30; cols=16; mines=99; }
        else {
            int v = stoi(arg);
            rows = cols = v;
            if (rows < 3) rows = cols = 3;
            if (argc > 2) mines = stoi(argv[2]);
            else mines = max(1, rows*cols/10);
        }
    }
    Minesweeper game(rows, cols, mines);
    string cmd;
    while (!game.isGameOver()) {
        game.render();
        cout << "Введите команду (строка столбец для открытия, f строка столбец для флага, q для выхода): ";
        getline(cin, cmd);
        if (cmd == "q") { cout << "Выход." << endl; return 0; }
        stringstream ss(cmd);
        string first;
        ss >> first;
        if (first == "f") {
            int r, c;
            if (ss >> r >> c) {
                if (r >= 1 && r <= rows && c >= 1 && c <= cols) {
                    game.toggleFlag(r-1, c-1);
                } else cout << "Координаты вне поля." << endl;
            } else cout << "Неверный формат." << endl;
        } else {
            int r, c;
            if (ss >> r >> c) {
                if (r >= 1 && r <= rows && c >= 1 && c <= cols) {
                    game.openCell(r-1, c-1);
                } else cout << "Координаты вне поля." << endl;
            } else cout << "Неверный формат." << endl;
        }
    }
    game.render();
    if (game.isWon()) {
        int elapsed = (int)(time(nullptr) - game.getStartTime());
        cout << "Поздравляем! Вы выиграли за " << game.getMoves() << " ходов и " << elapsed << " секунд." << endl;
    } else {
        cout << "Игра окончена. Попробуйте снова!" << endl;
    }
    return 0;
}
