require 'test_helper'

describe Aquatone::Validation do
  describe ".valid_domain_name?" do
    describe "when given invalid domain names" do
      it "returns false" do
        ["invalid", "invalid.", 1, "!!!!!", "1337", ".invalid"].each do |invalid_domain|
          Aquatone::Validation.valid_domain_name?(invalid_domain).must_equal false
        end
      end
    end

    describe "when given valid domain names" do
      it "returns true" do
        ["microsoft.com", "subdomain.example.net", "nsa.gov", "something.store", "verylongverylongverylong.verylong-verylongverylongverylong.loan"].each do |valid_domain|
          Aquatone::Validation.valid_domain_name?(valid_domain).must_equal true
        end
      end
    end
  end

  describe ".valid_ip?" do
    describe "when given invalid IPs" do
      it "returns false" do
        ["1.2", "!!!!!", 1234, "....", "543.513.1.54", "12.43.234.15."].each do |invalid_ip|
          Aquatone::Validation.valid_ip?(invalid_ip).must_equal false
        end
      end
    end

    describe "when given valid IPs" do
      it "returns true" do
        ["192.168.1.1", "10.0.0.42", "8.8.8.8", "104.40.211.35", "2001:db8:0:0:0:0:2:1", "2001:db8::2:1"].each do |valid_ip|
          Aquatone::Validation.valid_ip?(valid_ip).must_equal true
        end
      end
    end
  end

  describe ".valid_tcp_port?" do
    describe "when given invalid TCP ports" do
      it "returns false" do
        [0, "0", "9999999", 65536, -12, "-65", "!!!!", "invalid"].each do |invalid_port|
          Aquatone::Validation.valid_tcp_port?(invalid_port).must_equal false
        end
      end
    end

    describe "when given valid TCP ports" do
      it "returns true" do
        [1, "1", 65535, "65535", 80, 1024, "8080", "500", "32213"].each do |valid_port|
          Aquatone::Validation.valid_tcp_port?(valid_port).must_equal true
        end
      end
    end

    describe "built-in port lists" do
      describe Aquatone::PortLists::SMALL do
        it "must be valid" do
          Aquatone::PortLists::SMALL.each do |port|
            Aquatone::Validation.valid_tcp_port?(port).must_equal true
          end
        end
      end

      describe Aquatone::PortLists::MEDIUM do
        it "must be valid" do
          Aquatone::PortLists::MEDIUM.each do |port|
            Aquatone::Validation.valid_tcp_port?(port).must_equal true
          end
        end
      end

      describe Aquatone::PortLists::LARGE do
        it "must be valid" do
          Aquatone::PortLists::LARGE.each do |port|
            Aquatone::Validation.valid_tcp_port?(port).must_equal true
          end
        end
      end

      describe Aquatone::PortLists::HUGE do
        it "must be valid" do
          Aquatone::PortLists::HUGE.each do |port|
            Aquatone::Validation.valid_tcp_port?(port).must_equal true
          end
        end
      end
    end
  end
end
