module Aquatone
  module Collectors
    class Dictionary < Aquatone::Collector
      self.meta = {
        :name        => "Dictionary",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses a dictionary to find hostnames",
        :cli_options => {
          "wordlist WORDLIST" => "OPTIONAL: wordlist/dictionary file to use for subdomain bruteforcing"
        }
      }

      DEFAULT_DICTIONARY = File.join(Aquatone::AQUATONE_ROOT, "subdomains.lst").freeze

      def run
        if has_cli_option?("wordlist")
          file = File.expand_path(get_cli_option("wordlist"))
          if !File.readable?(file)
            failure("Wordlist file #{file} is not readable or does not exist")
          end
          dictionary = File.open(file, "r")
        else
          dictionary = File.open(DEFAULT_DICTIONARY, "r")
        end

        dictionary.each_line do |subdomain|
          add_host("#{subdomain.strip}.#{domain.name}")
        end
      end
    end
  end
end
