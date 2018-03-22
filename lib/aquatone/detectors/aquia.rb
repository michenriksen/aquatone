module Aquatone
  module Detectors
    class Aquia < Aquatone::Detector
      self.meta = {
        :service         => "Aquia",
        :service_website => "https://aquia.com/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "Community and e-commerce platform"
      }

      APEX_VALUE           = "69.172.201.153"
      CNAME_VALUE          = "aquia.com".freeze
      RESPONSE_FINGERPRINT = "If you are an Acquia Cloud customer and expect to see your site at this address".freeze

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
