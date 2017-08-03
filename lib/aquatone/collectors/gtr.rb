module Aquatone
  module Collectors
    class Gtr < Aquatone::Collector
      self.meta = {
        :name        => "Google Transparency Report",
        :author      => "Michael Henriksen (@michenriksen)",
        :description => "Uses Google Transparency Report to find hostnames",
        :slug        => "gtr",
        :cli_options  => {
          "gtr-pages PAGES" => "Number of Google Transparency Report pages to process (default: 30)"
        }
      }

      BASE_URI                 = "https://www.google.com/transparencyreport/jsonp/ct/search"
      DEFAULT_PAGES_TO_PROCESS = 30.freeze

      def run
        token = nil
        pages_to_process.times do
          response = parse_response(request_page(token))
          response["results"].each do |result|
            host = result["subject"]
            add_host(host) if valid_host?(host)
          end
          break if !response.key?("nextPageToken")
          token = response["nextPageToken"]
        end
      end

      private

      def request_page(token = nil)
        if token.nil?
          uri = "#{BASE_URI}?domain=#{url_escape(domain.name)}&incl_exp=true&incl_sub=true&c=_callbacks_._#{random_jsonp_callback}"
        else
          uri = "#{BASE_URI}?domain=#{url_escape(domain.name)}&incl_exp=true&incl_sub=true&token=#{url_escape(token)}&c=_callbacks_._#{random_jsonp_callback}"
        end

        get_request(uri,
          { :headers => { "Referer" => "https://www.google.com/transparencyreport/https/ct/?hl=en-US" } }
        )
      end

      def random_jsonp_callback
        "abcdefghijklmnopqrstuvwxyz0123456789".split("").sample(9).join
      end

      def parse_response(body)
        body = body.split("(", 2).last
        body.gsub!(");", "")
        JSON.parse(body)
      end

      def valid_host?(host)
        return false if host.start_with?("*.")
        return false unless host.end_with?(".#{domain.name}")
        true
      end

      def pages_to_process
        if has_cli_option?("gtr-pages")
          return get_cli_option("gtr-pages").to_i
        end
        DEFAULT_PAGES_TO_PROCESS
      end
    end
  end
end
