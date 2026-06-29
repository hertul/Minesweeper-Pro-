// minesweeper.java
import java.io.*;
import java.util.*;

public class minesweeper {
    private static final String RESET = "\u001B[0m";
    private static final String RED = "\u001B[91m";
    private static final String GREEN = "\u001B[92m";
    private static final String YELLOW = "\u001B[93m";
    private static final String BLUE = "\u001B[94m";
    private static final String MAGENTA = "\u001B[95m";
    private static final String CYAN = "\u001B[96m";
    private static final String WHITE = "\u001B[97m";
    private static final String GRAY = "\u001B[90m";
    private static final String BOLD = "\u001B[1m";

    private static String colorize(String text, String color) {
        return color + text + RESET;
    }

    static class MinesweeperGame {
        int rows, cols, minesTotal, minesRemaining;
        int[][] real;
        boolean[][] revealed, flags;
        boolean gameOver, won, firstMove;
        int moves;
        long startTime;

        MinesweeperGame(int rows, int cols, int mines) {
            this.rows = rows;
            this.cols = cols;
            minesTotal = mines;
            minesRemaining = mines;
            real = new int[rows][cols];
            revealed = new boolean[rows][cols];
            flags = new boolean[rows][cols];
            firstMove = true;
        }

        void generateBoard(int firstRow, int firstCol) {
            Random rand = new Random();
            int placed = 0;
            while (placed < minesTotal) {
                int r = rand.nextInt(rows);
                int c = rand.nextInt(cols);
                if (r == firstRow && c == firstCol) continue;
                if (Math.abs(r - firstRow) <= 1 && Math.abs(c - firstCol) <= 1) continue;
                if (real[r][c] == -1) continue;
                real[r][c] = -1;
                placed++;
            }
            for (int r = 0; r < rows; r++) {
                for (int c = 0; c < cols; c++) {
                    if (real[r][c] == -1) continue;
                    int cnt = 0;
                    for (int dr = -1; dr <= 1; dr++) {
                        for (int dc = -1; dc <= 1; dc++) {
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
            System.out.print("\033[H\033[2J");
            System.out.flush();
            System.out.print(colorize("   ", BOLD));
            for (int c = 0; c < cols; c++) {
                System.out.print(colorize(String.format("%2d ", c+1), BOLD));
            }
            System.out.println();
            System.out.print(colorize("   +", GRAY));
            for (int c = 0; c < cols; c++) System.out.print("---");
            System.out.println(colorize("+", GRAY));
            for (int r = 0; r < rows; r++) {
                System.out.print(colorize(String.format("%2d ", r+1), BOLD));
                System.out.print(colorize("|", GRAY));
                for (int c = 0; c < cols; c++) {
                    if (gameOver && real[r][c] == -1) {
                        System.out.print(colorize(" 💣", RED));
                    } else if (flags[r][c]) {
                        System.out.print(colorize(" ⚑", YELLOW));
                    } else if (revealed[r][c]) {
                        int val = real[r][c];
                        if (val == -1) System.out.print(colorize(" 💣", RED));
                        else if (val == 0) System.out.print("  ");
                        else {
                            String[] colors = {"", BLUE, GREEN, RED, MAGENTA, CYAN, YELLOW, WHITE, GRAY};
                            System.out.print(colorize(" " + val, colors[val]));
                        }
                    } else {
                        System.out.print(" ■");
                    }
                }
                System.out.println(colorize("|", GRAY));
            }
            System.out.print(colorize("   +", GRAY));
            for (int c = 0; c < cols; c++) System.out.print("---");
            System.out.println(colorize("+", GRAY));
            System.out.println(colorize("Мин осталось: " + minesRemaining, YELLOW));
            if (!gameOver && startTime != 0) {
                long elapsed = (System.currentTimeMillis() - startTime) / 1000;
                System.out.println(colorize("Время: " + elapsed + " сек", BLUE));
            }
            if (gameOver) {
                System.out.println(won ? colorize("🎉 ПОБЕДА!", GREEN) : colorize("💥 ПОРАЖЕНИЕ!", RED));
            }
        }

        void openCell(int r, int c) {
            if (gameOver || revealed[r][c] || flags[r][c]) return;
            if (firstMove) {
                generateBoard(r, c);
                firstMove = false;
                startTime = System.currentTimeMillis();
            }
            revealed[r][c] = true;
            moves++;
            if (real[r][c] == -1) {
                gameOver = true;
                return;
            }
            if (real[r][c] == 0) {
                for (int dr = -1; dr <= 1; dr++) {
                    for (int dc = -1; dc <= 1; dc++) {
                        int nr = r + dr, nc = c + dc;
                        if (nr >= 0 && nr < rows && nc >= 0 && nc < cols && !revealed[nr][nc] && !flags[nr][nc]) {
                            openCell(nr, nc);
                        }
                    }
                }
            }
            int revealedCount = 0;
            for (int i = 0; i < rows; i++)
                for (int j = 0; j < cols; j++)
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
    }

    public static void main(String[] args) throws IOException {
        int rows = 9, cols = 9, mines = 10;
        if (args.length > 0) {
            String level = args[0].toLowerCase();
            if (level.equals("easy")) { rows = 9; cols = 9; mines = 10; }
            else if (level.equals("medium")) { rows = 16; cols = 16; mines = 40; }
            else if (level.equals("hard")) { rows = 30; cols = 16; mines = 99; }
            else {
                try {
                    int v = Integer.parseInt(level);
                    rows = cols = v;
                    if (rows < 3) rows = cols = 3;
                    if (args.length > 1) mines = Integer.parseInt(args[1]);
                    else mines = Math.max(1, rows * cols / 10);
                } catch (NumberFormatException e) {
                    System.out.println("Неверный аргумент.");
                    return;
                }
            }
        }
        MinesweeperGame game = new MinesweeperGame(rows, cols, mines);
        BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));
        while (!game.gameOver) {
            game.render();
            System.out.print("Введите команду (строка столбец для открытия, f строка столбец для флага, q для выхода): ");
            String cmd = reader.readLine().trim();
            if (cmd.equals("q")) {
                System.out.println("Выход.");
                return;
            }
            String[] parts = cmd.split("\\s+");
            if (parts.length < 2) {
                System.out.println("Неверный формат.");
                continue;
            }
            if (parts[0].equals("f")) {
                if (parts.length != 3) {
                    System.out.println("Неверный формат. Используйте: f строка столбец");
                    continue;
                }
                try {
                    int r = Integer.parseInt(parts[1]) - 1;
                    int c = Integer.parseInt(parts[2]) - 1;
                    if (r >= 0 && r < rows && c >= 0 && c < cols)
                        game.toggleFlag(r, c);
                    else
                        System.out.println("Координаты вне поля.");
                } catch (NumberFormatException e) {
                    System.out.println("Неверный формат.");
                }
            } else {
                if (parts.length != 2) {
                    System.out.println("Неверный формат. Используйте: строка столбец");
                    continue;
                }
                try {
                    int r = Integer.parseInt(parts[0]) - 1;
                    int c = Integer.parseInt(parts[1]) - 1;
                    if (r >= 0 && r < rows && c >= 0 && c < cols)
                        game.openCell(r, c);
                    else
                        System.out.println("Координаты вне поля.");
                } catch (NumberFormatException e) {
                    System.out.println("Неверный формат.");
                }
            }
        }
        game.render();
        if (game.won) {
            long elapsed = (System.currentTimeMillis() - game.startTime) / 1000;
            System.out.printf("Поздравляем! Вы выиграли за %d ходов и %d секунд.\n", game.moves, elapsed);
        } else {
            System.out.println("Игра окончена. Попробуйте снова!");
        }
    }
}
