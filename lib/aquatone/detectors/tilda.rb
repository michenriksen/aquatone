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
      CNAME_VALUES         = %w(tilda.sh tilda.ws).freeze
      RESPONSE_FINGERPRINT = "Domain has been assigned".freeze

      def run
        if apex_resource?
          return false unless resource_value == APEX_VALUE
        elsif cname_resource?
          return false unless CNAME_VALUES.include?(resource_value)
        end
        return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
