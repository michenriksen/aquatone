module Aquatone
  module Collectors
    class Dictionary < Aquatone::Collector
      self.meta = {
        :name        => "Dictionary",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses a dictionary to find hostnames"
      }

      DICTIONARY = File.join(Aquatone::AQUATONE_ROOT, "subdomains.lst").freeze

      def run
        dictionary = File.open(DICTIONARY, "r")
        dictionary.each_line do |subdomain|
          add_host("#{subdomain.strip}.#{domain.name}")
        end
      end
    end
  end
end
