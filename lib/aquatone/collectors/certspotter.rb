module Aquatone
  module Collectors
    class CertSpotter < Aquatone::Collector
      self.meta = {
        :name         => "Cert Spotter",
        :author       => "Joel (@jolle)",
        :description  => "Uses Cert Spotter by SSLMate to find hostnames"
      }

      def run
        response = get_request("https://certspotter.com/api/v0/certs?domain=#{url_escape(domain.name)}")

        if response.code != 200
          failure("Cert Spotter API returned an unexpected response code: #{response.code}")
        end

        response.parsed_response.each do |cert|
          cert['dns_names'].each do |name|
            if /\.#{Regexp.escape(domain.name)}$/.match(name) or name == domain.name
              add_host(name)
            end
          end
        end
      end
    end
  end
end
