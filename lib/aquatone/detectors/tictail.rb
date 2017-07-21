module Aquatone
  module Detectors
    class Tictail < Aquatone::Detector
      self.meta = {
        :service         => "Tictail",
        :service_website => "https://tictail.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Social shopping platform"
      }

      APEX_VALUE           = "46.137.181.142".freeze
      CNAME_VALUE          = "domains.tictail.com".freeze
      RESPONSE_FINGERPRINT = 'class="MarketplaceHeader__tictailLogo"'.freeze

      def run
        if apex_resource?
          return false unless resource_value == APEX_VALUE
        elsif cname_resource?
          return false unless resource_value == CNAME_VALUE
        end
        get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
      end
    end
  end
end
