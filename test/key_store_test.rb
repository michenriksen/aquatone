require 'test_helper'

describe Aquatone::KeyStore do
  before { Aquatone::KeyStore.reset! }

  describe Aquatone::KeyStore::KEY_STORE_FILE_LOCATION do
    it "is a YAML file in the AQUATONE folder" do
      Aquatone::KeyStore::KEY_STORE_FILE_LOCATION.must_equal File.join(Aquatone.aquatone_path, ".keys.yml")
    end
  end

  describe ".get" do
    before do
      @key_file_contents = YAML.dump({
        "key" => "deadbeefdeadbeefdeadbeef"
      })
    end

    describe "when key store file does not exist" do
      it "returns nil" do
        File.stub :exists?, false do
          Aquatone::KeyStore.get("key").must_be_nil
        end
      end
    end

    describe "when given a key name that is not set" do
      it "returns nil" do
        File.stub :read, @key_file_contents do
          Aquatone::KeyStore.get("somethingelse").must_be_nil
        end
      end
    end

    describe "when given a key name that is set" do
      it "returns the key value" do
        File.stub :read, @key_file_contents do
          Aquatone::KeyStore.get("key").must_equal "deadbeefdeadbeefdeadbeef"
        end
      end
    end
  end

  describe ".set" do
    it "writes keys to key store file" do
      file_mock = Minitest::Mock.new
      file_mock.expect(:write, YAML.dump({"key" => "deadbeefdeadbeefdeadbeef"}).length) do |arg|
        arg == YAML.dump({"key" => "deadbeefdeadbeefdeadbeef"})
      end
      # TODO: Figure out why this stubbing doesn't work...
      File.stub :open, file_mock do
        Aquatone::KeyStore.set("key", "deadbeefdeadbeefdeadbeef")
      end
      file_mock.verify
    end
  end
end
