module MagicformulaPlugin
  class Handler
    def call(_args)
      template = File.open('formula.erb').read
    end
  end
end
