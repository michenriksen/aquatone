module Aquatone
  module Detectors
    class Smartling < Aquatone::Detector
      self.meta = {
        :service         => "Smartling",
        :service_website => "https://smartling.com/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "Content translation and localization"
      }

      CNAME_VALUE          = "smartling.com".freeze
      RESPONSE_FINGERPRINT = "Domain is not configured".freeze

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
