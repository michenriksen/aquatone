module Aquatone
  module Detectors
    class S3 < Aquatone::Detector
      self.meta = {
        :service         => "Amazon S3",
        :service_website => "https://aws.amazon.com/s3/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Cloud storage"
      }

      CNAME_VALUE          = ".amazonaws.com".freeze
      RESPONSE_FINGERPRINT = "NoSuchBucket".freeze

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
