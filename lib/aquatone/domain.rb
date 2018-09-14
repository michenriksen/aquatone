module Aquatone
  class Domain

    class UnresolvableDomain < StandardError; end

    attr_reader :name

    def initialize(name, options = {})
      @name    = name
      @options = options
    end

    def nameservers
      result = []
      parts  = name.split(".")
      parts.size.times do |n|
        lookup      = parts[n..-1].join('.') + "."
        nameservers = nameserver.getresources(lookup, Resolv::DNS::Resource::IN::NS)
        if !nameservers.count.zero?
          result = nameservers.map do |ns|
            begin
              nameserver.getaddress(ns.name.to_s).to_s
            rescue
              nil
            end
          end.compact
          break
        end
      end
      result
    end

    private

    def nameserver
      @nameserver ||= Resolv::DNS.new
    end
  end
end
