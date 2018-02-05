module Aquatone
  module Detectors
    class Tilda < Aquatone::Detector
      self.meta = {
        :service         => "Tilda",
        :service_website => "https://tida.ws/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "Web content hosting"
      }

      APEX_VALUE           = "178.248.234.146"
      CNAME_VALUE          = "tilda.ws".freeze
      RESPONSE_FINGERPRINT = "Domain has been assigned".freeze

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
