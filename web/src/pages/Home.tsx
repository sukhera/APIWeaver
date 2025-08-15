import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { FileText, Zap, CheckCircle, GitMerge } from 'lucide-react'

export default function Home() {
  const features = [
    {
      icon: Zap,
      title: 'Generate',
      description: 'Convert Markdown requirements to OpenAPI 3.1 specifications',
      href: '/generate',
      color: 'text-blue-500',
      bgColor: 'bg-blue-50 dark:bg-blue-900/20',
    },
    {
      icon: CheckCircle,
      title: 'Validate',
      description: 'Validate your OpenAPI specs with detailed error reporting',
      href: '/validate',
      color: 'text-green-500',
      bgColor: 'bg-green-50 dark:bg-green-900/20',
    },
    {
      icon: GitMerge,
      title: 'Amend',
      description: 'Update existing specs with changes and see diffs',
      href: '/amend',
      color: 'text-purple-500',
      bgColor: 'bg-purple-50 dark:bg-purple-900/20',
    },
  ]

  return (
    <div className="container mx-auto px-6 py-12">
      {/* Hero Section */}
      <div className="text-center mb-16">
        <div className="flex justify-center mb-6">
          <div className="p-4 bg-primary/10 rounded-full">
            <FileText className="h-12 w-12 text-primary" />
          </div>
        </div>
        
        <h1 className="text-display-lg font-bold text-text-primary mb-4">
          OpenAPI 3.1 Generator
        </h1>
        
        <p className="text-body-lg text-text-secondary max-w-2xl mx-auto mb-8">
          Transform your Markdown API requirements into valid OpenAPI 3.1 specifications. 
          Generate, validate, and amend your API documentation with ease.
        </p>
        
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Button size="lg" asChild>
            <Link to="/generate" className="flex items-center space-x-2">
              <Zap className="h-5 w-5" />
              <span>Get Started</span>
            </Link>
          </Button>
          
          <Button size="lg" variant="outline" asChild>
            <a 
              href="https://github.com/your-org/apiweaver" 
              target="_blank" 
              rel="noopener noreferrer"
              className="flex items-center space-x-2"
            >
              <FileText className="h-5 w-5" />
              <span>View Documentation</span>
            </a>
          </Button>
        </div>
      </div>
      
      {/* Features Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
        {features.map((feature) => {
          const Icon = feature.icon
          
          return (
            <Link
              key={feature.title}
              to={feature.href}
              className="group p-6 border border-border rounded-lg hover:border-primary/50 transition-colors"
            >
              <div className={`inline-flex p-3 rounded-lg ${feature.bgColor} mb-4`}>
                <Icon className={`h-6 w-6 ${feature.color}`} />
              </div>
              
              <h3 className="text-heading-sm font-semibold text-text-primary mb-2 group-hover:text-primary transition-colors">
                {feature.title}
              </h3>
              
              <p className="text-body-sm text-text-secondary">
                {feature.description}
              </p>
            </Link>
          )
        })}
      </div>
      
      {/* Quick Start */}
      <div className="bg-muted rounded-lg p-8">
        <h2 className="text-heading-md font-semibold text-text-primary mb-4">
          Quick Start
        </h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div>
            <h3 className="text-heading-xs font-medium text-text-primary mb-3">
              Via Web Interface
            </h3>
            <ol className="space-y-2 text-body-sm text-text-secondary">
              <li>1. Upload or paste your Markdown API requirements</li>
              <li>2. Click "Generate" to create your OpenAPI spec</li>
              <li>3. Review, validate, and download your specification</li>
            </ol>
          </div>
          
          <div>
            <h3 className="text-heading-xs font-medium text-text-primary mb-3">
              Via CLI
            </h3>
            <div className="bg-surface rounded border p-3 font-mono text-body-sm">
              <code className="text-text-primary">
                openapi-gen generate requirements.md -o api.yaml
              </code>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}