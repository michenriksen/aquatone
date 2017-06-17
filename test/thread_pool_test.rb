require 'test_helper'

describe Aquatone::ThreadPool do
  describe ".initialize" do
    it "initializes the thread pool with given size" do
      pool = Aquatone::ThreadPool.new(5)
      pool.pool.size.must_equal 5
      pool.pool.each do |t|
        t.must_be_kind_of Thread
      end
    end

    it "initializes an empty job queue" do
      pool = Aquatone::ThreadPool.new(5)
      pool.jobs.must_be_kind_of Queue
      pool.jobs.size.must_equal 0
    end
  end

  describe "#schedule" do
    before do
      @pool_size = 5
      @pool      = Aquatone::ThreadPool.new(@pool_size)
      @mutex     = Mutex.new
    end

    it "executes scheduled blocks" do
      iterations = @pool_size * 3
      results    = Array.new(iterations)

      iterations.times do |i|
        @pool.schedule do
          @mutex.synchronize do
            results[i] = i + 1
          end
        end
      end
      @pool.shutdown

      expected_results = (1.upto(@pool_size * 3)).to_a
      results.must_equal expected_results
    end

    it "executes scheduled blocks in parallel" do
      elapsed = time_taken do
        @pool_size.times do
          @pool.schedule { sleep 1 }
        end
        @pool.shutdown
      end

      elapsed.must_be :<, 4.5
    end
  end
end
