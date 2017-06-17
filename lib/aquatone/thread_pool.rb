module Aquatone
  class ThreadPool
    attr_reader :size, :jobs, :pool

    def initialize(size)
      @size = size.to_i
      @jobs = Queue.new
      @pool = Array.new(size) do
        Thread.new do
          catch(:exit) do
            loop do
              job, args = @jobs.pop
              job.call(*args)
            end
          end
        end
      end
    end

    def schedule(*args, &block)
      @jobs << [block, args]
    end

    def shutdown
      @size.times do
        schedule { throw :exit }
      end
      @pool.map(&:join)
    end
  end
end
