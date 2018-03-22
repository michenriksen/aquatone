module Aquatone
  module Collectors
    class Ptrarchive < Aquatone::Collector
      self.meta = {
        :name        => "PTRArchive",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses ptrarchive.com to find subdomains"
      }

      def run
        response = get_request("http://ptrarchive.com/tools/search.htm?label=#{url_escape(domain.name)}&date=ALL", 
          :headers => {
            "User-Agent" => "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
          }
        )
        if response.code != 200
          failure("PTRArchive returned unexpected response code: #{response.code}")
        end
        response.body.scan(/[a-z0-9\.\-_]+\.#{regex_escape(domain.name)}/).each do |host|
          add_host(host)
        end
      end
    end
  end
end
