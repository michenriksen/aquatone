module Aquatone
  module Detectors
    class Unbounce < Aquatone::Detector
      self.meta = {
        :service         => "Unbounce",
        :service_website => "https://unbounce.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Landing page builder and conversion marketing platform"
      }

      APEX_VALUE           = "54.84.104.245".freeze
      CNAME_VALUE          = "unbouncepages.com".freeze
      RESPONSE_FINGERPRINT = "The requested URL was not found on this server.".freeze

      def run
        if apex_resource?
          return false unless resource_value == APEX_VALUE
        elsif cname_resource?
          return false unless resource_value == CNAME_VALUE
        end
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
