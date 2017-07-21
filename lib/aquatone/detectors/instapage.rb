module Aquatone
  module Detectors
    class Instapage < Aquatone::Detector
      self.meta = {
        :service         => "Instapage",
        :service_website => "https://instapage.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Landing page platform"
      }

      CNAME_VALUES          = %w(pageserve.co secure.pageserve.co).freeze
      RESPONSE_FINGERPRINT = "You've Discovered A Missing Link. Our Apologies!".freeze

      def run
        return false unless cname_resource?
        return false unless CNAME_VALUES.include?(resource_value)
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
