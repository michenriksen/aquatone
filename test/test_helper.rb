$LOAD_PATH.unshift File.expand_path('../../lib', __FILE__)
require 'aquatone'

require 'minitest/autorun'

def time_taken
  now = Time.now.to_f
  yield
  Time.now.to_f - now
end
