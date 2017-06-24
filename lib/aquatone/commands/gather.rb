module Aquatone
  module Commands
    class Gather < Aquatone::Command
      def execute!
        if !options[:domain]
          output("Please specify a domain to assess\n")
          exit 1
        end

        @assessment = Aquatone::Assessment.new(options[:domain])

        banner("Gather")
        check_prerequisites
        prepare_tasks
        make_directories
        process_pages
        generate_report
      end

      private

      def check_prerequisites
        if !has_executable?("node")
          output(red("node executable not found!\n\n"))
          output("Please make sure Node.js is installed on your system.\n")
          exit 1
        end

        if !has_executable?("npm")
          output(red("npm executable not found!\n\n"))
          output("Please make sure NPM package manager is installed on your system.\n")
          exit 1
        end

        if !Dir.exists?(File.join(Aquatone::AQUATONE_ROOT, "node_modules"))
          output("Installing Nightmare.js package, please wait...")
          Dir.chdir(Aquatone::AQUATONE_ROOT) do
            if system("npm install nightmare >/dev/null 2>&1")
              output(" Done\n\n")
            else
              output(red(" Failed\n"))
              exit 1
            end
          end
        end
      end

      def prepare_tasks
        if !@assessment.has_file?("hosts.json")
          output(red("#{@assessment.path} does not contain hosts.json file\n\n"))
          output("Did you run aquatone-discover first?\n")
          exit 1
        end
        if !@assessment.has_file?("open_ports.txt")
          output(red("#{@assessment.path} does not contain open_ports.txt file\n\n"))
          output("Did you run aquatone-scan first?\n")
          exit 1
        end
        @tasks = []
        @hosts = JSON.parse(@assessment.read_file("hosts.json"))
        @open_ports = parse_open_ports_file
        @hosts.each_pair do |domain, host|
          next unless @open_ports.key?(host)
          @open_ports[host].each do |port|
            @tasks << [host, port, domain]
          end
        end
      end

      def process_pages
        pool        = thread_pool
        @task_count = 0
        @successful = 0
        @failed     = 0
        @start_time = Time.now.to_i
        @visits     = []
        output("Processing #{bold(@tasks.count)} pages...\n")
        @tasks.shuffle.each do |task|
          host, port, domain = task
          pool.schedule do
            begin
              output_progress if asked_for_progress?
              visit = visit_page(host, port, domain)
              if visit['success']
                output("#{green('Processed:')} #{Aquatone::UrlMaker.make(host, port)} (#{domain}) - #{visit['status']}\n")
                @successful += 1
              else
                output("   #{red('Failed:')} #{Aquatone::UrlMaker.make(host, port)} (#{domain}) - #{visit['error']} #{visit['details']}\n")
                @failed += 1
              end
              jitter_sleep
              @task_count += 1
            rescue Aquatone::Browser::Drivers::IncompatabilityError => e
              output("\n")
              output(red("Incompatability Error:") + " #{e.message}\n")
              exit 1
            end
          end
        end
        pool.shutdown
        output("\nFinished processing pages:\n\n")
        output(" - Successful : #{bold(green(@successful))}\n")
        output(" - Failed     : #{bold(red(@failed))}\n\n")
      end

      def generate_report
        output("Generating report...")
        report = Aquatone::Report.new(options[:domain], @visits)
        report.generate(File.join(@assessment.path, "report"))
        output("done\n")
        report_pages = Dir[File.join(@assessment.path, "report", "report_page_*.html")]
        output("Report pages generated:\n\n")
        sort_report_pages(report_pages).each do |report_page|
          output(" - file://#{report_page}\n")
        end
        output("\n")
      end

      def parse_open_ports_file
        contents = @assessment.read_file("open_ports.txt")
        result   = {}
        lines    = contents.split("\n").map(&:strip)
        lines.each do |line|
          values = line.split(",").map(&:strip)
          result[values[0]] = values[1..-1].map(&:to_i)
        end
        result
      end

      def make_directories
        @assessment.make_directory("headers")
        @assessment.make_directory("html")
        @assessment.make_directory("report")
        @assessment.make_directory("screenshots")
      end

      def make_file_basename(host, port, domain)
        "#{domain}__#{host}__#{port}".downcase.gsub(".", "_")
      end

      def output_progress
        output("Stats: #{seconds_to_time(Time.now.to_i - @start_time)} elapsed; " \
               "#{@task_count} out of #{@tasks.count} pages processed (#{@successful} successful, #{@failed} failed); " \
               "#{(@task_count.to_f / @tasks.count.to_f * 100.00).round(1)}% done\n")
      end

      def sort_report_pages(pages)
        pages.sort_by { |f| File.basename(f).split("_").last.split(".").first.to_i }
      end

      def visit_page(host, port, domain)
        file_basename          = make_file_basename(host, port, domain)
        url                    = Aquatone::UrlMaker.make(host, port)
        html_destination       = File.join(@assessment.path, "html", "#{file_basename}.html")
        headers_destination    = File.join(@assessment.path, "headers", "#{file_basename}.txt")
        screenshot_destination = File.join(@assessment.path, "screenshots", "#{file_basename}.png")
        visit = Aquatone::Browser.visit(url, domain, html_destination, headers_destination, screenshot_destination, :timeout => options[:timeout])
        if visit['success']
          @visits.push({
            :host          => host,
            :port          => port,
            :domain        => domain,
            :url           => url,
            :file_basename => file_basename,
            :headers       => visit['headers'],
            :status        => visit['status']
          })
        end
        visit
      end
    end
  end
end
