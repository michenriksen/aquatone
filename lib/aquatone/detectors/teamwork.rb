module Aquatone
  module Detectors
    class Teamwork < Aquatone::Detector
      self.meta = {
        :service         => "Teamwork",
        :service_website => "https://www.teamwork.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Project management, help desk and chat software"
      }

      CNAME_VALUE          = ".teamwork.com".freeze
      RESPONSE_FINGERPRINT = "<title>Oops - We didn't find your site.</title>".freeze

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
