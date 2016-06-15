# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'magicformula/version'

Gem::Specification.new do |spec|
  spec.name          = 'magicformula'
  spec.version       = Magicformula::VERSION
  spec.authors       = ['Yasuaki Uechi']
  spec.email         = ['uetchy@randompaper.co']

  spec.summary       = 'Generate and upload Homebrew Formula like magic.'
  spec.description   = 'Generate and upload Homebrew Formula like magic.'
  spec.homepage      = "https://github.com/uetchy/magicformula"
  spec.license       = 'MIT'
  spec.files         = `git ls-files -z`.split("\x0").reject { |f| f.match(%r{^(test|spec|features)/}) }
  spec.bindir        = 'exe'
  spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = ['lib']

  spec.add_dependency 'thor'
  spec.add_dependency 'baby_erubis'
  spec.add_development_dependency 'bundler', '~> 1.12'
  spec.add_development_dependency 'rake', '~> 11.2'
  spec.add_development_dependency 'rspec', '~> 3.4'
end
