module Aquatone
  module Detectors
    class Uservoice < Aquatone::Detector
      self.meta = {
        :service         => "UserVoice",
        :service_website => "https://www.uservoice.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Product management software"
      }

      CNAME_VALUE           = ".uservoice.com".freeze
      RESPONSE_FINGERPRINTS = [
        "The page you have requested does not exist.",
        "This UserVoice subdomain is currently available!"
      ].freeze

      def run
        return false unless cname_resource?
        if resource_value.end_with?(CNAME_VALUE)
          response = get_request("http://#{host}/")
          RESPONSE_FINGERPRINTS.each do |fingerprint|
            return true if response.body.include?(fingerprint)
          end
        end
        false
      end
    end
  end
end
