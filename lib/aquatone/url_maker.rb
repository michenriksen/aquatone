module Aquatone
  class UrlMaker
    SSL_PORTS = [443,  454,  455,  832,  981,  1010, 1311, 2053, 2083, 2087, 
                 2096, 4016, 4018, 4020, 4022, 4712, 7000, 8081, 8172, 8243,
                 8333, 8443, 8834, 9443, 18091, 18092].freeze

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
