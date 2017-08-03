module Aquatone
  module Collectors
    class Censys < Aquatone::Collector
      self.meta = {
        :name         => "Censys",
        :author       => "James McLean (@vortexau)",
        :description  => "Uses the Censys API to find hostnames in TLS certificates",
        :require_keys => ["censys_secret","censys_id"],
        :cli_options  => {
          "censys-pages PAGES" => "Number of Censys API pages to process (default: 10)"
        }
      }

      API_BASE_URI             = "https://www.censys.io/api/v1".freeze
      API_RESULTS_PER_PAGE     = 100.freeze
      DEFAULT_PAGES_TO_PROCESS = 10.freeze

      def run
        request_censys_page
      end

      def request_censys_page(page=1)
          # Initial version only supporting Censys Certificates API

          # Censys expects Basic Auth for requests.
          auth = {
              :username => get_key('censys_id'),
              :password => get_key('censys_secret')
          }

          # Define this is JSON content
          headers = {
            'Content-Type' => 'application/json',
            'Accept' => 'application/json'
          }

          # The post body itself, as JSON
          query = {
              'query'   => url_escape("#{domain.name}"),
              'page'    => page,
              'fields'  => [ "parsed.names", "parsed.extensions.subject_alt_name.dns_names" ],
              'flatten' => true
          }

          # Search API documented at https://censys.io/api/v1/docs/search
          response = post_request(
              "#{API_BASE_URI}/search/certificates",
              query.to_json,
              {
                  :basic_auth => auth,
                  :headers => headers
              }
          )

          if response.code != 200
              failure(response.parsed_response["error"] || "Censys API encountered error: #{response.code}")
          end

          # If nothing returned from Censys, return:
          return unless response.parsed_response["results"]

          response.parsed_response["results"].each do |result|

            next unless result["parsed.extensions.subject_alt_name.dns_names"]
            result["parsed.extensions.subject_alt_name.dns_names"].each do |altdns|
                add_host(altdns) if altdns.end_with?(".#{domain.name}")
            end

            next unless result["parsed.names"]
            result["parsed.names"].each do |parsedname|
                add_host(parsedname) if parsedname.end_with?(".#{domain.name}")
            end
          end

          # Get the next page of results
          request_censys_page(page + 1) if next_page?(page, response.parsed_response)

      end

      def next_page?(page, body)
          page <= pages_to_process && body["metadata"]["pages"] && API_RESULTS_PER_PAGE * page < body["metadata"]["count"].to_i
      end

      def pages_to_process
        if has_cli_option?("censys-pages")
          return get_cli_option("censys-pages").to_i
        end
        DEFAULT_PAGES_TO_PROCESS
      end
    end
  end
end
