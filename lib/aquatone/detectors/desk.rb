module Aquatone
  module Detectors
    class Desk < Aquatone::Detector
      self.meta = {
        :service         => "Desk",
        :service_website => "https://www.desk.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Customer service and helpdesk ticket software"
      }

      CNAME_VALUE          = ".desk.com".freeze
      RESPONSE_FINGERPRINT = "Sorry, We Couldn't Find That Page".freeze

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
