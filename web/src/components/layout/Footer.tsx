export default function Footer() {
  return (
    <footer className="h-12 bg-surface border-t border-border px-6 flex items-center justify-between text-body-xs text-text-secondary">
      <div className="flex items-center space-x-4">
        <span>APIWeaver v1.0.0</span>
        <span>•</span>
        <span>OpenAPI 3.1 Generator</span>
      </div>
      
      <div className="flex items-center space-x-4">
        <span>Ready</span>
        <div className="flex items-center space-x-1">
          <div className="h-2 w-2 bg-success rounded-full"></div>
          <span>Connected</span>
        </div>
      </div>
    </footer>
  )
}