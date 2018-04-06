module Aquatone
  module Detectors
    class Cargo < Aquatone::Detector
      self.meta = {
        :service         => "Cargo",
        :service_website => "https://cargocollective.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Web publishing platform"
      }

      APEX_VALUE            = "173.203.204.123".freeze
      CNAME_VALUE           = "cargocollective.com".freeze
      RESPONSE_FINGERPRINTS = [
        "Use a personal domain name",
        "404 Not Found"
      ].freeze

      def run
        if apex_resource?
          return false unless resource_value == APEX_VALUE
        elsif cname_resource?
          return false unless resource_value == CNAME_VALUE
        end
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINTS)
      end
    end
  end
end
