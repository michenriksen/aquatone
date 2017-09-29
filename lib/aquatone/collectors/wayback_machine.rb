require 'uri'

module Aquatone
  module Collectors
    class WaybackMachine < Aquatone::Collector
      self.meta = {
        :name         => "Wayback Machine",
        :author       => "Joel (@jolle)",
        :description  => "Uses Wayback Machine by Internet Archive to find unique hostnames"
      }

      def run
        response = get_request("http://web.archive.org/cdx/search/cdx?url=*.#{url_escape(domain.name)}&output=json&fl=original&collapse=urlkey")

        response.parsed_response.each do |page|
          if page[0] != "original"
              add_host(URI.parse(page[0]).host)
          end
        end
      end
    end
  end
end
