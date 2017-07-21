module Aquatone
  class Detector
    class Error < StandardError; end
    class InvalidMetadataError < Error; end
    class MetadataNotSetError < Error; end

    attr_reader :host, :resource

    def self.meta
      @meta || fail(MetadataNotSetError, "Metadata has not been set")
    end

    def self.meta=(meta)
      validate_metadata(meta)
      @meta = meta
    end

    def self.descendants
      detectors = ObjectSpace.each_object(Class).select { |klass| klass < self }
      detectors.sort { |x, y| x.meta[:service] <=> y.meta[:service] }
    end

    def self.sluggified_name
      return meta[:slug].downcase if meta[:slug]
      meta[:service].strip.downcase.gsub(/[^a-z0-9]+/, '-').gsub("--", "-")
    end

    def initialize(host, resource)
      @host     = host
      @resource = resource
    end

    def run
      fail NotImplementedError
    end

    def positive?
      run
    rescue
      false
    end

    protected

    def cname_resource?
      resource.is_a?(Resolv::DNS::Resource::IN::CNAME)
    end

    def apex_resource?
      resource.is_a?(Resolv::DNS::Resource::IN::A)
    end

    def resource_value
      cname_resource? ? resource.name.to_s : resource.address.to_s
    end

    def get_request(uri, options={})
      options = {
        :timeout => 10
      }.merge(options)
      Aquatone::HttpClient.get(uri, options)
    end

    def post_request(uri, body=nil, options={})
      options = {
        :body    => body,
        :timeout => 10
      }.merge(options)
      Aquatone::HttpClient.post(uri, options)
    end

    def url_escape(string)
      CGI.escape(string)
    end

    def random_sleep(seconds)
      random_sleep = ((1 - (rand(30) * 0.01)) * seconds.to_i)
      sleep(random_sleep)
    end

    def failure(message)
      fail Error, message
    end

    def self.validate_metadata(meta)
      fail InvalidMetadataError, "Metadata is not a hash" unless meta.is_a?(Hash)
      fail InvalidMetadataError, "Metadata is empty" if meta.empty?
      fail InvalidMetadataError, "Metadata is missing key: service" unless meta.key?(:service)
      fail InvalidMetadataError, "Metadata is missing key: service_website" unless meta.key?(:service_website)
      fail InvalidMetadataError, "Metadata is missing key: author" unless meta.key?(:author)
      fail InvalidMetadataError, "Metadata is missing key: description" unless meta.key?(:description)
    end
  end
end
