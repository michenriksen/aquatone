module Aquatone
  module Collectors
    class Publicwww < Aquatone::Collector
      self.meta = {
        :name        => "PublicWWW",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses the publicwww.com source code search engine to find subdomains",
        :cli_options => {
          "publicwww-pages PAGES" => "Number of PublicWWW pages to process (default: 30)"
        }
      }

      DEFAULT_PAGES_TO_PROCESS = 30.freeze

      def run
        pages_to_process.times do |page|
          response = get_request("https://publicwww.com/websites/.#{url_escape(domain.name)}/#{page + 1}")
          response.body.gsub("<b>", "").gsub("</b>", "").scan(/[a-z0-9\.\-_]+\.#{regex_escape(domain.name)}/).each do |host|
            add_host(host)
          end
        end
      end

      private

      def pages_to_process
        if has_cli_option?("publicwww-pages")
          return get_cli_option("publicwww-pages").to_i
        end
        DEFAULT_PAGES_TO_PROCESS
      end
    end
  end
end
