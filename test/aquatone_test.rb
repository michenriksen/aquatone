require 'test_helper'

describe Aquatone do
  it "has a version" do
    Aquatone::VERSION.wont_be_nil
  end

  describe Aquatone::DEFAULT_AQUATONE_PATH do
    it "is a folder in the current user's home folder" do
      Aquatone::DEFAULT_AQUATONE_PATH.must_equal File.join(Dir.home, "aquatone")
    end
  end

  describe ".aquatone_path" do
    describe "when ENV[AQUATONEPATH] is not set" do
      it "returns default aquatone path" do
        ENV["AQUATONEPATH"] = nil
        Aquatone.aquatone_path.must_equal File.join(Dir.home, "aquatone")
      end
    end

    describe "when ENV[AQUATONEPATH] is set" do
      it "returns the custom path" do
        path = File.join(Dir.home, "someotherplace")
        ENV["AQUATONEPATH"] = path
        Aquatone.aquatone_path.must_equal path
        ENV["AQUATONEPATH"] = nil
      end
    end
  end
end
