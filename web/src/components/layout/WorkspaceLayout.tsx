import { ReactNode, useState } from 'react'
import { ResizablePanelGroup, ResizablePanel, ResizableHandle } from '@/components/ui/resizable'
import { cn } from '@/lib/utils'

interface WorkspaceLayoutProps {
  leftPanel?: ReactNode
  centerPanel: ReactNode
  rightPanel?: ReactNode
  leftWidth?: number
  centerWidth?: number
  rightWidth?: number
  showLeftPanel?: boolean
  showRightPanel?: boolean
  className?: string
}

export default function WorkspaceLayout({
  leftPanel,
  centerPanel,
  rightPanel,
  leftWidth = 25,
  centerWidth = 50,
  rightWidth = 25,
  showLeftPanel = true,
  showRightPanel = true,
  className
}: WorkspaceLayoutProps) {
  const [panels] = useState({
    left: showLeftPanel,
    right: showRightPanel
  })

  // Calculate panel configuration based on visible panels
  const getPanelConfig = () => {
    if (panels.left && panels.right) {
      return {
        leftSize: leftWidth,
        centerSize: centerWidth,
        rightSize: rightWidth
      }
    } else if (panels.left && !panels.right) {
      return {
        leftSize: leftWidth,
        centerSize: 100 - leftWidth,
        rightSize: 0
      }
    } else if (!panels.left && panels.right) {
      return {
        leftSize: 0,
        centerSize: 100 - rightWidth,
        rightSize: rightWidth
      }
    } else {
      return {
        leftSize: 0,
        centerSize: 100,
        rightSize: 0
      }
    }
  }

  const config = getPanelConfig()

  return (
    <div className={cn('h-full', className)}>
      <ResizablePanelGroup direction="horizontal" className="h-full">
        {/* Left Panel */}
        {panels.left && leftPanel && (
          <>
            <ResizablePanel
              defaultSize={config.leftSize}
              minSize={15}
              maxSize={50}
              className="h-full"
            >
              <div className="h-full overflow-hidden">
                {leftPanel}
              </div>
            </ResizablePanel>
            <ResizableHandle withHandle />
          </>
        )}

        {/* Center Panel */}
        <ResizablePanel
          defaultSize={config.centerSize}
          minSize={30}
          className="h-full"
        >
          <div className="h-full overflow-hidden">
            {centerPanel}
          </div>
        </ResizablePanel>

        {/* Right Panel */}
        {panels.right && rightPanel && (
          <>
            <ResizableHandle withHandle />
            <ResizablePanel
              defaultSize={config.rightSize}
              minSize={15}
              maxSize={50}
              className="h-full"
            >
              <div className="h-full overflow-hidden">
                {rightPanel}
              </div>
            </ResizablePanel>
          </>
        )}
      </ResizablePanelGroup>
    </div>
  )
}