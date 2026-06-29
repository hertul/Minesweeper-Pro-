// minesweeper.cs
using System;
using System.Collections.Generic;
using System.Linq;

class Minesweeper
{
    static string Colorize(string text, string color)
    {
        string col = color switch
        {
            "red" => "\x1b[91m",
            "green" => "\x1b[92m",
            "yellow" => "\x1b[93m",
            "blue" => "\x1b[94m",
            "magenta" => "\x1b[95m",
            "cyan" => "\x1b[96m",
            "white" => "\x1b[97m",
            "gray" => "\x1b[90m",
            "bold" => "\x1b[1m",
            _ => "\x1b[0m"
        };
        return col + text + "\x1b[0m";
    }

    private int rows, cols, minesTotal, minesRemaining;
    private char[,] board;
    private int[,] real;
    private bool[,] revealed, flags;
    private bool gameOver, won, firstMove;
    private int moves;
    private DateTime startTime;

    public Minesweeper(int rows, int cols, int mines)
    {
        this.rows = rows;
        this.cols = cols;
        minesTotal = mines;
        minesRemaining = mines;
        firstMove = true;
        board = new char[rows, cols];
        real = new int[rows, cols];
        revealed = new bool[rows, cols];
        flags = new bool[rows, cols];
    }

    private void GenerateBoard(int firstRow, int firstCol)
    {
        Random rnd = new Random();
        int placed = 0;
        while (placed < minesTotal)
        {
            int r = rnd.Next(rows);
            int c = rnd.Next(cols);
            if (r == firstRow && c == firstCol) continue;
            if (Math.Abs(r - firstRow) <= 1 && Math.Abs(c - firstCol) <= 1) continue;
            if (real[r, c] == -1) continue;
            real[r, c] = -1;
            placed++;
        }
        for (int r = 0; r < rows; r++)
            for (int c = 0; c < cols; c++)
            {
                if (real[r, c] == -1) continue;
                int cnt = 0;
                for (int dr = -1; dr <= 1; dr++)
                    for (int dc = -1; dc <= 1; dc++)
                    {
                        int nr = r + dr, nc = c + dc;
                        if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && real[nr, nc] == -1)
                            cnt++;
                    }
                real[r, c] = cnt;
            }
    }

    public void Render()
    {
        Console.Clear();
        Console.Write(Colorize("   ", "bold"));
        for (int c = 0; c < cols; c++)
            Console.Write(Colorize($"{c+1,2} ", "bold"));
        Console.WriteLine();
        Console.Write(Colorize("   +", "gray"));
        for (int c = 0; c < cols; c++) Console.Write("---");
        Console.WriteLine(Colorize("+", "gray"));
        for (int r = 0; r < rows; r++)
        {
            Console.Write(Colorize($"{r+1,2} ", "bold"));
            Console.Write(Colorize("|", "gray"));
            for (int c = 0; c < cols; c++)
            {
                if (gameOver && real[r, c] == -1)
                    Console.Write(Colorize(" 💣", "red"));
                else if (flags[r, c])
                    Console.Write(Colorize(" ⚑", "yellow"));
                else if (revealed[r, c])
                {
                    int val = real[r, c];
                    if (val == -1) Console.Write(Colorize(" 💣", "red"));
                    else if (val == 0) Console.Write("  ");
                    else
                    {
                        string[] colors = { "", "blue", "green", "red", "magenta", "cyan", "yellow", "white", "gray" };
                        Console.Write(Colorize($" {val}", colors[val]));
                    }
                }
                else Console.Write(" ■");
            }
            Console.WriteLine(Colorize("|", "gray"));
        }
        Console.Write(Colorize("   +", "gray"));
        for (int c = 0; c < cols; c++) Console.Write("---");
        Console.WriteLine(Colorize("+", "gray"));
        Console.WriteLine(Colorize($"Мин осталось: {minesRemaining}", "yellow"));
        if (!gameOver && startTime != DateTime.MinValue)
        {
            int elapsed = (int)(DateTime.Now - startTime).TotalSeconds;
            Console.WriteLine(Colorize($"Время: {elapsed} сек", "blue"));
        }
        if (gameOver)
            Console.WriteLine(won ? Colorize("🎉 ПОБЕДА!", "green") : Colorize("💥 ПОРАЖЕНИЕ!", "red"));
    }

    public void OpenCell(int r, int c)
    {
        if (gameOver || revealed[r, c] || flags[r, c]) return;
        if (firstMove)
        {
            GenerateBoard(r, c);
            firstMove = false;
            startTime = DateTime.Now;
        }
        revealed[r, c] = true;
        moves++;
        if (real[r, c] == -1)
        {
            gameOver = true;
            return;
        }
        if (real[r, c] == 0)
        {
            for (int dr = -1; dr <= 1; dr++)
                for (int dc = -1; dc <= 1; dc++)
                {
                    int nr = r + dr, nc = c + dc;
                    if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && !revealed[nr, nc] && !flags[nr, nc])
                        OpenCell(nr, nc);
                }
        }
        // Проверка победы
        int revealedCount = 0;
        for (int i = 0; i < rows; i++)
            for (int j = 0; j < cols; j++)
                if (revealed[i, j]) revealedCount++;
        if (revealedCount == rows * cols - minesTotal)
        {
            won = true;
            gameOver = true;
        }
    }

    public void ToggleFlag(int r, int c)
    {
        if (gameOver || revealed[r, c]) return;
        if (flags[r, c])
        {
            flags[r, c] = false;
            minesRemaining++;
        }
        else
        {
            flags[r, c] = true;
            minesRemaining--;
        }
    }

    public bool IsGameOver => gameOver;
    public bool IsWon => won;
    public int Moves => moves;
    public DateTime StartTime => startTime;

    static void Main(string[] args)
    {
        int rows = 9, cols = 9, mines = 10;
        if (args.Length > 0)
        {
            string level = args[0].ToLower();
            if (level == "easy") { rows = 9; cols = 9; mines = 10; }
            else if (level == "medium") { rows = 16; cols = 16; mines = 40; }
            else if (level == "hard") { rows = 30; cols = 16; mines = 99; }
            else if (int.TryParse(level, out int v))
            {
                rows = cols = v;
                if (rows < 3) rows = cols = 3;
                if (args.Length > 1 && int.TryParse(args[1], out int m)) mines = m;
                else mines = Math.Max(1, rows * cols / 10);
            }
        }
        Minesweeper game = new Minesweeper(rows, cols, mines);
        while (!game.IsGameOver)
        {
            game.Render();
            Console.Write("Введите команду (строка столбец для открытия, f строка столбец для флага, q для выхода): ");
            string input = Console.ReadLine().Trim();
            if (input.ToLower() == "q") { Console.WriteLine("Выход."); return; }
            string[] parts = input.Split(' ');
            if (parts.Length < 2) { Console.WriteLine("Неверный формат."); continue; }
            if (parts[0].ToLower() == "f")
            {
                if (parts.Length != 3) { Console.WriteLine("Неверный формат. Используйте: f строка столбец"); continue; }
                if (int.TryParse(parts[1], out int r) && int.TryParse(parts[2], out int c))
                {
                    if (r >= 1 && r <= rows && c >= 1 && c <= cols)
                        game.ToggleFlag(r-1, c-1);
                    else Console.WriteLine("Координаты вне поля.");
                }
                else Console.WriteLine("Неверный формат.");
            }
            else
            {
                if (parts.Length != 2) { Console.WriteLine("Неверный формат. Используйте: строка столбец"); continue; }
                if (int.TryParse(parts[0], out int r) && int.TryParse(parts[1], out int c))
                {
                    if (r >= 1 && r <= rows && c >= 1 && c <= cols)
                        game.OpenCell(r-1, c-1);
                    else Console.WriteLine("Координаты вне поля.");
                }
                else Console.WriteLine("Неверный формат.");
            }
        }
        game.Render();
        if (game.IsWon)
        {
            int elapsed = (int)(DateTime.Now - game.StartTime).TotalSeconds;
            Console.WriteLine($"Поздравляем! Вы выиграли за {game.Moves} ходов и {elapsed} секунд.");
        }
        else Console.WriteLine("Игра окончена. Попробуйте снова!");
    }
}
