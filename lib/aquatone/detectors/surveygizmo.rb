module Aquatone
  module Detectors
    class Surveygizmo < Aquatone::Detector
      self.meta = {
        :service         => "SurveyGizmo",
        :service_website => "https://www.surveygizmo.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Online survey software"
      }

      CNAME_VALUES          = %w(privatedomain.sgizmo.com privatedomain.surveygizmo.eu privatedomain.sgizmoca.com).freeze
      RESPONSE_FINGERPRINT = 'data-html-name="Header Logo Link"'.freeze

      def run
        return false unless cname_resource?
        CNAME_VALUES.each do |cname_value|
          if resource_value == cname_value
            return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
          end
        end
        false
      end
    end
  end
end
