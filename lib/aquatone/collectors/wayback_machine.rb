require 'uri'

module Aquatone
  module Collectors
    class WaybackMachine < Aquatone::Collector
      self.meta = {
        :name         => "Wayback Machine",
        :author       => "Joel (@jolle)",
        :description  => "Uses Wayback Machine by Internet Archive to find unique hostnames",
        :cli_options  => {
          "wayback-machine-timeout SECONDS" => "Timeout for Wayback Machine collector in seconds (default: 30)",
          "wayback-machine-page-timeout SECONDS" => "Timeout for one page fetch of the Wayback Machine collector (default: none)",
          "wayback-machine-pages AMOUNT" => "Amount of pages to fetch for Wayback Machine (amount in numbers or \"all\") (default: 3)"
        }
      }

      DEFAULT_TIMEOUT       = 30.freeze
      DEFAULT_PAGES         = 3.freeze
      DEFAULT_PAGE_TIMEOUT  = nil.freeze

      def run
        Timeout::timeout(timeout) do
          base_uri = "http://web.archive.org/cdx/search/cdx?url=#{url_escape(domain.name)}&matchType=domain&output=json".freeze

          pages_response = nil
          Timeout::timeout(page_timeout) do
            pages_response = get_request("#{base_uri}&showNumPages=true")
          end

          available_pages = pages_response.body.to_i
          fetchable_pages = pages

          if fetchable_pages == 0 && available_pages > 70
            print "\e[1m\e[33mWayback Machine found #{available_pages} pages of entries. Consider lowering the amount of pages with the wayback-machine-pages option.\e[0m"
          end

          if fetchable_pages > available_pages || fetchable_pages == 0
            fetchable_pages = available_pages
          end

          fetchable_pages.times do |page|
            response = nil
            Timeout::timeout(page_timeout) do
              response = get_request("#{base_uri}&page=#{page}&fl=original&collapse=urlkey")
            end

            response.parsed_response.each do |page|
              if page[0] != "original"
                begin
                  add_host(URI.parse(page[0]).host)
                rescue URI::Error; end
              end
            end
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

      def page_timeout
        if has_cli_option?("wayback-machine-page-timeout")
          return get_cli_option("wayback-machine-page-timeout").to_i
        end
        DEFAULT_PAGE_TIMEOUT
      end

      def pages
        if has_cli_option?("wayback-machine-pages")
          pages_option = get_cli_option("wayback-machine-pages")

          if pages_option == "all"
            return 0
          else
            return pages_option.to_i
          end
        end
        DEFAULT_PAGES
      end
    end
  end
end
