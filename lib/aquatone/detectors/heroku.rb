module Aquatone
  module Detectors
    class Heroku < Aquatone::Detector
      self.meta = {
        :service         => "Heroku",
        :service_website => "https://www.heroku.com/",
        :author          => "Michael Henriksen (@michenriksen)",
        :description     => "Cloud application platform"
      }

      CNAME_VALUES         = %w(.herokudns.com .herokussl.com .herokuapp.com).freeze
      RESPONSE_FINGERPRINT = "<title>No such app</title>".freeze

      def run
        return false unless cname_resource?
        CNAME_VALUES.each do |cname_value|
          if resource_value.end_with?(cname_value)
            return get_request("http://#{host}/").body.include?(RESPONSE_FINGERPRINT)
          end
        end
        false
      end
    end
  end
end
