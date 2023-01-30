$cols = `tput col`.chomp.to_i
$lines = `tput lines`.chomp.to_i

class Bar
  def initialize(w, h)
    @vars = []
    @width, @height = w, h

    @pos_self_left = proc do |x|
      print("\033[#{$lines - @height+x};0H")
    end

    @pos_self_right = proc do |x|
      print("\033[#{$lines - @height+x};#{@width}H")
    end
  end

  def display
    print("\0337")

    @pos_self_left.call(0)
    @width.times { print("\033[0;4;1;43;35m_\033[0m") }

    @height.times do |n|
      @pos_self_left.call(n+1)
      print("\033[0;4;1;43;35m|")

      @width.times do
        print(" ")
      end

      @pos_self_right.call(n+1)
      print("|\033[0m")
    end

    print("\0338")
  end
end

system("clear")
b = Bar.new($cols, $lines / 8)
b.display
