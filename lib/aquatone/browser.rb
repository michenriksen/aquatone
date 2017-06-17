module Aquatone
  class Browser
    def self.visit(url, vhost, html_destination, headers_destination, screenshot_destination, options)
      driver = make_driver(url, vhost, html_destination, headers_destination, screenshot_destination, options)
      visit  = driver.visit
      if !visit["success"]
        visit = driver.visit
      end
      visit
    end

    protected

    def self.make_driver(url, vhost, html_destination, headers_destination, screenshot_destination, options)
      Aquatone::Browser::Drivers::Nightmare.new(url, vhost, html_destination, headers_destination, screenshot_destination, options)
    end
  end
end
