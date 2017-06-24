module Aquatone
  module Collectors
    class Crtsh < Aquatone::Collector
      self.meta = {
        :name         => "crtsh",
        :author       => "Joel (@jolle)",
        :description  => "Uses crt.sh by COMODO CA to find hostnames"
      }

      def run
        response = get_request("https://crt.sh/?dNSName=%25.#{url_escape(domain.name)}")

        parsed = parse_html(response.body)

        parsed.css('table').css('table').css('tr').css('td:nth-child(4)').map do |column|
          add_host(column.text.gsub("*.", ""))
        end
      end
    end
  end
end