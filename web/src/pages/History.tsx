import { useState, useEffect, useCallback } from 'react'
import { toast } from 'sonner'
import { History as HistoryIcon, FileText, Download, Eye, Trash2, Search, Calendar, CheckCircle, AlertCircle } from 'lucide-react'

// Components
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import MonacoEditor from '@/components/editor/MonacoEditor'

// Hooks and Store
import { useAppStore } from '@/store/useAppStore'

// Types
interface ConversionHistoryItem {
  id: string
  timestamp: Date
  inputType: 'markdown' | 'openapi'
  inputContent: string
  outputContent?: string
  outputFormat: 'yaml' | 'json'
  success: boolean
  errors?: string[]
  warnings?: string[]
  processingTime?: number
  inputSize: number
  outputSize?: number
  operation: 'generate' | 'validate' | 'amend'
  title?: string
}

interface HistoryFilters {
  search: string
  operation: 'all' | 'generate' | 'validate' | 'amend'
  status: 'all' | 'success' | 'error'
  dateRange: 'all' | 'today' | 'week' | 'month'
}

export default function History() {
  const { updateMarkdown, updateSpec, setFormat } = useAppStore()
  
  const [history, setHistory] = useState<ConversionHistoryItem[]>([])
  const [filteredHistory, setFilteredHistory] = useState<ConversionHistoryItem[]>([])
  const [selectedItem, setSelectedItem] = useState<ConversionHistoryItem | null>(null)
  const [isDetailDialogOpen, setIsDetailDialogOpen] = useState(false)
  const [filters, setFilters] = useState<HistoryFilters>({
    search: '',
    operation: 'all',
    status: 'all',
    dateRange: 'all'
  })

  const createMockHistory = useCallback((): ConversionHistoryItem[] => {
    const now = new Date()
    return [
      {
        id: '1',
        timestamp: new Date(now.getTime() - 2 * 60 * 60 * 1000),
        inputType: 'markdown',
        inputContent: '# Task API\n\n## GET /tasks\nRetrieve all tasks.',
        outputContent: 'openapi: 3.1.0\ninfo:\n  title: Task API\n  version: 1.0.0\npaths:\n  /tasks:\n    get:\n      summary: Retrieve all tasks',
        outputFormat: 'yaml',
        success: true,
        processingTime: 245,
        inputSize: 45,
        outputSize: 156,
        operation: 'generate',
        title: 'Task API Generation'
      }
    ]
  }, [])

  const saveHistory = useCallback((historyData: ConversionHistoryItem[]) => {
    try {
      localStorage.setItem('apiweaver-history', JSON.stringify(historyData))
    } catch (error) {
      console.error('Failed to save history:', error)
      toast.error('Failed to save conversion history')
    }
  }, [])

  const loadHistory = useCallback(() => {
    try {
      const storedHistory = localStorage.getItem('apiweaver-history')
      if (storedHistory) {
        const parsed = JSON.parse(storedHistory)
        const historyWithDates = parsed.map((item: ConversionHistoryItem) => ({
          ...item,
          timestamp: new Date(item.timestamp)
        }))
        setHistory(historyWithDates)
      } else {
        const mockHistory = createMockHistory()
        setHistory(mockHistory)
        saveHistory(mockHistory)
      }
    } catch (error) {
      console.error('Failed to load history:', error)
      toast.error('Failed to load conversion history')
    }
  }, [createMockHistory, saveHistory])

  const applyFilters = useCallback(() => {
    let filtered = [...history]

    if (filters.search) {
      const searchLower = filters.search.toLowerCase()
      filtered = filtered.filter(item => 
        (item.title?.toLowerCase().includes(searchLower)) ||
        item.inputContent.toLowerCase().includes(searchLower) ||
        (item.outputContent?.toLowerCase().includes(searchLower))
      )
    }

    if (filters.operation !== 'all') {
      filtered = filtered.filter(item => item.operation === filters.operation)
    }

    if (filters.status !== 'all') {
      filtered = filtered.filter(item => 
        filters.status === 'success' ? item.success : !item.success
      )
    }

    if (filters.dateRange !== 'all') {
      const now = new Date()
      const cutoff = new Date()
      
      switch (filters.dateRange) {
        case 'today':
          cutoff.setHours(0, 0, 0, 0)
          break
        case 'week':
          cutoff.setDate(now.getDate() - 7)
          break
        case 'month':
          cutoff.setMonth(now.getMonth() - 1)
          break
      }
      
      filtered = filtered.filter(item => item.timestamp >= cutoff)
    }

    filtered.sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime())
    setFilteredHistory(filtered)
  }, [history, filters])

  useEffect(() => {
    loadHistory()
  }, [loadHistory])

  useEffect(() => {
    applyFilters()
  }, [applyFilters])

  const handleViewDetails = (item: ConversionHistoryItem) => {
    setSelectedItem(item)
    setIsDetailDialogOpen(true)
  }

  const handleRestoreToEditor = (item: ConversionHistoryItem) => {
    updateMarkdown(item.inputContent)
    if (item.outputContent) {
      updateSpec(item.outputContent)
      setFormat(item.outputFormat)
    }
    toast.success('Restored to editor')
  }

  const handleClearHistory = () => {
    setHistory([])
    saveHistory([])
    toast.success('History cleared')
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-heading-lg font-semibold text-text-primary">
              Conversion History
            </h1>
            <p className="text-body-md text-text-secondary mt-2">
              View and manage your API documentation conversion history.
            </p>
          </div>
          
          <div className="flex items-center space-x-2">
            <Button variant="outline" onClick={loadHistory}>
              <HistoryIcon className="h-4 w-4 mr-2" />
              Refresh
            </Button>
            <Button 
              variant="destructive" 
              onClick={handleClearHistory}
              disabled={history.length === 0}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Clear All
            </Button>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-4 p-4 bg-muted/20 rounded-lg">
          <div className="flex items-center space-x-2 flex-1 min-w-64">
            <Search className="h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search conversions..."
              value={filters.search}
              onChange={(e) => setFilters(prev => ({ ...prev, search: e.target.value }))}
              className="flex-1"
            />
          </div>
        </div>
      </div>

      {filteredHistory.length === 0 ? (
        <Card>
          <CardContent className="flex items-center justify-center py-16">
            <div className="text-center">
              <HistoryIcon className="h-16 w-16 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-heading-sm font-medium mb-2">No conversion history</h3>
              <p className="text-body-sm text-muted-foreground max-w-sm">
                Start converting API documentation to see your history here.
              </p>
            </div>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-4">
          {filteredHistory.map((item) => (
            <Card key={item.id} className="hover:shadow-md transition-shadow">
              <CardContent className="p-4">
                <div className="flex items-start justify-between">
                  <div className="flex items-start space-x-4 flex-1">
                    <div className="p-2 rounded-full bg-blue-500">
                      <FileText className="h-4 w-4" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center space-x-2 mb-2">
                        <h3 className="text-heading-sm font-medium truncate">
                          {item.title || `${item.operation} - ${item.id}`}
                        </h3>
                        <Badge variant={item.success ? 'default' : 'destructive'} className="text-xs">
                          {item.success ? (
                            <>
                              <CheckCircle className="h-3 w-3 mr-1" />
                              Success
                            </>
                          ) : (
                            <>
                              <AlertCircle className="h-3 w-3 mr-1" />
                              Error
                            </>
                          )}
                        </Badge>
                      </div>
                      <div className="flex items-center space-x-4 text-sm text-muted-foreground mb-2">
                        <span className="flex items-center">
                          <Calendar className="h-3 w-3 mr-1" />
                          {item.timestamp.toLocaleString()}
                        </span>
                      </div>
                      <p className="text-sm text-muted-foreground truncate">
                        {item.inputContent.split('\n')[0].substring(0, 100)}...
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2 ml-4">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleViewDetails(item)}
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleRestoreToEditor(item)}
                    >
                      <HistoryIcon className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        const blob = new Blob([item.outputContent || item.inputContent])
                        const url = URL.createObjectURL(blob)
                        const a = document.createElement('a')
                        a.href = url
                        a.download = `conversion_${item.id}.txt`
                        a.click()
                        URL.revokeObjectURL(url)
                      }}
                    >
                      <Download className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <Dialog open={isDetailDialogOpen} onOpenChange={setIsDetailDialogOpen}>
        <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {selectedItem?.title || `Conversion Details - ${selectedItem?.id}`}
            </DialogTitle>
          </DialogHeader>
          
          {selectedItem && (
            <div className="space-y-4">
              <Card>
                <CardHeader>
                  <CardTitle className="text-base">Input Content</CardTitle>
                </CardHeader>
                <CardContent>
                  <MonacoEditor
                    value={selectedItem.inputContent}
                    language={selectedItem.inputType === 'markdown' ? 'markdown' : 'yaml'}
                    height="200px"
                    readOnly
                    options={{
                      minimap: { enabled: false },
                      readOnly: true,
                      scrollBeyondLastLine: false
                    }}
                  />
                </CardContent>
              </Card>

              {selectedItem.outputContent && (
                <Card>
                  <CardHeader>
                    <CardTitle className="text-base">Output Content</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <MonacoEditor
                      value={selectedItem.outputContent}
                      language={selectedItem.outputFormat}
                      height="200px"
                      readOnly
                      options={{
                        minimap: { enabled: false },
                        readOnly: true,
                        scrollBeyondLastLine: false
                      }}
                    />
                  </CardContent>
                </Card>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}
