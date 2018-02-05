module Aquatone
  module Detectors
    class BitBucket < Aquatone::Detector
      self.meta = {
        :service         => "BitBucket",
        :service_website => "https://bitbucket.org/",
        :author          => "Alessandro De Micheli (@eur0pa_)",
        :description     => "Bitbucket static page hosting"
      }

      CNAME_VALUE          = "bitbucket.org".freeze
      RESPONSE_FINGERPRINT = "The page you have requested does not exist".freeze

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
