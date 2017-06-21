module Aquatone
  module Collectors
    class Censys < Aquatone::Collector
      self.meta = {
        :name         => "Censys",
        :author       => "James McLean (vortexau)",
        :description  => "Uses the Censys API to find hostnames in TLS certificates",
        :require_keys => ["censys"]
        :require_id   => ["censys"]
      }

      API_BASE_URI         = "https://www.censys.io/api/v1".freeze
      API_RESULTS_PER_PAGE = 100.freeze
      PAGES_TO_PROCESS     = 10.freeze

      def run
        request_censys_page
      end

      private

      def request_censys_page()



      end
    end
  end
end
