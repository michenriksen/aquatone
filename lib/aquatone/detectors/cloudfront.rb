module Aquatone
  module Detectors
    class Cloudfront < Aquatone::Detector
      self.meta = {
        :service         => "Cloudfront",
        :service_website => "https://aws.amazon.com/cloudfront/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Content delivery network"
      }

      CNAME_VALUE          = ".cloudfront.net".freeze
      RESPONSE_FINGERPRINT = "The request could not be satisfied".freeze

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
