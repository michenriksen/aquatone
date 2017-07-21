module Aquatone
  module Detectors
    class Feedpress < Aquatone::Detector
      self.meta = {
        :service         => "FeedPress",
        :service_website => "https://feed.press/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Feed analytics and Podcast hosting"
      }

      CNAME_VALUE          = "redirect.feedpress.me".freeze
      RESPONSE_FINGERPRINT = "The feed has not been found.".freeze

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
