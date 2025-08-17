import { useState } from 'react'
import { toast } from 'sonner'
import { GitMerge, FileText, Upload, Download, Copy, Eye, Play, RotateCcw } from 'lucide-react'

// Components
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import WorkspaceLayout from '@/components/layout/WorkspaceLayout'
import MonacoEditor from '@/components/editor/MonacoEditor'
import FileUpload from '@/components/common/FileUpload'

// Hooks and Store
import { useAmendMutation, useDiffMutation } from '@/hooks/useApiQueries'
import { useAppStore } from '@/store/useAppStore'
import type { AmendResponse, DiffResponse, DiffLine } from '@/types/api'

export default function Amend() {
  const { 
    editor, 
    updateSpec, 
    setFormat,
    markAsModified 
  } = useAppStore()

  const [originalSpec, setOriginalSpec] = useState('')
  const [changes, setChanges] = useState('')
  const [amendedSpec, setAmendedSpec] = useState('')
  const [activeTab, setActiveTab] = useState<'original' | 'changes' | 'upload'>('original')
  const [outputTab, setOutputTab] = useState<'amended' | 'diff' | 'preview'>('amended')
  const [dryRun, setDryRun] = useState(true)

  // API hooks
  const amendMutation = useAmendMutation()
  const diffMutation = useDiffMutation()

  const handleFileSelect = (file: File | null, target: 'original' | 'changes') => {
    if (!file) return

    const reader = new FileReader()
    reader.onload = (e) => {
      const content = e.target?.result as string
      if (target === 'original') {
        setOriginalSpec(content)
      } else {
        setChanges(content)
      }
      markAsModified()
      toast.success(`Loaded ${file.name}`)
    }
    reader.onerror = () => {
      toast.error('Failed to read file')
    }
    reader.readAsText(file)
  }

  const handleOriginalSpecChange = (value: string | undefined) => {
    const newSpec = value || ''
    setOriginalSpec(newSpec)
    markAsModified()
  }

  const handleChangesChange = (value: string | undefined) => {
    const newChanges = value || ''
    setChanges(newChanges)
    markAsModified()
  }

  const handlePreviewAmendment = async () => {
    if (!originalSpec.trim()) {
      toast.error('Please provide an original OpenAPI specification')
      return
    }
    
    if (!changes.trim()) {
      toast.error('Please provide changes to apply')
      return
    }

    try {
      const result = await amendMutation.mutateAsync({
        originalSpec,
        changes,
        format: editor.format,
        dryRun: true
      })
      
      if (result.amendedSpec) {
        setAmendedSpec(result.amendedSpec)
        setOutputTab('amended')
        toast.success('Amendment preview generated')
      }
    } catch (error) {
      toast.error('Failed to preview amendment')
    }
  }

  const handleApplyAmendment = async () => {
    if (!originalSpec.trim() || !changes.trim()) {
      toast.error('Please provide both original spec and changes')
      return
    }

    try {
      const result = await amendMutation.mutateAsync({
        originalSpec,
        changes,
        format: editor.format,
        dryRun: false
      })
      
      if (result.amendedSpec) {
        setAmendedSpec(result.amendedSpec)
        updateSpec(result.amendedSpec)
        setOutputTab('amended')
        toast.success('Amendment applied successfully')
      }
    } catch (error) {
      toast.error('Failed to apply amendment')
    }
  }

  const handleGenerateDiff = async () => {
    if (!originalSpec.trim() || !amendedSpec.trim()) {
      toast.error('Please provide both original and amended specifications')
      return
    }

    try {
      await diffMutation.mutateAsync({
        originalSpec,
        modifiedSpec: amendedSpec,
        format: editor.format
      })
      
      setOutputTab('diff')
      toast.success('Diff generated')
    } catch (error) {
      toast.error('Failed to generate diff')
    }
  }

  const handleCopyResult = async () => {
    const textToCopy = outputTab === 'amended' ? amendedSpec : originalSpec
    try {
      await navigator.clipboard.writeText(textToCopy)
      toast.success('Specification copied to clipboard')
    } catch (error) {
      toast.error('Failed to copy to clipboard')
    }
  }

  const handleDownloadResult = () => {
    const textToDownload = outputTab === 'amended' ? amendedSpec : originalSpec
    const filename = outputTab === 'amended' ? 'amended-spec' : 'original-spec'
    
    const blob = new Blob([textToDownload], { 
      type: editor.format === 'yaml' ? 'text/yaml' : 'application/json' 
    })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${filename}.${editor.format}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    toast.success('Specification downloaded')
  }

  const handleReset = () => {
    setOriginalSpec('')
    setChanges('')
    setAmendedSpec('')
    setActiveTab('original')
    setOutputTab('amended')
    amendMutation.reset()
    diffMutation.reset()
    toast.success('Amendment workspace reset')
  }

  const isAmending = amendMutation.isPending
  const isGeneratingDiff = diffMutation.isPending
  const amendResult = amendMutation.data
  const diffResult = diffMutation.data

  // Left Panel - Input
  const leftPanel = (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b border-border">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-heading-sm font-semibold">Input Specifications</h2>
          <div className="flex items-center space-x-2">
            <Select value={editor.format} onValueChange={(value: 'yaml' | 'json') => setFormat(value)}>
              <SelectTrigger className="w-20">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="yaml">YAML</SelectItem>
                <SelectItem value="json">JSON</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <Tabs value={activeTab} onValueChange={(value: 'original' | 'changes' | 'upload') => setActiveTab(value)}>
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="original">Original</TabsTrigger>
            <TabsTrigger value="changes">Changes</TabsTrigger>
            <TabsTrigger value="upload">Upload</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      <div className="flex-1 overflow-hidden">
        {activeTab === 'original' && (
          <div className="h-full p-4">
            <div className="mb-2">
              <label className="text-sm font-medium text-muted-foreground">Original OpenAPI Specification</label>
            </div>
            <MonacoEditor
              value={originalSpec}
              onChange={handleOriginalSpecChange}
              language={editor.format === 'yaml' ? 'yaml' : 'json'}
              height="100%"
              options={{
                minimap: { enabled: false },
                scrollBeyondLastLine: false,
                fontSize: 14,
                lineNumbers: 'on',
                wordWrap: 'on'
              }}
            />
          </div>
        )}

        {activeTab === 'changes' && (
          <div className="h-full p-4">
            <div className="mb-2">
              <label className="text-sm font-medium text-muted-foreground">Changes to Apply</label>
            </div>
            <MonacoEditor
              value={changes}
              onChange={handleChangesChange}
              language="markdown"
              height="100%"
              options={{
                minimap: { enabled: false },
                scrollBeyondLastLine: false,
                fontSize: 14,
                lineNumbers: 'on',
                wordWrap: 'on'
              }}
            />
          </div>
        )}

        {activeTab === 'upload' && (
          <div className="h-full p-4 space-y-4">
            <div>
              <label className="text-sm font-medium text-muted-foreground mb-2 block">Upload Original Spec</label>
              <FileUpload
                onFileSelect={(file) => handleFileSelect(file, 'original')}
                accept=".yaml,.yml,.json"
                maxSize={10 * 1024 * 1024}
                className="h-32"
              />
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground mb-2 block">Upload Changes</label>
              <FileUpload
                onFileSelect={(file) => handleFileSelect(file, 'changes')}
                accept=".md,.markdown,.txt"
                maxSize={5 * 1024 * 1024}
                className="h-32"
              />
            </div>
          </div>
        )}
      </div>

      <div className="p-4 border-t border-border">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center space-x-2">
            <label className="flex items-center space-x-2 text-sm">
              <input
                type="checkbox"
                checked={dryRun}
                onChange={(e) => setDryRun(e.target.checked)}
                className="rounded border-border"
              />
              <span>Dry run (preview only)</span>
            </label>
          </div>
          
          <Button
            variant="outline"
            size="sm"
            onClick={handleReset}
          >
            <RotateCcw className="h-4 w-4 mr-2" />
            Reset
          </Button>
        </div>
        
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            onClick={handlePreviewAmendment}
            disabled={!originalSpec || !changes || isAmending}
            className="flex-1"
          >
            <Eye className="h-4 w-4 mr-2" />
            Preview
          </Button>
          <Button
            onClick={handleApplyAmendment}
            disabled={!originalSpec || !changes || isAmending || dryRun}
            className="flex-1"
          >
            {isAmending ? (
              <div className="flex items-center">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                Amending...
              </div>
            ) : (
              <>
                <GitMerge className="h-4 w-4 mr-2" />
                Apply
              </>
            )}
          </Button>
        </div>
      </div>
    </div>
  )

  // Right Panel - Results
  const rightPanel = (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b border-border">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-heading-sm font-semibold">Amendment Results</h2>
          {amendResult && (
            <div className="flex items-center space-x-2">
              <Badge variant="default" className="bg-green-500">
                <GitMerge className="h-3 w-3 mr-1" />
                Amended
              </Badge>
            </div>
          )}
        </div>

        <Tabs value={outputTab} onValueChange={(value: 'amended' | 'diff' | 'preview') => setOutputTab(value)}>
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="amended">Amended</TabsTrigger>
            <TabsTrigger value="diff">Diff</TabsTrigger>
            <TabsTrigger value="preview">Preview</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      <div className="flex-1 overflow-auto">
        <AmendmentResults
          amendResult={amendResult}
          diffResult={diffResult}
          amendedSpec={amendedSpec}
          isAmending={isAmending}
          isGeneratingDiff={isGeneratingDiff}
          activeTab={outputTab}
          onGenerateDiff={handleGenerateDiff}
        />
      </div>

      <div className="p-4 border-t border-border">
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleCopyResult}
            disabled={!amendedSpec}
          >
            <Copy className="h-4 w-4 mr-2" />
            Copy
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={handleDownloadResult}
            disabled={!amendedSpec}
          >
            <Download className="h-4 w-4 mr-2" />
            Download
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={handleGenerateDiff}
            disabled={!originalSpec || !amendedSpec || isGeneratingDiff}
            className="flex-1"
          >
            {isGeneratingDiff ? (
              <div className="flex items-center">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary mr-2" />
                Generating...
              </div>
            ) : (
              <>
                <Play className="h-4 w-4 mr-2" />
                Generate Diff
              </>
            )}
          </Button>
        </div>
      </div>
    </div>
  )

  return (
    <WorkspaceLayout
      leftPanel={leftPanel}
      centerPanel={
        <div className="h-full flex items-center justify-center bg-muted/20">
          <div className="text-center">
            <GitMerge className="h-16 w-16 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-heading-sm font-medium mb-2">OpenAPI Amendment</h3>
            <p className="text-body-sm text-muted-foreground max-w-sm">
              Apply changes to existing OpenAPI specifications and preview the differences before committing.
            </p>
          </div>
        </div>
      }
      rightPanel={rightPanel}
      leftWidth={35}
      centerWidth={30}
      rightWidth={35}
    />
  )
}

// Amendment Results Component
interface AmendmentResultsProps {
  amendResult?: AmendResponse
  diffResult?: DiffResponse
  amendedSpec: string
  isAmending: boolean
  isGeneratingDiff: boolean
  activeTab: 'amended' | 'diff' | 'preview'
  onGenerateDiff: () => void
}

function AmendmentResults({ 
  amendResult, 
  diffResult, 
  amendedSpec,
  isAmending, 
  isGeneratingDiff,
  activeTab, 
  onGenerateDiff 
}: AmendmentResultsProps) {
  if (isAmending) {
    return (
      <div className="flex items-center justify-center h-32 p-4">
        <div className="flex items-center space-x-2">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary" />
          <span className="text-sm text-muted-foreground">Applying amendments...</span>
        </div>
      </div>
    )
  }

  if (activeTab === 'amended') {
    return <AmendedSpecView spec={amendedSpec} />
  }

  if (activeTab === 'diff') {
    return <DiffView diffResult={diffResult} isGenerating={isGeneratingDiff} onGenerate={onGenerateDiff} />
  }

  if (activeTab === 'preview') {
    return <AmendmentPreview amendResult={amendResult} />
  }

  return null
}

// Amended Spec View Component
function AmendedSpecView({ spec }: { spec: string }) {
  if (!spec) {
    return (
      <div className="flex items-center justify-center h-32 p-4">
        <div className="text-center">
          <FileText className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">No amended specification yet</p>
          <p className="text-xs text-muted-foreground mt-1">
            Apply amendments to see the result
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="h-full p-4">
      <MonacoEditor
        value={spec}
        language="yaml"
        height="100%"
        readOnly
        options={{
          minimap: { enabled: false },
          scrollBeyondLastLine: false,
          fontSize: 14,
          lineNumbers: 'on',
          wordWrap: 'on',
          readOnly: true
        }}
      />
    </div>
  )
}

// Diff View Component
function DiffView({ 
  diffResult, 
  isGenerating, 
  onGenerate 
}: { 
  diffResult?: DiffResponse
  isGenerating: boolean
  onGenerate: () => void 
}) {
  if (isGenerating) {
    return (
      <div className="flex items-center justify-center h-32 p-4">
        <div className="flex items-center space-x-2">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary" />
          <span className="text-sm text-muted-foreground">Generating diff...</span>
        </div>
      </div>
    )
  }

  if (!diffResult) {
    return (
      <div className="flex items-center justify-center h-32 p-4">
        <div className="text-center">
          <GitMerge className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">No diff generated yet</p>
          <Button variant="outline" onClick={onGenerate} className="mt-2">
            Generate Diff
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="p-4 space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Diff Summary</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <div className="text-xl font-semibold text-green-500">{diffResult.summary.linesAdded}</div>
              <div className="text-sm text-muted-foreground">Added</div>
            </div>
            <div>
              <div className="text-xl font-semibold text-red-500">{diffResult.summary.linesRemoved}</div>
              <div className="text-sm text-muted-foreground">Removed</div>
            </div>
            <div>
              <div className="text-xl font-semibold text-yellow-500">{diffResult.summary.linesModified}</div>
              <div className="text-sm text-muted-foreground">Modified</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Changes</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2 max-h-96 overflow-y-auto">
            {diffResult.diff.added.map((line, index) => (
              <DiffLineComponent key={`added-${index}`} line={line} />
            ))}
            {diffResult.diff.removed.map((line, index) => (
              <DiffLineComponent key={`removed-${index}`} line={line} />
            ))}
            {diffResult.diff.modified.map((line, index) => (
              <DiffLineComponent key={`modified-${index}`} line={line} />
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

// Diff Line Component
function DiffLineComponent({ line }: { line: DiffLine }) {
  const getLineStyle = (type: string) => {
    switch (type) {
      case 'added':
        return 'bg-green-50 border-green-200 text-green-800'
      case 'removed':
        return 'bg-red-50 border-red-200 text-red-800'
      case 'modified':
        return 'bg-yellow-50 border-yellow-200 text-yellow-800'
      default:
        return 'bg-muted/20 border-border'
    }
  }

  const getPrefix = (type: string) => {
    switch (type) {
      case 'added':
        return '+'
      case 'removed':
        return '-'
      case 'modified':
        return '~'
      default:
        return ' '
    }
  }

  return (
    <div className={`p-2 rounded border font-mono text-sm ${getLineStyle(line.type)}`}>
      <div className="flex items-start space-x-2">
        <span className="text-xs text-muted-foreground min-w-8">{line.lineNumber}</span>
        <span className="font-bold">{getPrefix(line.type)}</span>
        <span className="flex-1">{line.content}</span>
      </div>
      {line.path && (
        <div className="text-xs text-muted-foreground mt-1 ml-10">
          Path: {line.path}
        </div>
      )}
    </div>
  )
}

// Amendment Preview Component
function AmendmentPreview({ amendResult }: { amendResult?: AmendResponse }) {
  if (!amendResult) {
    return (
      <div className="flex items-center justify-center h-32 p-4">
        <div className="text-center">
          <Eye className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">No preview available</p>
          <p className="text-xs text-muted-foreground mt-1">
            Generate an amendment to see the preview
          </p>
        </div>
      </div>
    )
  }

  const errorCount = amendResult.errors?.length || 0
  const warningCount = amendResult.warnings?.length || 0

  return (
    <div className="p-4 space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Amendment Status</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-red-500">{errorCount}</div>
              <div className="text-sm text-muted-foreground">Errors</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-yellow-500">{warningCount}</div>
              <div className="text-sm text-muted-foreground">Warnings</div>
            </div>
          </div>
        </CardContent>
      </Card>

      {(amendResult.errors?.length || 0) > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Errors</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {amendResult.errors?.map((error, index) => (
                <div key={index} className="p-2 bg-red-50 border border-red-200 rounded text-sm">
                  {error.message}
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {(amendResult.warnings?.length || 0) > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Warnings</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {amendResult.warnings?.map((warning, index) => (
                <div key={index} className="p-2 bg-yellow-50 border border-yellow-200 rounded text-sm">
                  {warning.message}
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}