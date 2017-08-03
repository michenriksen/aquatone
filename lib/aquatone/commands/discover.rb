module Aquatone
  module Commands
    class Discover < Aquatone::Command
      def execute!
        if !options[:domain]
          output("Please specify a domain to assess\n")
          exit 1
        end

        @domain          = Aquatone::Domain.new(options[:domain])
        @assessment      = Aquatone::Assessment.new(options[:domain])
        @hosts           = [options[:domain]]
        @host_dictionary = {}

        banner("Discover")
        setup_resolver
        identify_wildcard_ips
        run_collectors
        resolve_hosts
        output_summary
        write_to_hosts_file
      rescue Aquatone::Domain::UnresolvableDomain => e
        output(red("Error: #{e.message}\n"))
      end

      private

      def setup_resolver
        if options[:nameservers]
          nameservers = options[:nameservers]
        else
          output("Identifying nameservers for #{@domain.name}... ")
          nameservers = @domain.nameservers
          output("Done\n")
          if nameservers.empty?
            output(yellow("#{@domain.name} has no nameservers. Using fallback nameservers.\n\n"))
            nameservers = []
          end
        end

        if !nameservers.empty?
          output("Using nameservers:\n\n")
          nameservers.each do |ns|
            output(" - #{ns}\n")
          end
          output("\n")
        end
        @resolver = Aquatone::Resolver.new(
          :nameservers          => nameservers,
          :fallback_nameservers => options[:fallback_nameservers]
        )
      end

      def identify_wildcard_ips
        output("Checking for wildcard DNS... ")
        @wildcard_ips   = []
        wildcard_domain = "#{random_string}.#{@domain.name}"
        if @resolver.resolve(wildcard_domain).nil?
          output("Done\n")
          return
        end
        output(yellow("Wildcard detected!\n"))
        output("Identifying wildcard IPs... ")
        20.times do
          wildcard_domain = "#{random_string}.#{@domain.name}"
          if wildcard_ip = @resolver.resolve(wildcard_domain)
            @wildcard_ips << wildcard_ip unless @wildcard_ips.include?(wildcard_ip)
          end
        end
        output("Done\n")
        output("Filtering out hosts resolving to wildcard IPs\n")
      end

      def run_collectors
        output("\n")
        Aquatone::Collector.descendants.each do |collector|
          next if skip_collector?(collector)
          output("Running collector: #{bold(collector.meta[:name])}... ")
          begin
            collector_instance = collector.new(@domain, options)
            hosts = collector_instance.execute!
            output("Done (#{hosts.count} #{hosts.count == 1 ? 'host' : 'hosts'})\n")
            @hosts += hosts
          rescue Aquatone::Collector::MissingKeyRequirement => e
            output(yellow("Skipped\n"))
            output(yellow(" -> #{e.message}\n"))
          rescue => e
            output(red("Error\n"))
            output(red(" -> #{e.message}\n"))
          end
        end
        @hosts = @hosts.sort.uniq
      end

      def resolve_hosts
        output("\nResolving #{bold(@hosts.count)} unique hosts...\n")
        task_count = 0
        @start_time = Time.now.to_i
        @hosts.each do |host|
          if asked_for_progress?
            output("Stats: #{seconds_to_time(Time.now.to_i - @start_time)} elapsed; " \
                   "#{task_count} out of #{@hosts.count} hosts checked (#{@host_dictionary.keys.count} discovered); " \
                   "#{(task_count.to_f / @hosts.count.to_f * 100.00).round(1)}% done\n")
          end
          if ip = @resolver.resolve(host)
            next if exclude_ip?(ip)
            @host_dictionary[host] = ip
            output("#{ip.ljust(15)} #{bold(host)}\n")
          end
          jitter_sleep
          task_count += 1
        end
        output("\n", true)
      end

      def output_summary
        subnets = find_subnets
        if !subnets.keys.count.zero?
          output("Found subnets:\n\n")
          subnets.each_pair do |subnet, count|
            next if count == 1
            subnet = "#{subnet}.0-255"
            output(" - #{subnet.ljust(17)} : #{count} hosts\n")
          end
        end
        output("\n")
      end

      def find_subnets
        subnets = {}
        @host_dictionary.values.each do |ip|
          subnet = ip.split(".")[0..2].join(".")
          if subnets.key?(subnet)
            subnets[subnet] += 1
          else
            subnets[subnet] = 1
          end
        end
        Hash[subnets.sort_by{|k, v| v}.reverse]
      end

      def write_to_hosts_file
        @hosts_file_contents = ""
        @host_dictionary.each_pair do |host, ip|
          @hosts_file_contents += "#{host},#{ip}\n"
        end
        @assessment.write_file("hosts.txt", @hosts_file_contents)
        @assessment.write_file("hosts.json", @host_dictionary.to_json)
        output("Wrote #{bold(@host_dictionary.keys.count)} hosts to:\n\n")
        output(" - #{bold('file://' + File.join(@assessment.path, 'hosts.txt'))}\n")
        output(" - #{bold('file://' + File.join(@assessment.path, 'hosts.json'))}\n")
      end

      def random_string
        %w(a b c d e f g h i j k l m n o p q r s t u v w x y z
           0 1 2 3 4 5 6 7 8 9).shuffle.take(10).join
      end

      def exclude_ip?(ip)
        wildcard_ip?(ip) || (options[:ignore_private] && private_ip?(ip)) || broadcast_ip?(ip)
      end

      def wildcard_ip?(ip)
        @wildcard_ips.include?(ip)
      end

      def private_ip?(ip)
        ip =~ /(\A127\.)|(\A10\.)|(\A172\.1[6-9]\.)|(\A172\.2[0-9]\.)|(\A172\.3[0-1]\.)|(\A192\.168\.)/
      end

      def broadcast_ip?(ip)
        ip == "255.255.255.255"
      end

      def skip_collector?(collector)
        if options[:only_collectors]
          if options[:only_collectors].include?(collector.sluggified_name)
            false
          else
            true
          end
        elsif options[:disable_collectors]
          if options[:disable_collectors].include?(collector.sluggified_name)
            true
          else
            false
          end
        else
          false
        end
      end
    end
  end
end
