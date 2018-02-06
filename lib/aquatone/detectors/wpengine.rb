module Aquatone
  module Detectors
    class Wpengine < Aquatone::Detector
      self.meta = {
        :service         => "WPEngine",
        :service_website => "https://wpengine.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "WordPress blog hosting"
      }

      APEX_VALUE           = "130.211.160.56".freeze
      CNAME_VALUE          = ".wpengine.com".freeze
      RESPONSE_FINGERPRINT = "The site you were looking for is no longer available at this IP address".freeze

      def run
        if apex_resource?
          return false unless resource_value == APEX_VALUE
        elsif cname_resource?
          return false unless resource_value.end_with?(CNAME_VALUE)
        end
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
