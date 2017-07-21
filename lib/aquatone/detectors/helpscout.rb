module Aquatone
  module Detectors
    class Helpscout < Aquatone::Detector
      self.meta = {
        :service         => "Help Scout",
        :service_website => "https://www.helpscout.net/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Customer service software and education platform"
      }

      CNAME_VALUE          = ".helpscoutdocs.com".freeze
      RESPONSE_FINGERPRINT = "No settings were found for this company".freeze

      def run
        return false unless cname_resource?
        if resource_value.end_with?(CNAME_VALUE)
          return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
        end
        false
      end
    end
  end
end
