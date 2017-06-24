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

        response.body.to_enum(:scan, /<TD>([a-zA-Z0-9\*_.-]+\.#{domain.name})<\/TD>/).map do |column|
          add_host(column[0].gsub("*.", ""))
        end
      end
    end
  end
end