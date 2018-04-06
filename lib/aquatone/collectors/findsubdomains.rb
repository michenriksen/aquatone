module Aquatone
  module Collectors
    class FindSubdomains < Aquatone::Collector
      self.meta = {
        :name         => "FindSubdomains",
        :author       => "Alessandro De Micheli (@eur0pa_)",
        :description  => "Uses findsubdomains.com to find hostnames"
      }

      def run
        response = get_request("https://findsubdomains.com/subdomains-of/#{url_escape(domain.name)}")

        response.body.to_enum(:scan, /<a\s.*?>\s+?([a-zA-Z0-9\*_.-]+\.#{Regexp.escape(domain.name)})<\/a><\/h4>/).map do |column|
          add_host(column[0].gsub("*.", ""))
        end
      end
    end
  end
end
