require "resolv"
require "ipaddr"
require "socket"
require "timeout"
require "shellwords"
require "optparse"

require "httparty"
require "childprocess"

require "aquatone/version"
require "aquatone/port_lists"
require "aquatone/url_maker"
require "aquatone/validation"
require "aquatone/thread_pool"
require "aquatone/http_client"
require "aquatone/browser"
require "aquatone/browser/drivers/nightmare"
require "aquatone/domain"
require "aquatone/resolver"
require "aquatone/assessment"
require "aquatone/report"
require "aquatone/command"
require "aquatone/collector"
require "aquatone/detector"

module Aquatone
  AQUATONE_ROOT         = File.expand_path(File.join(File.dirname(__FILE__), "..")).freeze
  DEFAULT_AQUATONE_PATH = File.join(Dir.home, "aquatone").freeze

  def self.aquatone_path
    ENV['AQUATONEPATH'] || DEFAULT_AQUATONE_PATH
  end
end

require "aquatone/key_store"

Dir[File.join(File.dirname(__FILE__), "aquatone", "collectors", "*.rb")].each do |collector|
  require collector
end

Dir[File.join(File.dirname(__FILE__), "aquatone", "detectors", "*.rb")].each do |detector|
  require detector
end

require "aquatone/commands/discover"
require "aquatone/commands/scan"
require "aquatone/commands/gather"
require "aquatone/commands/takeover"
