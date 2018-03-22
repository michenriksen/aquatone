module Aquatone
  module PortLists
    SMALL  = [80, 443].freeze

    MEDIUM = [80, 443, 8000, 8080, 8443].freeze

    LARGE  = [80,   81,   443,  591,  2082, 2087, 2095, 2096, 3000, 8000, 8001,
              8008, 8080, 8083, 8443, 8834, 8888].freeze

    HUGE   = [80,    81,    300,   443,   591,   593,   832,   981,   1010,  1311,
              2082,  2087,  2095,  2096,  2480,  3000,  3128,  3333,  4243,  4567,
              4711,  4712,  4993,  5000,  5104,  5108,  5800,  6543,  6379,  7000,
              7396,  7474,  8000,  8001,  8008,  8014,  8042,  8069,  8080,  8081,
              8088,  8090,  8091,  8118,  8123,  8172,  8222,  8243,  8280,  8281,
              8333,  8443,  8500,  8834,  8880,  8888,  8983,  9000,  9043,  9060,
              9080,  9090,  9091,  9200,  9443,  9800,  9981,  1000, 12443,  16080,
              18091, 18092, 27018, 20720, 28017].freeze
              

    class UnknownPortListName < StandardError; end

    def self.port_list_by_name(name)
      case name.downcase.strip
      when "small"
        return self::SMALL
      when "medium", "default"
        return self::MEDIUM
      when "large"
        return self::LARGE
      when "huge", "xlarge"
        return self::HUGE
      else
        fail UnknownPortListName, "Unknown port list name: #{name}"
      end
    end
  end
end
