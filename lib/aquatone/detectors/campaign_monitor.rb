module Aquatone
  module Detectors
    class CampaignMonitor < Aquatone::Detector
      self.meta = {
        :service         => "Campaign Monitor",
        :service_website => "https://www.zendesk.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Email marketing"
      }

      CNAME_VALUE          = "name.createsend.com".freeze
      RESPONSE_FINGERPRINT = "<strong>Trying to access your account?</strong>".freeze

      def run
        return false unless cname_resource?
        if resource_value == CNAME_VALUE
          return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
        end
        false
      end
    end
  end
end
