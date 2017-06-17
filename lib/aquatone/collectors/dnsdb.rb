module Aquatone
  module Collectors
    class Dnsdb < Aquatone::Collector
      self.meta = {
        :name        => "DNSDB",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses dnsdb.org to find hostnames"
      }

      LINK_REGEX = /<a href="(.*?)\">(.*?)<\/a>/.freeze
      INDEX_REGEX = /<a href="([0-9a-z]{1})">[0-9a-z]{1}<\/a>/.freeze

      def run
        @base_url = "http://www.dnsdb.org/f/#{url_escape(domain.name)}.dnsdb.org/"
        parse_page(@base_url)
      end

      private

      def parse_page(url)
        response = get_request(url)
        if response.code != 200
          failure("DNSDB returned unexpected response code: #{response.code}")
        end

        if response.body.include?("index for")
          response.body.scan(INDEX_REGEX) do |index|
            response = get_request("#{@base_url}#{url_escape(index.first)}")
            extract_hosts(response.body)
          end
        else
          extract_hosts(response.body)
        end
      end

      def extract_hosts(body)
        body.scan(LINK_REGEX) do |href, hostname|
          if hostname.end_with?(".#{domain.name}")
            add_host(hostname)
          end
        end
      end
    end
  end
end
