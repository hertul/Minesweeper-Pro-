#!/usr/bin/env ruby
# minesweeper.rb
# encoding: UTF-8

require 'io/console'
require 'timeout'

COLORS = {
  reset: "\e[0m",
  red: "\e[91m",
  green: "\e[92m",
  yellow: "\e[93m",
  blue: "\e[94m",
  magenta: "\e[95m",
  cyan: "\e[96m",
  white: "\e[97m",
  gray: "\e[90m",
  bold: "\e[1m"
}

def colorize(text, color)
  "#{COLORS[color]}#{text}#{COLORS[:reset]}"
end

class Minesweeper
  attr_reader :rows, :cols, :mines_total, :mines_remaining, :game_over, :won, :moves, :start_time

  def initialize(rows, cols, mines)
    @rows = rows
    @cols = cols
    @mines_total = mines
    @mines_remaining = mines
    @board = Array.new(rows) { Array.new(cols, ' ') }
    @real = Array.new(rows) { Array.new(cols, 0) }
    @revealed = Array.new(rows) { Array.new(cols, false) }
    @flags = Array.new(rows) { Array.new(cols, false) }
    @game_over = false
    @won = false
    @first_move = true
    @moves = 0
    @start_time = nil
  end

  def generate_board(first_row, first_col)
    placed = 0
    while placed < @mines_total
      r = rand(@rows)
      c = rand(@cols)
      next if r == first_row && c == first_col
      next if (r - first_row).abs <= 1 && (c - first_col).abs <= 1
      next if @real[r][c] == -1
      @real[r][c] = -1
      placed += 1
    end
    @rows.times do |r|
      @cols.times do |c|
        next if @real[r][c] == -1
        cnt = 0
        (-1..1).each do |dr|
          (-1..1).each do |dc|
            nr, nc = r+dr, c+dc
            cnt += 1 if nr >= 0 && nr < @rows && nc >= 0 && nc < @cols && @real[nr][nc] == -1
          end
        end
        @real[r][c] = cnt
      end
    end
  end

  def render
    system('clear') || system('cls')
    print colorize('   ', :bold)
    @cols.times { |c| print colorize("%2d " % (c+1), :bold) }
    puts
    print colorize('   +', :gray)
    @cols.times { print '---' }
    puts colorize('+', :gray)
    @rows.times do |r|
      print colorize("%2d " % (r+1), :bold)
      print colorize('|', :gray)
      @cols.times do |c|
        if @game_over && @real[r][c] == -1
          print colorize(' 💣', :red)
        elsif @flags[r][c]
          print colorize(' ⚑', :yellow)
        elsif @revealed[r][c]
          val = @real[r][c]
          if val == -1
            print colorize(' 💣', :red)
          elsif val == 0
            print '  '
          else
            colors = ['', :blue, :green, :red, :magenta, :cyan, :yellow, :white, :gray]
            print colorize(" #{val}", colors[val])
          end
        else
          print ' ■'
        end
      end
      puts colorize('|', :gray)
    end
    print colorize('   +', :gray)
    @cols.times { print '---' }
    puts colorize('+', :gray)
    puts colorize("Мин осталось: #{@mines_remaining}", :yellow)
    if !@game_over && @start_time
      elapsed = (Time.now - @start_time).to_i
      puts colorize("Время: #{elapsed} сек", :blue)
    end
    puts @game_over ? (@won ? colorize('🎉 ПОБЕДА!', :green) : colorize('💥 ПОРАЖЕНИЕ!', :red)) : ''
  end

  def open_cell(r, c)
    return if @game_over || @revealed[r][c] || @flags[r][c]
    if @first_move
      generate_board(r, c)
      @first_move = false
      @start_time = Time.now
    end
    @revealed[r][c] = true
    @moves += 1
    if @real[r][c] == -1
      @game_over = true
      return
    end
    if @real[r][c] == 0
      (-1..1).each do |dr|
        (-1..1).each do |dc|
          nr, nc = r+dr, c+dc
          if nr >= 0 && nr < @rows && nc >= 0 && nc < @cols && !@revealed[nr][nc] && !@flags[nr][nc]
            open_cell(nr, nc)
          end
        end
      end
    end
    # Проверка победы
    revealed_count = @rows.times.sum { |i| @cols.times.count { |j| @revealed[i][j] } }
    if revealed_count == @rows * @cols - @mines_total
      @won = true
      @game_over = true
    end
  end

  def toggle_flag(r, c)
    return if @game_over || @revealed[r][c]
    if @flags[r][c]
      @flags[r][c] = false
      @mines_remaining += 1
    else
      @flags[r][c] = true
      @mines_remaining -= 1
    end
  end
end

def main
  rows, cols, mines = 9, 9, 10
  if ARGV.length > 0
    level = ARGV[0].downcase
    case level
    when 'easy'
      rows, cols, mines = 9, 9, 10
    when 'medium'
      rows, cols, mines = 16, 16, 40
    when 'hard'
      rows, cols, mines = 30, 16, 99
    else
      v = level.to_i
      if v > 0
        rows = cols = v
        rows = cols = 3 if rows < 3
        mines = ARGV[1].to_i if ARGV.length > 1
        mines = [1, rows*cols/10].max if mines == 0
      end
    end
  end
  game = Minesweeper.new(rows, cols, mines)

  while !game.game_over
    game.render
    print "Введите команду (строка столбец для открытия, f строка столбец для флага, q для выхода): "
    cmd = STDIN.gets.chomp.strip
    if cmd == 'q'
      puts "Выход."
      return
    end
    parts = cmd.split
    if parts.size < 2
      puts "Неверный формат."
      next
    end
    if parts[0] == 'f'
      if parts.size != 3
        puts "Неверный формат. Используйте: f строка столбец"
        next
      end
      r = parts[1].to_i - 1
      c = parts[2].to_i - 1
      if r >= 0 && r < rows && c >= 0 && c < cols
        game.toggle_flag(r, c)
      else
        puts "Координаты вне поля."
      end
    else
      if parts.size != 2
        puts "Неверный формат. Используйте: строка столбец"
        next
      end
      r = parts[0].to_i - 1
      c = parts[1].to_i - 1
      if r >= 0 && r < rows && c >= 0 && c < cols
        game.open_cell(r, c)
      else
        puts "Координаты вне поля."
      end
    end
  end
  game.render
  if game.won
    elapsed = (Time.now - game.start_time).to_i
    puts "Поздравляем! Вы выиграли за #{game.moves} ходов и #{elapsed} секунд."
  else
    puts "Игра окончена. Попробуйте снова!"
  end
end

main if __FILE__ == $0
