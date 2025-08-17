import { useState } from 'react'
import { toast } from 'sonner'
import { CheckCircle, AlertCircle, FileText, Upload, Download, Copy } from 'lucide-react'

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
import { useValidateMutation, useRealtimeValidation } from '@/hooks/useApiQueries'
import { useAppStore } from '@/store/useAppStore'
import type { ValidationError, ValidationSummary } from '@/types/api'

export default function Validate() {
  const { 
    editor, 
    updateSpec, 
    setFormat,
    markAsModified 
  } = useAppStore()

  const [inputSpec, setInputSpec] = useState('')
  const [activeTab, setActiveTab] = useState<'editor' | 'upload'>('editor')
  const [outputTab, setOutputTab] = useState<'results' | 'summary' | 'errors'>('results')
  const [realtimeValidation, setRealtimeValidation] = useState(true)

  // API hooks
  const validateMutation = useValidateMutation()
  const realtimeQuery = useRealtimeValidation(
    realtimeValidation ? inputSpec : '', 
    editor.format
  )

  const handleFileSelect = (file: File | null) => {
    if (!file) return

    const reader = new FileReader()
    reader.onload = (e) => {
      const content = e.target?.result as string
      setInputSpec(content)
      updateSpec(content)
      markAsModified()
      toast.success(`Loaded ${file.name}`)
    }
    reader.onerror = () => {
      toast.error('Failed to read file')
    }
    reader.readAsText(file)
  }

  const handleSpecChange = (value: string | undefined) => {
    const newSpec = value || ''
    setInputSpec(newSpec)
    updateSpec(newSpec)
    markAsModified()
  }

  const handleManualValidate = async () => {
    if (!inputSpec.trim()) {
      toast.error('Please provide an OpenAPI specification to validate')
      return
    }

    try {
      await validateMutation.mutateAsync({
        spec: inputSpec,
        format: editor.format
      })
      toast.success('Validation completed')
    } catch (error) {
      toast.error('Validation failed')
    }
  }

  const handleCopySpec = async () => {
    try {
      await navigator.clipboard.writeText(inputSpec)
      toast.success('Specification copied to clipboard')
    } catch (error) {
      toast.error('Failed to copy to clipboard')
    }
  }

  const handleDownloadSpec = () => {
    const blob = new Blob([inputSpec], { 
      type: editor.format === 'yaml' ? 'text/yaml' : 'application/json' 
    })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `openapi-spec.${editor.format}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    toast.success('Specification downloaded')
  }

  // Get validation results from either realtime or manual validation
  const validationResult = realtimeValidation && realtimeQuery.data 
    ? realtimeQuery.data 
    : validateMutation.data

  const isValidating = validateMutation.isPending || 
    (realtimeValidation && realtimeQuery.isFetching)

  const hasErrors = validationResult && !validationResult.valid
  const errorCount = validationResult?.errors?.length || 0
  const warningCount = validationResult?.warnings?.length || 0

  // Left Panel - Input
  const leftPanel = (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b border-border">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-heading-sm font-semibold">OpenAPI Specification</h2>
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

        <Tabs value={activeTab} onValueChange={(value: 'editor' | 'upload') => setActiveTab(value)}>
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="editor">Editor</TabsTrigger>
            <TabsTrigger value="upload">Upload</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      <div className="flex-1 overflow-hidden">
        {activeTab === 'editor' && (
          <div className="h-full p-4">
            <MonacoEditor
              value={inputSpec}
              onChange={handleSpecChange}
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

        {activeTab === 'upload' && (
          <div className="h-full p-4 flex items-center justify-center">
            <FileUpload
              onFileSelect={handleFileSelect}
              accept=".yaml,.yml,.json"
              maxSize={10 * 1024 * 1024} // 10MB
              className="w-full h-full"
            />
          </div>
        )}
      </div>

      <div className="p-4 border-t border-border">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <label className="flex items-center space-x-2 text-sm">
              <input
                type="checkbox"
                checked={realtimeValidation}
                onChange={(e) => setRealtimeValidation(e.target.checked)}
                className="rounded border-border"
              />
              <span>Real-time validation</span>
            </label>
          </div>
          
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleCopySpec}
              disabled={!inputSpec}
            >
              <Copy className="h-4 w-4 mr-2" />
              Copy
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={handleDownloadSpec}
              disabled={!inputSpec}
            >
              <Download className="h-4 w-4 mr-2" />
              Download
            </Button>
            <Button
              onClick={handleManualValidate}
              disabled={!inputSpec || isValidating}
              className="min-w-24"
            >
              {isValidating ? (
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                  Validating...
                </div>
              ) : (
                <>
                  <CheckCircle className="h-4 w-4 mr-2" />
                  Validate
                </>
              )}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )

  // Right Panel - Results
  const rightPanel = (
    <div className="h-full flex flex-col">
      <div className="p-4 border-b border-border">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-heading-sm font-semibold">Validation Results</h2>
          {validationResult && (
            <div className="flex items-center space-x-2">
              {validationResult.valid ? (
                <Badge variant="default" className="bg-green-500">
                  <CheckCircle className="h-3 w-3 mr-1" />
                  Valid
                </Badge>
              ) : (
                <Badge variant="destructive">
                  <AlertCircle className="h-3 w-3 mr-1" />
                  Invalid
                </Badge>
              )}
            </div>
          )}
        </div>

        <Tabs value={outputTab} onValueChange={(value: 'results' | 'summary' | 'errors') => setOutputTab(value)}>
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="results">Results</TabsTrigger>
            <TabsTrigger value="summary">Summary</TabsTrigger>
            <TabsTrigger value="errors" className="relative">
              Errors
              {errorCount > 0 && (
                <Badge variant="destructive" className="ml-1 h-5 w-5 rounded-full p-0 text-xs">
                  {errorCount}
                </Badge>
              )}
            </TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      <div className="flex-1 overflow-auto p-4">
        <ValidationResults
          result={validationResult}
          isLoading={isValidating}
          activeTab={outputTab}
          onErrorClick={(error) => {
            // TODO: Jump to error line in editor
            console.log('Jump to error:', error)
          }}
        />
      </div>
    </div>
  )

  return (
    <WorkspaceLayout
      leftPanel={leftPanel}
      centerPanel={
        <div className="h-full flex items-center justify-center bg-muted/20">
          <div className="text-center">
            <FileText className="h-16 w-16 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-heading-sm font-medium mb-2">OpenAPI Validator</h3>
            <p className="text-body-sm text-muted-foreground max-w-sm">
              Upload or paste your OpenAPI specification to validate it against the OpenAPI 3.1 standard.
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

// Validation Results Component
interface ValidationResultsProps {
  result?: any
  isLoading: boolean
  activeTab: 'results' | 'summary' | 'errors'
  onErrorClick: (error: ValidationError) => void
}

function ValidationResults({ result, isLoading, activeTab, onErrorClick }: ValidationResultsProps) {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-32">
        <div className="flex items-center space-x-2">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary" />
          <span className="text-sm text-muted-foreground">Validating specification...</span>
        </div>
      </div>
    )
  }

  if (!result) {
    return (
      <div className="flex items-center justify-center h-32">
        <div className="text-center">
          <AlertCircle className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
          <p className="text-sm text-muted-foreground">No validation results yet</p>
          <p className="text-xs text-muted-foreground mt-1">
            Enter a specification to see validation results
          </p>
        </div>
      </div>
    )
  }

  if (activeTab === 'results') {
    return <ValidationOverview result={result} />
  }

  if (activeTab === 'summary') {
    return <ValidationSummaryView summary={result.summary} />
  }

  if (activeTab === 'errors') {
    return <ValidationErrorsList errors={result.errors || []} warnings={result.warnings || []} onErrorClick={onErrorClick} />
  }

  return null
}

// Validation Overview Component
function ValidationOverview({ result }: { result: any }) {
  const errorCount = result.errors?.length || 0
  const warningCount = result.warnings?.length || 0

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            {result.valid ? (
              <CheckCircle className="h-5 w-5 text-green-500" />
            ) : (
              <AlertCircle className="h-5 w-5 text-red-500" />
            )}
            <span>Validation Status</span>
          </CardTitle>
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

      {result.summary && (
        <Card>
          <CardHeader>
            <CardTitle>Specification Summary</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-3 gap-4 text-center">
              <div>
                <div className="text-xl font-semibold">{result.summary.endpoints || 0}</div>
                <div className="text-sm text-muted-foreground">Endpoints</div>
              </div>
              <div>
                <div className="text-xl font-semibold">{result.summary.schemas || 0}</div>
                <div className="text-sm text-muted-foreground">Schemas</div>
              </div>
              <div>
                <div className="text-xl font-semibold">{result.summary.parameters || 0}</div>
                <div className="text-sm text-muted-foreground">Parameters</div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

// Validation Summary Component
function ValidationSummaryView({ summary }: { summary?: ValidationSummary }) {
  if (!summary) {
    return (
      <div className="text-center py-8">
        <p className="text-muted-foreground">No summary available</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Specification Details</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-sm font-medium">Total Errors:</span>
              <Badge variant={summary.totalErrors > 0 ? "destructive" : "default"}>
                {summary.totalErrors}
              </Badge>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-medium">Total Warnings:</span>
              <Badge variant={summary.totalWarnings > 0 ? "secondary" : "default"}>
                {summary.totalWarnings}
              </Badge>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-medium">Endpoints:</span>
              <span className="text-sm">{summary.endpoints}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-medium">Schemas:</span>
              <span className="text-sm">{summary.schemas}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-medium">Parameters:</span>
              <span className="text-sm">{summary.parameters}</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

// Validation Errors List Component
function ValidationErrorsList({ 
  errors, 
  warnings, 
  onErrorClick 
}: { 
  errors: ValidationError[]
  warnings: ValidationError[]
  onErrorClick: (error: ValidationError) => void 
}) {
  const allIssues = [
    ...errors.map(e => ({ ...e, type: 'error' as const })),
    ...warnings.map(w => ({ ...w, type: 'warning' as const }))
  ].sort((a, b) => (a.line || 0) - (b.line || 0))

  if (allIssues.length === 0) {
    return (
      <div className="text-center py-8">
        <CheckCircle className="h-8 w-8 text-green-500 mx-auto mb-2" />
        <p className="text-sm text-muted-foreground">No errors or warnings found</p>
      </div>
    )
  }

  return (
    <div className="space-y-2">
      {allIssues.map((issue, index) => (
        <Card 
          key={index}
          className={`cursor-pointer transition-colors hover:bg-muted/50 ${
            issue.type === 'error' ? 'border-red-200' : 'border-yellow-200'
          }`}
          onClick={() => onErrorClick(issue)}
        >
          <CardContent className="p-3">
            <div className="flex items-start space-x-3">
              <div className="flex-shrink-0 mt-0.5">
                {issue.type === 'error' ? (
                  <AlertCircle className="h-4 w-4 text-red-500" />
                ) : (
                  <AlertCircle className="h-4 w-4 text-yellow-500" />
                )}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center space-x-2 mb-1">
                  <Badge variant={issue.type === 'error' ? 'destructive' : 'secondary'} className="text-xs">
                    {issue.type.toUpperCase()}
                  </Badge>
                  {issue.line && (
                    <Badge variant="outline" className="text-xs">
                      Line {issue.line}{issue.column ? `:${issue.column}` : ''}
                    </Badge>
                  )}
                  {issue.code && (
                    <Badge variant="outline" className="text-xs">
                      {issue.code}
                    </Badge>
                  )}
                </div>
                <p className="text-sm text-foreground">{issue.message}</p>
                {issue.path && (
                  <p className="text-xs text-muted-foreground mt-1">Path: {issue.path}</p>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}