require 'open3'

module Aquatone
  class Browser
    module Drivers

      class IncompatabilityError < StandardError; end

      class Nightmare
        attr_reader :url, :vhost, :html_destination, :headers_destination, :screenshot_destination, :options

        def initialize(url, vhost, html_destination, headers_destination, screenshot_destination, options)
          @url                    = url
          @vhost                  = vhost
          @html_destination       = html_destination
          @headers_destination    = headers_destination
          @screenshot_destination = screenshot_destination
          @options                = options
        end

        def visit
          Open3.popen3(construct_command) do |_, stdout, stderr, wait_thr|
            output, pid = [], wait_thr.pid
            begin
              Timeout.timeout(options[:timeout]) do
                output = [stdout.read, stderr.read]
                Process.wait(pid)
              end
            rescue Errno::ECHILD => e
            rescue Timeout::Error => e
              Process.kill('HUP', pid)
              return {
                "success" => false,
                "error"   => e.is_a?(Timeout::Error) ? "Timeout" : "#{e.class}: #{e.message}",
                "code"    => 0,
                "details" => ""
              }
            end
            JSON.parse(output[0])
          end
          rescue JSON::ParserError
            fail IncompatabilityError, "Nightmarejs must be run on a system with a graphical desktop session (X11)"  
        end

        private

        def construct_command
          [
            "node",
            File.join(Aquatone::AQUATONE_ROOT, "aquatone.js"),
            Shellwords.escape(url),
            Shellwords.escape(vhost),
            Shellwords.escape(html_destination),
            Shellwords.escape(headers_destination),
            Shellwords.escape(screenshot_destination)
          ].join(' ')
        end
      end
    end
  end
end
