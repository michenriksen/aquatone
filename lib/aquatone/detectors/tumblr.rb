module Aquatone
  module Detectors
    class Tumblr < Aquatone::Detector
      self.meta = {
        :service         => "Tumblr",
        :service_website => "https://www.tumblr.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Microblogging and social networking platform"
      }

      APEX_VALUE           = "66.6.44.4".freeze
      CNAME_VALUE          = "domains.tumblr.com".freeze
      RESPONSE_FINGERPRINT = "Whatever you were looking for doesn't currently exist at this address.".freeze

      def run
        if apex_resource?
          return false unless resource_value == APEX_VALUE
        elsif cname_resource?
          return false unless resource_value == CNAME_VALUE
        end
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
