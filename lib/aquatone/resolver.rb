module Aquatone
  class Resolver
    attr_reader :options

    TIMEOUT = 3.freeze

    def initialize(options)
      @options = {
        :timeout => TIMEOUT
      }.merge(options)
    end

    def resolve(host)
      retried = false
      host    = "#{host}." unless host.end_with?(".")
      return (resolve_with_nameserver(host) || resolve_with_fallback_nameserver(host))
    rescue
      if !retried
        retried = true
        retry
      end
      nil
    end

    private

    def resolve_with_nameserver(host)
      _resolve(host, options[:nameservers].sample)
    end

    def resolve_with_fallback_nameserver(host)
      _resolve(host, options[:fallback_nameservers].sample)
    end

    def _resolve(host, nameserver_ip)
      nameserver = Resolv::DNS.new(:nameserver => nameserver_ip).tap do |ns|
        ns.timeouts = options[:timeout]
      end
      ip = nameserver.getaddress(host).to_s
      nameserver.close
      ip
    rescue Resolv::ResolvError
      nameserver.close
      nil
    end
  end
end
