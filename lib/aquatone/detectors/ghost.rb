module Aquatone
  module Detectors
    class Ghost < Aquatone::Detector
      self.meta = {
        :service         => "Ghost",
        :service_website => "https://ghost.org/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Publishing platform"
      }

      CNAME_VALUE          = ".ghost.io".freeze
      RESPONSE_FINGERPRINT = "The thing you were looking for is no longer here, or never was".freeze

      def run
        return false unless cname_resource?
        if resource_value.end_with?(CNAME_VALUE)
          response = get_request("http://#{host}/",
            # Set a proper User-Agent to avoid potential CloudFlare CAPTCHA wall
            :headers => {
              "User-Agent" => "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
            }
          )
          return response.body.include?(RESPONSE_FINGERPRINT)
        end
        false
      end
    end
  end
end
