module Aquatone
  module Detectors
    class Statuspage < Aquatone::Detector
      self.meta = {
        :service         => "StatusPage",
        :service_website => "https://www.statuspage.io/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Status page hosting"
      }

      CNAME_VALUE          = ".stspg-customer.com".freeze
      RESPONSE_FINGERPRINT = "<title>Hosted Status Pages for Your Company</title>".freeze

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
