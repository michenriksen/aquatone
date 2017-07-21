module Aquatone
  module Detectors
    class GithubPages < Aquatone::Detector
      self.meta = {
        :service         => "GitHub Pages",
        :service_website => "https://pages.github.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "GitHub static website hosting"
      }

      APEX_VALUES          = %w(192.30.252.153 192.30.252.154).freeze
      CNAME_VALUE          = ".github.io".freeze
      RESPONSE_FINGERPRINT = "There isn't a GitHub Pages site here.".freeze

      def run
        if apex_resource?
          return false unless APEX_VALUES.include?(resource_value)
        elsif cname_resource?
          return false unless resource_value.end_with?(CNAME_VALUE)
        end
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
