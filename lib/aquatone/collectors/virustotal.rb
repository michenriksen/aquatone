module Aquatone
  module Collectors
    class Virustotal < Aquatone::Collector
      self.meta = {
        :name         => "VirusTotal",
        :author       => "Michael Henriksen (@michenriksen)",
        :description  => "Uses virustotal.com domain search to find hostnames",
        :require_keys => ["virustotal"]
      }

      API_URI = "http://www.virustotal.com/vtapi/v2/domain/report".freeze

      def run
        response = get_request("#{API_URI}?domain=#{url_escape(domain.name)}&apikey=#{get_key('virustotal')}")
        if response.code != 200
          failure("VirusTotal API returned unexpected status code: #{response.code}")
        end
        if response.parsed_response.key?("subdomains")
          response.parsed_response["subdomains"].each { |host| add_host(host) }
        end
      end
    end
  end
end
