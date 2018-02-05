module Aquatone
  module Detectors
    class Shopify < Aquatone::Detector
      self.meta = {
        :service         => "Shopify",
        :service_website => "https://www.shopify.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Ecommerce platform"
      }

      APEX_VALUE           = "23.227.38.32".freeze
      CNAME_VALUE          = "shops.myshopify.com".freeze
      RESPONSE_FINGERPRINT = "Sorry, this shop is currently unavailable.".freeze

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
