module Aquatone
  module Detectors
    class Wordpress < Aquatone::Detector
      self.meta = {
        :service         => "Wordpress",
        :service_website => "https://wordpress.com/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "WordPress blog hosting"
      }

      APEX_VALUES          = %w(192.0.78.12 192.0.78.13).freeze
      CNAME_VALUE          = "lb.wordpress.com".freeze
      RESPONSE_FINGERPRINT = "Domain mapping upgrade for this domain not found".freeze

      def run
        if apex_resource?
          return false unless APEX_VALUES.include?(resource_value)
        elsif cname_resource?
          return false unless resource_value.end_with?(CNAME_VALUE)
        end
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
