module Aquatone
  module Collectors
    class Netcraft < Aquatone::Collector
      self.meta = {
        :name        => "Netcraft",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses searchdns.netcraft.com to find hostnames"
      }

      BASE_URI         = "http://searchdns.netcraft.com/".freeze
      HOSTNAME_REGEX   = /<a href="http:\/\/(.*?)\/" rel="nofollow">/.freeze
      RESULTS_PER_PAGE = 20.freeze
      PAGES_TO_PROCESS = 10.freeze

      def run
        last  = nil
        count = 0
        PAGES_TO_PROCESS.times do |i|
          page = i + 1
          if page == 1
            uri = "#{BASE_URI}/?restriction=site+contains&host=*.#{url_escape(domain.name)}&lookup=wait..&position=limited"
          else
            uri = "#{BASE_URI}/?host=*.#{url_escape(domain.name)}&last=#{url_escape(last)}&from=#{count + 1}&restriction=site%20contains&position=limited"
          end
          response = get_request(uri,
            { :headers => { "Referer" => "http://searchdns.netcraft.com/" } }
          )
          hosts = extract_hostnames_from_response(response.body)
          last  = hosts.last
          count += hosts.count
          hosts.each { |host| add_host(host) }
          break if hosts.count != RESULTS_PER_PAGE
          random_sleep(5)
        end
      end

      private

      def extract_hostnames_from_response(body)
        hosts = []
        body.scan(HOSTNAME_REGEX).each do |match|
          hosts << match.last.to_s.strip.downcase
        end
        hosts
      end
    end
  end
end
