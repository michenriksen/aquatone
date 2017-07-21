module Aquatone
  module Detectors
    class Zendesk < Aquatone::Detector
      self.meta = {
        :service         => "Zendesk",
        :service_website => "https://www.zendesk.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Customer service software and support ticket system"
      }

      CNAME_VALUE          = ".zendesk.com".freeze
      RESPONSE_FINGERPRINT = "<title>Help Center Closed | Zendesk</title>".freeze

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
