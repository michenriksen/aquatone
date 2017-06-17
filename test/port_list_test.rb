require 'test_helper'

describe Aquatone::PortLists do
  describe Aquatone::PortLists::SMALL do
    it "includes expected ports" do
      [80, 443].each do |expected_port|
        Aquatone::PortLists::SMALL.include?(expected_port).must_equal true
      end
    end
  end

  describe Aquatone::PortLists::MEDIUM do
    it "includes expected ports" do
      [80, 443, 8000, 8080, 8443].each do |expected_port|
        Aquatone::PortLists::MEDIUM.include?(expected_port).must_equal true
      end
    end

    it "includes all ports from SMALL list" do
      Aquatone::PortLists::SMALL.each do |port|
        Aquatone::PortLists::MEDIUM.include?(port).must_equal true
      end
    end
  end

  describe Aquatone::PortLists::LARGE do
    it "includes expected ports" do
      [80, 81, 443, 591, 2082, 2095, 2096, 3000, 8000, 8001, 8008, 8080, 8083, 8443, 8834, 8888, 55672].each do |expected_port|
        Aquatone::PortLists::LARGE.include?(expected_port).must_equal true
      end
    end

    it "includes all ports from MEDIUM list" do
      Aquatone::PortLists::MEDIUM.each do |port|
        Aquatone::PortLists::LARGE.include?(port).must_equal true
      end
    end
  end

  describe Aquatone::PortLists::HUGE do
    it "includes expected ports" do
      [80, 81, 300, 443, 591, 593, 832, 981, 1010, 1311, 2082, 2095, 2096, 2480, 3000, 3128, 3333, 4243, 4567, 4711, 4712, 4993, 5000, 5104, 5108, 5280, 5281, 5800, 6543, 7000, 7396, 7474, 8000, 8001, 8008, 8014, 8042, 8069, 8080, 8081, 8083, 8088, 8090, 8091, 8118, 8123, 8172, 8222, 8243, 8280, 8281, 8333, 8337, 8443, 8500, 8834, 8880, 8888, 8983, 9000, 9043, 9060, 9080, 9090, 9091, 9200, 9443, 9800, 9981, 11371, 12443, 16080, 18091, 18092, 20720, 55672].each do |expected_port|
        Aquatone::PortLists::HUGE.include?(expected_port).must_equal true
      end
    end

    it "includes all ports from LARGE list" do
      Aquatone::PortLists::LARGE.each do |port|
        Aquatone::PortLists::HUGE.include?(port).must_equal true
      end
    end
  end

  describe ".port_list_by_name" do
    describe "when given small" do
      it "returns SMALL list" do
        Aquatone::PortLists.port_list_by_name("small").must_equal Aquatone::PortLists::SMALL
      end
    end

    describe "when given medium" do
      it "returns MEDIUM list" do
        Aquatone::PortLists.port_list_by_name("medium").must_equal Aquatone::PortLists::MEDIUM
      end
    end

    describe "when given default" do
      it "returns MEDIUM list" do
        Aquatone::PortLists.port_list_by_name("default").must_equal Aquatone::PortLists::MEDIUM
      end
    end

    describe "when given large" do
      it "returns LARGE list" do
        Aquatone::PortLists.port_list_by_name("large").must_equal Aquatone::PortLists::LARGE
      end
    end

    describe "when given huge" do
      it "returns HUGE list" do
        Aquatone::PortLists.port_list_by_name("huge").must_equal Aquatone::PortLists::HUGE
      end
    end

    describe "when given unknown list name" do
      it "raises an exception" do
        proc { Aquatone::PortLists.port_list_by_name("unknown") }.must_raise Aquatone::PortLists::UnknownPortListName
      end
    end
  end
end
