module Aquatone
  module Collectors
    class Ptrarchive < Aquatone::Collector
      self.meta = {
        :name        => "PTRArchive",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses ptrarchive.com to find subdomains"
      }

      def run
        response = get_request("http://ptrarchive.com/tools/search3.htm?label=#{url_escape(domain.name)}&date=ALL",
         :headers => {
           "User-Agent" => "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10.6; en-US; rv:1.9.2.3) Gecko/20100402 Prism/1.0b4"
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
