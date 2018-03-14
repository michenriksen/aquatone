require "ipaddr"

module Aquatone
  module Detectors
    class Squarespace < Aquatone::Detector
      self.meta = {
        :service         => "Squarespace",
        :service_website => "https://www.squarespace.com/",
        :author          => "Duarte Duarte (@dduarte)",
        :description     => "Website builder"
      }

      RESPONSE_FINGERPRINT = "Squarespace - Claim This Domain".freeze
      APEX_VALUES          = [IPAddr.new("198.185.159.0/24"), IPAddr.new("198.49.23.0/24")].freeze

      def run
        if apex_resource?
          return false unless APEX_VALUES.include?(IPAddr.new(resource_value))
        end
        return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
