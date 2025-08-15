import { Link } from 'react-router-dom'
import { ThemeToggle } from '@/components/theme-toggle'
import { Button } from '@/components/ui/button'
import { FileText, Github } from 'lucide-react'

export default function Header() {
  return (
    <header className="h-16 bg-surface border-b border-border px-6 flex items-center justify-between">
      {/* Left side - Logo and navigation */}
      <div className="flex items-center space-x-6">
        {/* Logo */}
        <Link to="/" className="flex items-center space-x-2">
          <FileText className="h-6 w-6 text-primary" />
          <span className="text-heading-sm font-bold text-text-primary">
            APIWeaver
          </span>
        </Link>
        
        {/* Navigation */}
        <nav className="hidden md:flex items-center space-x-4">
          <Link 
            to="/generate" 
            className="text-text-secondary hover:text-text-primary transition-colors text-body-sm"
          >
            Generate
          </Link>
          <Link 
            to="/validate" 
            className="text-text-secondary hover:text-text-primary transition-colors text-body-sm"
          >
            Validate
          </Link>
          <Link 
            to="/amend" 
            className="text-text-secondary hover:text-text-primary transition-colors text-body-sm"
          >
            Amend
          </Link>
          <Link 
            to="/history" 
            className="text-text-secondary hover:text-text-primary transition-colors text-body-sm"
          >
            History
          </Link>
        </nav>
      </div>
      
      {/* Right side - Actions */}
      <div className="flex items-center space-x-4">
        <Button variant="outline" size="sm" asChild>
          <a 
            href="https://github.com/your-org/apiweaver" 
            target="_blank" 
            rel="noopener noreferrer"
            className="flex items-center space-x-2"
          >
            <Github className="h-4 w-4" />
            <span className="hidden sm:inline">GitHub</span>
          </a>
        </Button>
        <ThemeToggle />
      </div>
    </header>
  )
}