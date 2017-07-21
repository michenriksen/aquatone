module Aquatone
  module Detectors
    class Freshdesk < Aquatone::Detector
      self.meta = {
        :service         => "Freshdesk",
        :service_website => "https://freshdesk.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Customer support software and ticketing system"
      }

      CNAME_VALUE          = ".freshdesk.com".freeze
      RESPONSE_FINGERPRINT = "You can claim it now at".freeze

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
