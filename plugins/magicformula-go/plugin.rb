module MagicformulaPlugin
  class GolangFormulaGenerator < Magicformula::AbstractPlugin
    def render(_args)
      template = File.open('formula.erb').read
    end

    # Collect args for golang project
    def collect(project_path)

    end
  end
end
