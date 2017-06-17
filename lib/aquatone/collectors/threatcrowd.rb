module Aquatone
  module Collectors
    class Threatcrowd < Aquatone::Collector
      self.meta = {
        :name         => "Threat Crowd",
        :author       => "Michael Henriksen (@michenriksen)",
        :description  => "Uses threadcrowd.org API to find hostnames",
        :slug         => "threatcrowd"
      }

      API_URI = "https://www.threatcrowd.org/searchApi/v2/domain/report/".freeze

      def run
        response = get_request("#{API_URI}?domain=#{url_escape(domain.name)}")
        if response.code != 200
          failure("Threat Crowd API returned unexpected status code: #{response.code}")
        end
        body = JSON.parse(response.body)
        if body.key?("subdomains")
          body["subdomains"].each { |host| add_host(host) }
        end
      end
    end
  end
end
