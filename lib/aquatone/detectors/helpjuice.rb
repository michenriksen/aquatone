module Aquatone
  module Detectors
    class Helpjuice < Aquatone::Detector
      self.meta = {
        :service         => "Helpjuice",
        :service_website => "https://helpjuice.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Knowledge base software"
      }

      CNAME_VALUE          = ".helpjuice.com".freeze
      RESPONSE_FINGERPRINT = "<title>No such app</title>".freeze

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
