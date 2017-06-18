module Aquatone
  module Commands
    class Scan < Aquatone::Command
      def execute!
        if !options[:domain]
          output("Please specify a domain to assess\n")
          exit 1
        end

        @assessment      = Aquatone::Assessment.new(options[:domain])
        @tasks           = []
        @host_dictionary = {}
        @results         = {}
        @urls            = []

        banner("Scan")
        prepare_host_dictionary
        scan_ports
        write_open_ports_file
        write_urls_file
      end

      def prepare_host_dictionary
        if !@assessment.has_file?("hosts.json")
          output(red("#{@assessment.path} does not contain hosts.json file\n\n"))
          output("Did you run aquatone-discover first?\n")
          exit 1
        end
        hosts = JSON.parse(@assessment.read_file("hosts.json"))
        output("Loaded #{bold(hosts.count)} hosts from #{bold(File.join(@assessment.path, 'hosts.json'))}\n\n")
        hosts.each_pair do |domain, ip|
          if !@host_dictionary.key?(ip)
            @host_dictionary[ip] = [domain]
            options[:ports].each do |port|
              @tasks << [ip, port]
            end
          else
            @host_dictionary[ip] << domain
          end
        end
      end

      def scan_ports
        pool        = thread_pool
        @task_count = 1
        @ports_open = 0
        @start_time = Time.now.to_i
        output("Probing #{bold(@tasks.count)} ports...\n")
        @tasks.shuffle.each do |task|
          host, port = task
          pool.schedule do
            output_progress if asked_for_progress?
            if port_open?(host, port)
              output_open_port(host, port)
              @ports_open += 1
              @results[host] ||= []
              @results[host] << port
              @host_dictionary[host].each do |hostname|
                @urls << Aquatone::UrlMaker.make(hostname, port)
              end
            end
            jitter_sleep
            @task_count += 1
          end
        end
        pool.shutdown
      end

      def write_open_ports_file
        contents = []
        @results.each_pair do |host, ports|
          contents << "#{host},#{ports.join(',')}"
        end
        @assessment.write_file("open_ports.txt", contents.sort.join("\n"))
        output("\nWrote open ports to #{bold('file://' + File.join(@assessment.path, 'open_ports.txt'))}\n")
      end

      def write_urls_file
        @assessment.write_file("urls.txt", @urls.uniq.sort.join("\n"))
        output("Wrote URLs to #{bold('file://' + File.join(@assessment.path, 'urls.txt'))}\n")
      end

      def port_open?(host, port)
        Timeout::timeout(options[:timeout]) do
          TCPSocket.new(host, port).close
          true
        end
      rescue Timeout::Error, Errno::ECONNREFUSED, Errno::EHOSTUNREACH, Errno::ENETUNREACH, SocketError
        false
      end

      def output_progress
        output("Stats: #{seconds_to_time(Time.now.to_i - @start_time)} elapsed; " \
               "#{@task_count} out of #{@tasks.count} ports checked (#{@ports_open} open); " \
               "#{(@task_count.to_f / @tasks.count.to_f * 100.00).round(1)}% done\n")
      end

      def output_open_port(host, port)
        if (@host_dictionary[host].count > 3)
          domains = @host_dictionary[host].shuffle.take(3).join(", ") + " and #{@host_dictionary[host].count - 3} more"
        else
          domains = @host_dictionary[host].take(3).join(", ")
        end
        output("#{green((port.to_s + '/tcp').ljust(9))} #{host.ljust(15)} #{domains}\n")
      end
    end
  end
end
