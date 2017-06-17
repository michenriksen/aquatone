module Aquatone
  class Assessment
    attr_reader :domain

    def initialize(domain)
      @domain = domain
      initialize_assessment_directory
    end

    def has_file?(name)
      File.exist?(File.join(path, name))
    end

    def read_file(name)
      File.read(File.join(path, name))
    end

    def write_file(name, data, mode = "w")
      File.open(File.join(path, name), mode) do |file|
        file.write(data)
      end
    end

    def make_directory(name)
      dir = File.join(path, name)
      Dir.mkdir(dir) unless Dir.exist?(dir)
    end

    def path
      File.join(Aquatone.aquatone_path, domain)
    end

    private

    def initialize_assessment_directory
      Dir.mkdir(Aquatone.aquatone_path) unless Dir.exist?(Aquatone.aquatone_path)
      Dir.mkdir(path) unless Dir.exist?(path)
    end
  end
end
