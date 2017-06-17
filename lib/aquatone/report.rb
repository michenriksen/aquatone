module Aquatone
  class Report
    def initialize(domain, visits, options = {})
      @domain  = domain
      @visits  = visits
      @options = {
        :per_page => 100,
        :template => "default"
      }.merge(options)
    end

    def generate(destination)
      sorted_visits = @visits.sort { |x,y| x[:domain] <=> y[:domain] }
      slices        = sorted_visits.each_slice(@options[:per_page].to_i)
      report        = load_template
      slices.each_with_index do |h, i|
        b            = binding
        @visit_slice = h
        @page_number = i

        if i + 1 == slices.count
          @link_to_next_page = false
        else
          @link_to_next_page = true
          @next_page_path    = File.basename(report_file_name(destination, i + 1))
        end

        if i.zero?
          @link_to_previous_page = false
        else
          @link_to_previous_page = true
          @previous_page_path    = File.basename(report_file_name(destination, i - 1))
        end

        File.open(report_file_name(destination, i), "w") do |f|
          f.write(report.result(b))
        end
      end
    end

    private

    def load_template
      ERB.new(File.read(File.join(Aquatone::AQUATONE_ROOT, "templates", "#{@options[:template]}.html.erb")))
    end

    def h(unsafe)
      CGI.escapeHTML(unsafe.to_s)
    end

    def report_file_name(destination, page_number)
     File.join(destination, "report_page_#{page_number}.html")
    end

    def url(domain, port)
      Aquatone::UrlMaker.make(domain, port)
    end

    def header_row_class?(header, value)
      case header.downcase
      when 'server', 'x-powered-by'
        :danger
      when 'access-control-allow-origin'
        :danger if value == '*'
      when 'content-security-policy'
        :success
      when 'x-permitted-cross-domain-policies'
        :success if value.downcase == 'master-only'
      when 'x-content-type-options'
        :success if value.downcase == 'nosniff'
      when 'strict-transport-security'
        :success
      when 'x-frame-options'
        :success
      when 'public-key-pins'
        :success
      when 'x-xss-protection'
        if value.start_with?('1')
          :success
        else
          :danger
        end
      else
        false
      end
    end
  end
end
