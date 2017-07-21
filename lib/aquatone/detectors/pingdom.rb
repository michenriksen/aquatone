module Aquatone
  module Detectors
    class Pingdom < Aquatone::Detector
      self.meta = {
        :service         => "Pingdom",
        :service_website => "https://www.pingdom.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Website and performance monitoring"
      }

      CNAME_VALUE          = "stats.pingdom.com".freeze
      RESPONSE_FINGERPRINT = "Sorry, couldn&rsquo;t find the status page".freeze

      def run
        return false unless cname_resource?
        return false unless resource_value == CNAME_VALUE
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
