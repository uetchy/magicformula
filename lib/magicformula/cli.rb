module Magicformula
  class CLI
    def load_plugins
      files = $LOAD_PATH.map do |_lp|
      end
      Dir.entries('plugins').each do |name|
        path = "plugins/#{name}"
        next if File.ftype(path) != 'directory' || %w(. ..).include?(name)

        require_relative path + '/plugin.rb'
        name[0] = name[0].upcase
        plugins << Module.const_get("#{name}Plugin").new
      end

      loop do
        msg = gets.chomp!
        plugins.each { |plugin| plugin.on_message(msg) }
        puts msg
      end
    end
  end
end
