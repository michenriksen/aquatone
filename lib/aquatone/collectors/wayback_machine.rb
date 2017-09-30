require 'uri'

module Aquatone
  module Collectors
    class WaybackMachine < Aquatone::Collector
      self.meta = {
        :name         => "Wayback Machine",
        :author       => "Joel (@jolle)",
        :description  => "Uses Wayback Machine by Internet Archive to find unique hostnames",
        :cli_options  => {
          "wayback-machine-timeout SECONDS" => "Timeout for Wayback Machine collector in seconds (default: 30)"
        }
      }

      DEFAULT_TIMEOUT = 30.freeze

      def run
        response = nil
        Timeout::timeout(timeout) do
          response = get_request("http://web.archive.org/cdx/search/cdx?url=*.#{url_escape(domain.name)}&output=json&fl=original&collapse=urlkey")
        end
        response.parsed_response.each do |page|
          if page[0] != "original"
            begin
              add_host(URI.parse(page[0]).host)
            rescue URI::Error; end
          end
        end
      end

      private

      def timeout
        if has_cli_option?("wayback-machine-timeout")
          return get_cli_option("wayback-machine-timeout").to_i
        end
        DEFAULT_TIMEOUT
      end
    end
  end
end
