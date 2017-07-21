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

    def resource(host)
      resource = resource_with_nameserver(host,
        Resolv::DNS::Resource::IN::CNAME) ||
      resource_with_fallback_nameserver(host,
        Resolv::DNS::Resource::IN::CNAME)
      if !resource
        resource = resource_with_nameserver(host,
          Resolv::DNS::Resource::IN::A) ||
        resource_with_fallback_nameserver(host,
          Resolv::DNS::Resource::IN::A)
      end
      resource
    end

    private

    def resolve_with_nameserver(host)
      _resolve(host, options[:nameservers].sample)
    end

    def resolve_with_fallback_nameserver(host)
      _resolve(host, options[:fallback_nameservers].sample)
    end

    def resource_with_nameserver(host, typeclass)
      _resource(host, typeclass, options[:nameservers].sample)
    end

    def resource_with_fallback_nameserver(host, typeclass)
      _resource(host, typeclass, options[:fallback_nameservers].sample)
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

    def _resource(host, typeclass, nameserver_ip)
      nameserver = Resolv::DNS.new(:nameserver => nameserver_ip).tap do |ns|
        ns.timeouts = options[:timeout]
      end
      resource = nameserver.getresource(host, typeclass)
      nameserver.close
      resource
    rescue Resolv::ResolvError
      nameserver.close
      nil
    end
  end
end
