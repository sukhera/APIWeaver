import { Link, useLocation } from 'react-router-dom'
import { cn } from '@/lib/utils'
import { 
  Home, 
  FileText, 
  CheckCircle, 
  GitMerge, 
  History, 
  BookOpen,
  Zap
} from 'lucide-react'

const navigationItems = [
  {
    name: 'Home',
    href: '/',
    icon: Home,
  },
  {
    name: 'Generate',
    href: '/generate',
    icon: Zap,
  },
  {
    name: 'Validate',
    href: '/validate',
    icon: CheckCircle,
  },
  {
    name: 'Amend',
    href: '/amend',
    icon: GitMerge,
  },
  {
    name: 'History',
    href: '/history',
    icon: History,
  },
]

const quickActions = [
  {
    name: 'Example APIs',
    icon: BookOpen,
    action: 'examples',
  },
  {
    name: 'Recent Files',
    icon: FileText,
    action: 'recent',
  },
]

export default function Sidebar() {
  const location = useLocation()
  
  return (
    <aside className="w-64 bg-surface border-r border-border overflow-y-auto">
      <div className="p-6">
        {/* Navigation */}
        <nav className="space-y-1">
          {navigationItems.map((item) => {
            const isActive = location.pathname === item.href
            const Icon = item.icon
            
            return (
              <Link
                key={item.name}
                to={item.href}
                className={cn(
                  'flex items-center px-3 py-2 text-sm rounded-md transition-colors',
                  isActive
                    ? 'bg-primary/10 text-primary border-r-2 border-primary'
                    : 'text-text-secondary hover:bg-surface-elevated hover:text-text-primary'
                )}
              >
                <Icon className="mr-3 h-5 w-5" />
                {item.name}
              </Link>
            )
          })}
        </nav>
        
        {/* Quick Actions */}
        <div className="mt-8">
          <h3 className="text-label-sm font-medium text-text-tertiary uppercase tracking-wider mb-3">
            Quick Actions
          </h3>
          <div className="space-y-1">
            {quickActions.map((item) => {
              const Icon = item.icon
              
              return (
                <button
                  key={item.name}
                  className="flex items-center w-full px-3 py-2 text-sm text-text-secondary hover:bg-surface-elevated hover:text-text-primary rounded-md transition-colors"
                  onClick={() => {
                    // TODO: Handle quick actions
                    console.log(`Quick action: ${item.action}`)
                  }}
                >
                  <Icon className="mr-3 h-5 w-5" />
                  {item.name}
                </button>
              )
            })}
          </div>
        </div>
        
        {/* Status */}
        <div className="mt-8 p-3 bg-muted rounded-md">
          <div className="flex items-center space-x-2">
            <div className="h-2 w-2 bg-success rounded-full"></div>
            <span className="text-body-xs text-text-secondary">
              Service Online
            </span>
          </div>
        </div>
      </div>
    </aside>
  )
}