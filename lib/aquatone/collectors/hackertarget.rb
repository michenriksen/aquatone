module Aquatone
  module Collectors
    class Hackertarget < Aquatone::Collector
      self.meta = {
        :name        => "HackerTarget",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses hackertarget.com to find hostnames"
      }

      API_BASE_URI = "https://api.hackertarget.com"

      def run
        response = get_request("#{API_BASE_URI}/hostsearch/?q=#{url_escape(domain.name)}")
        if response.code != 200
          failure("HackerTarget API returned unexpected response code: #{response.code}")
        end
        response.body.each_line do |line|
          host = line.split(",", 2).first.strip
          add_host(host) unless host.empty?
        end
      end
    end
  end
end
