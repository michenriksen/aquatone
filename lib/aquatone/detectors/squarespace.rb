module Aquatone
  module Detectors
    class Squarespace < Aquatone::Detector
      self.meta = {
        :service         => "Squarespace",
        :service_website => "https://www.squarespace.com/",
        :author          => "Duarte Duarte (@dduarte)",
        :description     => "Website builder"
      }

      RESPONSE_FINGERPRINT = "Squarespace - Claim This Domain".freeze

      def run
        return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
