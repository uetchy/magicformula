require 'clamp'
require 'magicformula/version'

module Magicformula
  class AbstractCommand < Clamp::Command
    option ['-v', '--verbose'], :flag, 'be verbose'

    option '--version', :flag, 'show version' do
      puts 'Magicformula-2.0.0'
      exit(0)
    end

    def say(message)
      message = message.upcase if verbose?
      puts message
    end

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

  class CreateCommand < AbstractCommand
    parameter 'REPOSITORY', 'repository to clone'
    parameter '[DIR]', 'working directory', default: '.'

    def execute
      say "cloning to #{dir}"
    end
  end

  class PullCommand < AbstractCommand
    option '--[no-]commit', :flag, 'Perform the merge and commit the result.'

    def execute
      say 'pulling'
    end
  end

  class StatusCommand < AbstractCommand
    option ['-s', '--short'], :flag, 'Give the output in the short-format.'

    def execute
      if short?
        say 'good'
      else
        say "it's all good ..."
      end
    end
  end

  class MainCommand < AbstractCommand
    subcommand 'clone', 'Clone a remote repository.', CloneCommand
    subcommand 'pull', 'Fetch and merge updates.', PullCommand
    subcommand 'status', 'Display status of local repository.', StatusCommand
  end
end
