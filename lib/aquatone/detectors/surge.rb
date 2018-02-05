module Aquatone
  module Detectors
    class Surge < Aquatone::Detector
      self.meta = {
        :service         => "Surge",
        :service_website => "https://surge.sh/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "Project hosting"
      }

      APEX_VALUE           = "138.197.235.123"
      CNAME_VALUE          = "tilda.sh".freeze
      RESPONSE_FINGERPRINT = "project not found".freeze

      def run
        if apex_resource?
          return false unless APEX_VALUES.include?(resource_value)
        elsif cname_resource?
          return false unless resource_value.end_with?(CNAME_VALUE)
        end
        return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
