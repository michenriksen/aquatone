module Aquatone
  module Detectors
    class Edgecast < Aquatone::Detector
      self.meta = {
        :service         => "Edgecast",
        :service_website => "https://www.verizondigitalmedia.com/platform/edgecast-cdn/",
        :author          => "Duarte Duarte (@dduarte)",
        :description     => "Content delivery network"
      }

      CNAME_VALUE          = ".edgecastcdn.net".freeze
      RESPONSE_FINGERPRINT = "404 - Not Found".freeze

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
