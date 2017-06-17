module Aquatone
  class KeyStore
    class Error < StandardError; end
    class KeyStoreFileNotReadable < Error; end
    class KeyStoreFileNotWritable < Error; end
    class KeyStoreFileCorrupt < Error; end

    KEY_STORE_FILE_LOCATION = File.join(Aquatone.aquatone_path, ".keys.yml").freeze

    def self.get(name)
      keys[name]
    end

    def self.set(name, value)
      k = keys
      k[name] = value
      write_key_store_file(keys)
    end

    def self.keys
      @keys ||= read_key_store_file
    end

    def self.key?(name)
      keys.key?(name)
    end

    def self.reset!
      @keys = nil
    end

    private

    def self.read_key_store_file
      return {} unless key_store_exists?
      fail KeyStoreFileNotReadable, "Key store file is not readable" unless key_store_readable?
      deserialize(File.read(KEY_STORE_FILE_LOCATION))
    end

    def self.write_key_store_file(keys)
      if key_store_exists?
        fail KeyStoreFileNotWritable, "Key store file is not writable" unless key_store_writable?
      end
      File.open(KEY_STORE_FILE_LOCATION, "w") do |file|
        file.write(serialize(keys))
      end
      @keys = nil
    end

    def self.key_store_readable?
      File.readable?(KEY_STORE_FILE_LOCATION)
    end

    def self.key_store_writable?
      File.writable?(KEY_STORE_FILE_LOCATION)
    end

    def self.key_store_exists?
      File.exists?(KEY_STORE_FILE_LOCATION)
    end

    def self.serialize(keys)
      YAML.dump(keys)
    end

    def self.deserialize(keys)
      YAML.parse(keys).to_ruby
    rescue Psych::SyntaxError
      fail KeyStoreFileCorrupt, "Key store file contains invalid YAML"
    end
  end
end
