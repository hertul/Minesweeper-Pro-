# minesweeper.py
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import random
import os
import time
import json
from pathlib import Path

# ANSI colors
COLORS = {
    'reset': '\033[0m',
    'red': '\033[91m',
    'green': '\033[92m',
    'yellow': '\033[93m',
    'blue': '\033[94m',
    'magenta': '\033[95m',
    'cyan': '\033[96m',
    'white': '\033[97m',
    'gray': '\033[90m',
    'bold': '\033[1m'
}

def colorize(text, color):
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

# Конфигурация по умолчанию
LEVELS = {
    'easy': (9, 9, 10),
    'medium': (16, 16, 40),
    'hard': (30, 16, 99)
}

RECORD_FILE = Path.home() / '.minesweeper_records.json'

def load_records():
    if RECORD_FILE.exists():
        try:
            with open(RECORD_FILE, 'r') as f:
                return json.load(f)
        except:
            pass
    return {}

def save_records(records):
    with open(RECORD_FILE, 'w') as f:
        json.dump(records, f, indent=2)

class Minesweeper:
    def __init__(self, rows, cols, mines):
        self.rows = rows
        self.cols = cols
        self.mines_total = mines
        self.board = [[' ' for _ in range(cols)] for _ in range(rows)]  # видимое поле
        self.real = [[0 for _ in range(cols)] for _ in range(rows)]    # реальные значения: -1 = мина, 0-8 = число
        self.revealed = [[False for _ in range(cols)] for _ in range(rows)]
        self.flags = [[False for _ in range(cols)] for _ in range(rows)]
        self.game_over = False
        self.won = False
        self.first_move = True
        self.mines_remaining = mines
        self.start_time = None
        self.moves = 0

    def generate_board(self, first_row, first_col):
        # Размещаем мины, избегая первой клетки и её соседей
        mines = []
        while len(mines) < self.mines_total:
            r = random.randint(0, self.rows-1)
            c = random.randint(0, self.cols-1)
            if (r, c) == (first_row, first_col):
                continue
            if abs(r - first_row) <= 1 and abs(c - first_col) <= 1:
                continue
            if (r, c) not in mines:
                mines.append((r, c))
                self.real[r][c] = -1

        # Подсчёт чисел
        for r in range(self.rows):
            for c in range(self.cols):
                if self.real[r][c] == -1:
                    continue
                count = 0
                for dr in (-1,0,1):
                    for dc in (-1,0,1):
                        nr, nc = r+dr, c+dc
                        if 0 <= nr < self.rows and 0 <= nc < self.cols and self.real[nr][nc] == -1:
                            count += 1
                self.real[r][c] = count

    def render(self):
        os.system('clear' if os.name == 'posix' else 'cls')
        # Верхняя строка с номерами столбцов
        print(colorize('   ' + ' '.join(f'{i+1:2}' for i in range(self.cols)), 'bold'))
        # Разделитель
        print(colorize('   ' + '+' + '---' * self.cols + '+', 'gray'))
        for r in range(self.rows):
            line = colorize(f'{r+1:2} ', 'bold') + colorize('|', 'gray')
            for c in range(self.cols):
                if self.game_over and self.real[r][c] == -1:
                    line += colorize(' 💣', 'red')
                elif self.flags[r][c]:
                    line += colorize(' ⚑', 'yellow')
                elif self.revealed[r][c]:
                    val = self.real[r][c]
                    if val == -1:
                        line += colorize(' 💣', 'red')
                    elif val == 0:
                        line += '  '
                    else:
                        colors = ['', 'blue', 'green', 'red', 'magenta', 'cyan', 'yellow', 'white', 'gray']
                        line += colorize(f' {val}', colors[val])
                else:
                    line += ' ■'
            line += colorize('|', 'gray')
            print(line)
        print(colorize('   ' + '+' + '---' * self.cols + '+', 'gray'))
        print(colorize(f"Мин осталось: {self.mines_remaining}", 'yellow'))
        if self.start_time and not self.game_over:
            elapsed = int(time.time() - self.start_time)
            print(colorize(f"Время: {elapsed} сек", 'blue'))
        if self.game_over:
            if self.won:
                print(colorize("🎉 ПОБЕДА!", 'green'))
            else:
                print(colorize("💥 ПОРАЖЕНИЕ!", 'red'))

    def open_cell(self, row, col):
        if self.game_over:
            return
        if self.revealed[row][col] or self.flags[row][col]:
            return
        if self.first_move:
            self.generate_board(row, col)
            self.first_move = False
            self.start_time = time.time()

        self.revealed[row][col] = True
        self.moves += 1

        if self.real[row][col] == -1:
            self.game_over = True
            return

        # Если пустая клетка (0), открываем соседей рекурсивно
        if self.real[row][col] == 0:
            for dr in (-1,0,1):
                for dc in (-1,0,1):
                    nr, nc = row+dr, col+dc
                    if 0 <= nr < self.rows and 0 <= nc < self.cols and not self.revealed[nr][nc] and not self.flags[nr][nc]:
                        self.open_cell(nr, nc)

        # Проверка победы
        total_cells = self.rows * self.cols
        revealed_count = sum(sum(row) for row in self.revealed)
        if revealed_count == total_cells - self.mines_total:
            self.won = True
            self.game_over = True

    def toggle_flag(self, row, col):
        if self.game_over:
            return
        if self.revealed[row][col]:
            return
        if self.flags[row][col]:
            self.flags[row][col] = False
            self.mines_remaining += 1
        else:
            self.flags[row][col] = True
            self.mines_remaining -= 1

    def play(self):
        while not self.game_over:
            self.render()
            try:
                cmd = input("Введите команду (строка столбец для открытия, f строка столбец для флага, q для выхода): ").strip()
                if cmd.lower() == 'q':
                    print("Выход.")
                    sys.exit(0)
                parts = cmd.split()
                if parts[0].lower() == 'f':
                    if len(parts) != 3:
                        print("Неверный формат. Используйте: f строка столбец")
                        continue
                    r = int(parts[1]) - 1
                    c = int(parts[2]) - 1
                    if 0 <= r < self.rows and 0 <= c < self.cols:
                        self.toggle_flag(r, c)
                    else:
                        print("Координаты вне поля.")
                else:
                    if len(parts) != 2:
                        print("Неверный формат. Используйте: строка столбец")
                        continue
                    r = int(parts[0]) - 1
                    c = int(parts[1]) - 1
                    if 0 <= r < self.rows and 0 <= c < self.cols:
                        self.open_cell(r, c)
                    else:
                        print("Координаты вне поля.")
            except ValueError:
                print("Неверный ввод.")
            except KeyboardInterrupt:
                print("\nВыход.")
                sys.exit(0)

        # Конец игры
        self.render()
        if self.won:
            elapsed = int(time.time() - self.start_time)
            print(colorize(f"Поздравляем! Вы выиграли за {self.moves} ходов и {elapsed} секунд.", 'green'))
            # Обновление рекорда
            records = load_records()
            key = f"{self.rows}x{self.cols}"
            if key not in records or elapsed < records[key].get('time', float('inf')):
                records[key] = {'time': elapsed, 'moves': self.moves}
                save_records(records)
                print(colorize("🏆 Новый рекорд!", 'yellow'))
        else:
            print(colorize("Игра окончена. Попробуйте снова!", 'red'))

def main():
    args = sys.argv[1:]
    rows, cols, mines = 9, 9, 10
    if args:
        level = args[0].lower()
        if level in LEVELS:
            rows, cols, mines = LEVELS[level]
            if len(args) > 1:
                mines = int(args[1])
        else:
            try:
                rows = int(level)
                if rows < 3:
                    rows = 3
                cols = rows
                if len(args) > 1:
                    mines = int(args[1])
                else:
                    mines = max(1, rows * cols // 10)
            except ValueError:
                print("Неверный аргумент. Используйте easy, medium, hard или число.")
                sys.exit(1)

    game = Minesweeper(rows, cols, mines)
    game.play()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\nВыход.")
        sys.exit(0)
