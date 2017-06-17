module Aquatone
  class UrlMaker
    SSL_PORTS = [443,  832,  981,  1010, 1311, 2083, 2087,  2095,  2096,  4712,
                 7000, 8172, 8243, 8333, 8443, 8834, 9443,  12443, 18091, 18092].freeze

    def self.make(host, port)
      case port
      when 80
        "http://#{host}/"
      when 443
        "https://#{host}/"
      else
        if ssl_port?(port)
          "https://#{host}:#{port}/"
        else
          "http://#{host}:#{port}/"
        end
      end
    end

    private

    def self.ssl_port?(port)
      SSL_PORTS.include?(port.to_i)
    end
  end
end
