module Aquatone
  module Collectors
    class ThreatMiner < Aquatone::Collector
      self.meta = {
        :name         => "ThreatMiner",
        :author       => "Joel (@jolle)",
        :description  => "Uses ThreatMiner to find hostnames"
      }

      def run
        response = get_request("https://www.threatminer.org/getData.php?e=subdomains_container&q=#{url_escape(domain.name)}&t=0&rt=10&p=1")

        if response.code != 200
          failure("ThreatMiner returned an unexpected response code: #{response.code}")
        end

        response.body.to_enum(:scan, /"domain\.php\?q=([a-zA-Z0-9\*_.-]+\.#{Regexp.escape(domain.name)})"/).map do |link|
          add_host(link[0])
        end
      end
    end
  end
end
