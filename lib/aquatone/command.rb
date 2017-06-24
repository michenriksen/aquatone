module Aquatone
  class Command
    attr_reader :options

    def initialize(options)
      @options = options
    end

    def self.run(options)
      self.new(options).execute!
    rescue Interrupt
      output("Caught interrupt; exiting.\n", true)
    rescue => e
      output(red("An unexpected error occurred: #{e.class}: #{e.message}\n"))
      output("#{e.backtrace.join("\n")}\n")
    end

    def execute!
      fail NotImplementedError, "Commands must overwrite #execute! method"
    end

    protected

    def output(text, clear_line = false)
      self.class.output(text, clear_line)
    end

    def status(message)
      self.class.status(message)
    end

    def self.output(text, clear_line = false)
      if clear_line
        text = "\r\e[0K#{text}"
      end
      with_output_mutex { print text }
    end

    def self.status(message)
      output(message, true)
    end

    def self.colorize(text, color_code)
      "\e[1m\e[#{color_code}m#{text}\e[0m"
    end

    def uncolorize(text)
      self.class.uncolorize(text)
    end

    def red(text)
      self.class.red(text)
    end

    def green(text)
      self.class.green(text)
    end

    def blue(text)
      self.class.blue(text)
    end

    def yellow(text)
      self.class.yellow(text)
    end

    def bold(text)
      self.class.bold(text)
    end

    def banner(subtitle)
      output("                           __\n")
      output("  ____ _____ ___  ______ _/ /_____  ____  ___\n")
      output(" / __ `/ __ `/ / / / __ `/ __/ __ \\/ __ \\/ _ \\\n")
      output("/ /_/ / /_/ / /_/ / /_/ / /_/ /_/ / / / /  __/\n")
      output("\\__,_/\\__, /\\__,_/\\__,_/\\__/\\____/_/ /_/\\___/\n")
      output("        /_/  #{subtitle.downcase} v#{Aquatone::VERSION} - by @michenriksen\n\n")
    end

    def self.uncolorize(text); colorize(text.to_s, 0); end
    def self.red(text); colorize(text.to_s, 31); end
    def self.green(text); colorize(text.to_s, 32); end
    def self.blue(text); colorize(text.to_s, 34); end
    def self.yellow(text); colorize(text.to_s, 33); end
    def self.bold(text); colorize(text.to_s, 1); end

    def thread_pool
      if options[:sleep]
        Aquatone::ThreadPool.new(1)
      else
        Aquatone::ThreadPool.new(options[:threads])
      end
    end

    def jitter_sleep
      return unless options[:sleep]
      seconds = options[:sleep].to_i
      if options[:jitter]
        jitter  = (options[:jitter].to_f / 100) * seconds
        if rand.round == 0
          seconds = seconds - Random.rand(0..jitter.round)
        else
          seconds = seconds + Random.rand(0..jitter.round)
        end
        seconds = 1 if seconds < 1
      end
      sleep seconds
    end

    def seconds_to_time(seconds)
      Time.at(seconds).utc.strftime("%H:%M:%S")
    end

    def asked_for_progress?
      begin
        while c = STDIN.read_nonblock(1)
          return true if c == "\n"
        end
        false
      rescue IO::EAGAINWaitReadable, Errno::EBADF
        false
      rescue Errno::EINTR, Errno::EAGAIN, EOFError
        true
      end
    rescue NameError
      false
    end

    def has_executable?(name)
      exts = ENV['PATHEXT'] ? ENV['PATHEXT'].split(';') : ['']
      ENV["PATH"].split(File::PATH_SEPARATOR).each do |path|
        exts.each do |ext|
          exe = File.join(path, "#{name}#{ext}")
          return true if File.executable?(exe) && !File.directory?(exe)
        end
      end
      false
    end

    def self.jitter
      seconds = @options[:jitter]
      if seconds != 0
        random_sleep = ((1 - (rand(30) * 0.01)) * seconds)
        sleep(random_sleep)
      end
    end

    def self.with_output_mutex
      output_mutex.synchronize { yield }
    end

    def self.output_mutex
      @output_mutex ||= Mutex.new
    end
  end
end
