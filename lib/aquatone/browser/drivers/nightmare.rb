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
          rout, wout = IO.pipe
          process           = ChildProcess.build(*construct_command)
          process.cwd       = Aquatone::AQUATONE_ROOT
          process.io.stdout = wout
          process.start
          process.poll_for_exit(options[:timeout])
          wout.close
          command_output = rout.readlines.join("\n").strip
          JSON.parse(command_output)
        rescue JSON::ParserError
          fail IncompatabilityError, "Nightmarejs must be run on a system with a graphical desktop session (X11)"
        rescue => e
          process.stop if process
          return {
            "success" => false,
            "error"   => e.is_a?(ChildProcess::TimeoutError) ? "Timeout" : "#{e.class}: #{e.message}",
            "code"    => 0,
            "details" => ""
          }
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
          ]
        end
      end
    end
  end
end
