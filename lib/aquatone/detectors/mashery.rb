module Aquatone
  module Detectors
    class Mashery < Aquatone::Detector
      self.meta = {
        :service         => "Mashery",
        :service_website => "https://www.mashery.com/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "Content hosting"
      }

      CNAME_VALUE          = "wildcard-mashery-com".freeze
      RESPONSE_FINGERPRINT = "Unrecognized domain".freeze

      def run
        return false unless cname_resource?
        if resource_value.includes?(CNAME_VALUE)
          return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
        end
        false
      end
    end
  end
end
