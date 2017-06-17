module Aquatone
  module Validation
    DOMAIN_REGEX = /\A([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?\z/.freeze
    MIN_PORT     = 1.freeze
    MAX_PORT     = 65535.freeze

    def self.valid_domain_name?(value)
      value.to_s =~ DOMAIN_REGEX ? true : false
    end

    def self.valid_ip?(value)
      IPAddr.new(value)
      true
    rescue IPAddr::Error
      false
    end

    def self.valid_tcp_port?(value)
      value.to_i.between?(MIN_PORT, MAX_PORT)
    end
  end
end
