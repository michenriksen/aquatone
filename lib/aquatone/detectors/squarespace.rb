module Aquatone
  module Detectors
    class SquareSpace < Aquatone::Detector
      self.meta = {
        :service         => "Squarespace",
        :service_website => "https://squarespace.com/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "Website hosting"
      }

      CNAME_VALUE          = ".squarespace.com".freeze
      RESPONSE_FINGERPRINT = "No Such Account".freeze

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
