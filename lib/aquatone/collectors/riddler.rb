module Aquatone
  module Collectors
    class Riddler < Aquatone::Collector
      self.meta = {
        :name         => "Riddler",
        :author       => "Joel (@jolle)",
        :description  => "Uses Riddler by F-Secure to find hostnames",
        :require_keys => ["riddler_username", "riddler_password"]
      }

      API_BASE_URI = "https://riddler.io"

      def run
        auth_response = post_request("#{API_BASE_URI}/auth/login", {
          :email    => get_key("riddler_username"),
          :password => get_key("riddler_password")
        }.to_json, {
          :headers  => { "Content-Type" => "application/json" }
        })

        if auth_response.code == 400 or auth_response.parsed_response["meta"]["code"] == 400
          failure("Invalid credentials to Riddler.io")
        elsif auth_response.code != 200 or auth_response.parsed_response["meta"]["code"] != 200
          failure("Riddler.io auth API returned unexpected response code: #{response.code}")
        end
        
        token = auth_response.parsed_response["response"]["user"]["authentication_token"]

        response = post_request("#{API_BASE_URI}/api/search", {
          :query  => "pld:#{url_escape(domain.name)}",
          :output => "host",
          :limit  => 500
        }.to_json, {
          :headers  => { "Content-Type" => "application/json", "Authentication-Token" => token }
        })

        response.parsed_response.each do |record|
          add_host(record["host"]) unless record["host"].empty?
        end
      end
    end
  end
end