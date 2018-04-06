module Aquatone
  module Detectors
    class Pantheon < Aquatone::Detector
      self.meta = {
        :service         => "Pantheon",
        :service_website => "https://pantheonsite.io/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "Web content hosting"
      }

      CNAME_VALUE          = "pantheonsite.io".freeze
      RESPONSE_FINGERPRINT = "The gods are wise".freeze

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
