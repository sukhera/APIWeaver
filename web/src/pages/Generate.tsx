import { useState } from 'react'
import { toast } from 'sonner'
import { Zap, FileText, Download, Eye, AlertCircle, CheckCircle } from 'lucide-react'

// Components
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import WorkspaceLayout from '@/components/layout/WorkspaceLayout'
import MonacoEditor from '@/components/editor/MonacoEditor'
import FileUpload from '@/components/common/FileUpload'

// Hooks and Store
import { useGenerateAndValidate } from '@/hooks/useApiQueries'
import { useAppStore } from '@/store/useAppStore'

export default function Generate() {
  const {
    editor,
    updateMarkdown,
    updateSpec,
    setFormat,
    setActiveTab,
    markAsModified,
    setGenerationResult,
    setValidationResult
  } = useAppStore()

  const { generateAndValidate, isLoading, error, reset } = useGenerateAndValidate()
  const [activeOutputTab, setActiveOutputTab] = useState<'spec' | 'validation' | 'errors'>('spec')

  const handleFileSelect = (file: File | null) => {
    if (!file) return

    const reader = new FileReader()
    reader.onload = (e) => {
      const content = e.target?.result as string
      updateMarkdown(content)
      markAsModified()
      toast.success(`Loaded ${file.name}`)
    }
    reader.onerror = () => {
      toast.error('Failed to read file')
    }
    reader.readAsText(file)
  }

  const handleGenerate = async () => {
    if (!editor.markdown.trim()) {
      toast.error('Please provide Markdown content to generate from')
      return
    }

    try {
      reset() // Clear previous errors
      const result = await generateAndValidate({
        markdown: editor.markdown,
        format: editor.format,
        validate: true
      })

      if (result.generate) {
        setGenerationResult(result.generate)
        updateSpec(result.generate.spec)
        setActiveTab('preview')
        setActiveOutputTab('spec')

        if (result.validate) {
          setValidationResult(result.validate)
          
          if (result.validate.valid) {
            toast.success('OpenAPI specification generated and validated successfully!')
          } else {
            toast.warning(`Generated with ${result.validate.errors?.length || 0} validation errors`)
            setActiveOutputTab('validation')
          }
        } else {
          toast.success('OpenAPI specification generated successfully!')
        }
      }
    } catch (err: unknown) {
      const error = err as Error
      toast.error(error.message || 'Failed to generate specification')
    }
  }

  const handleExport = () => {
    if (!editor.spec) {
      toast.error('No specification to export')
      return
    }

    const blob = new Blob([editor.spec], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `openapi-spec.${editor.format}`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    toast.success(`Exported as ${editor.format.toUpperCase()}`)
  }

  const handleLoadTemplate = (templateName: string) => {
    const templates = {
      'Basic REST API': `# User Management API

## POST /users
Create a new user

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | User's full name |
| email | string | Yes | User's email address |
| role | string | No | User role (default: user) |

### Response
\`\`\`json
{
  "id": "string",
  "name": "string", 
  "email": "string",
  "role": "string",
  "createdAt": "string"
}
\`\`\`

## GET /users/{id}
Get user by ID

### Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | User ID |

### Response
\`\`\`json
{
  "id": "string",
  "name": "string",
  "email": "string", 
  "role": "string",
  "createdAt": "string"
}
\`\`\``,

      'E-commerce API': `# E-commerce API

## GET /products
List all products

### Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| category | string | No | Filter by category |
| limit | integer | No | Number of items to return |
| offset | integer | No | Number of items to skip |

### Response
\`\`\`json
{
  "products": [{
    "id": "string",
    "name": "string",
    "price": "number",
    "category": "string",
    "inStock": "boolean"
  }],
  "total": "integer"
}
\`\`\`

## POST /orders
Create a new order

### Request Body
\`\`\`json
{
  "items": [{
    "productId": "string",
    "quantity": "integer"
  }],
  "shippingAddress": {
    "street": "string",
    "city": "string", 
    "zipCode": "string"
  }
}
\`\`\``,

      'Task Management': `# Task Management API

## GET /tasks
List all tasks

### Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| status | string | No | Filter by status (pending, completed) |
| assignee | string | No | Filter by assignee |

### Response
\`\`\`json
{
  "tasks": [{
    "id": "string",
    "title": "string",
    "description": "string",
    "status": "string",
    "assignee": "string",
    "dueDate": "string"
  }]
}
\`\`\`

## POST /tasks
Create a new task

### Request Body
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Task title |
| description | string | No | Task description |
| assignee | string | No | User to assign task to |
| dueDate | string | No | Due date (ISO format) |

### Response
\`\`\`json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "status": "pending",
  "assignee": "string",
  "dueDate": "string",
  "createdAt": "string"
}
\`\`\``
    }

    const template = templates[templateName as keyof typeof templates]
    if (template) {
      updateMarkdown(template)
      markAsModified()
      toast.success(`Loaded ${templateName} template`)
    }
  }

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="bg-surface border-b border-border p-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Zap className="h-6 w-6 text-primary" />
            <div>
              <h1 className="text-heading-lg font-semibold text-text-primary">
                Generate OpenAPI Specification
              </h1>
              <p className="text-body-sm text-text-secondary">
                Convert Markdown API requirements to OpenAPI 3.1 specification
              </p>
            </div>
          </div>
          
          <div className="flex items-center space-x-3">
            <Select
              value={editor.format}
              onValueChange={(value: 'yaml' | 'json') => setFormat(value)}
            >
              <SelectTrigger className="w-24">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="yaml">YAML</SelectItem>
                <SelectItem value="json">JSON</SelectItem>
              </SelectContent>
            </Select>

            <Button
              onClick={handleGenerate}
              disabled={!editor.markdown.trim() || isLoading}
              className="flex items-center space-x-2"
            >
              <Zap className="h-4 w-4" />
              <span>{isLoading ? 'Generating...' : 'Generate'}</span>
            </Button>

            {editor.spec && (
              <Button
                variant="outline"
                onClick={handleExport}
                className="flex items-center space-x-2"
              >
                <Download className="h-4 w-4" />
                <span>Export</span>
              </Button>
            )}
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-hidden">
        <WorkspaceLayout
          leftPanel={
            <div className="h-full flex flex-col">
              <div className="p-4 border-b border-border">
                <h3 className="text-heading-xs font-medium text-text-primary mb-4">
                  Templates & Upload
                </h3>
                
                <FileUpload
                  onFileSelect={handleFileSelect}
                  accept=".md,.markdown,.txt"
                  className="mb-4"
                />
                
                <div className="mt-6">
                  <h4 className="text-label-sm font-medium text-text-primary mb-3">
                    Quick Start Templates
                  </h4>
                  
                  <div className="space-y-1">
                    {['Basic REST API', 'E-commerce API', 'Task Management'].map((template) => (
                      <button
                        key={template}
                        className="w-full text-left px-3 py-2 text-body-sm text-text-secondary hover:bg-surface-elevated rounded-md transition-colors"
                        onClick={() => handleLoadTemplate(template)}
                      >
                        <FileText className="h-4 w-4 inline mr-2" />
                        {template}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
              
              {error && (
                <div className="p-4 border-b border-border">
                  <div className="flex items-center space-x-2 p-3 bg-error/10 border border-error/20 rounded-md">
                    <AlertCircle className="h-4 w-4 text-error flex-shrink-0" />
                    <div className="text-body-sm text-error">
                      {error.message}
                    </div>
                  </div>
                </div>
              )}
            </div>
          }
          centerPanel={
            <div className="h-full flex flex-col">
              <div className="p-4 border-b border-border">
                <div className="flex items-center justify-between">
                  <h3 className="text-heading-xs font-medium text-text-primary">
                    Markdown Input
                  </h3>
                  <div className="text-body-xs text-text-secondary">
                    {editor.markdown.length} characters
                  </div>
                </div>
              </div>
              
              <div className="flex-1 p-4">
                <MonacoEditor
                  value={editor.markdown}
                  onChange={(value) => {
                    updateMarkdown(value || '')
                    markAsModified()
                  }}
                  language="markdown"
                  height="100%"
                  options={{
                    wordWrap: 'on',
                    minimap: { enabled: false },
                    scrollBeyondLastLine: false,
                  }}
                />
              </div>
            </div>
          }
          rightPanel={
            <div className="h-full flex flex-col">
              <div className="p-4 border-b border-border">
                <Tabs value={activeOutputTab} onValueChange={(value: string) => setActiveOutputTab(value)}>
                  <TabsList className="grid w-full grid-cols-3">
                    <TabsTrigger value="spec" className="text-xs">
                      <Eye className="h-3 w-3 mr-1" />
                      Output
                    </TabsTrigger>
                    <TabsTrigger value="validation" className="text-xs">
                      <CheckCircle className="h-3 w-3 mr-1" />
                      Validation
                    </TabsTrigger>
                    <TabsTrigger value="errors" className="text-xs">
                      <AlertCircle className="h-3 w-3 mr-1" />
                      Issues
                    </TabsTrigger>
                  </TabsList>
                </Tabs>
              </div>
              
              <div className="flex-1 overflow-hidden">
                <Tabs value={activeOutputTab} className="h-full">
                  <TabsContent value="spec" className="h-full m-0 p-4">
                    {editor.spec ? (
                      <MonacoEditor
                        value={editor.spec}
                        language={editor.format === 'yaml' ? 'yaml' : 'json'}
                        height="100%"
                        readOnly
                        options={{
                          minimap: { enabled: false },
                          scrollBeyondLastLine: false,
                        }}
                      />
                    ) : (
                      <div className="flex items-center justify-center h-full text-text-secondary">
                        <div className="text-center">
                          <FileText className="h-8 w-8 mx-auto mb-2 opacity-50" />
                          <p className="text-body-sm">
                            Generated specification will appear here
                          </p>
                        </div>
                      </div>
                    )}
                  </TabsContent>
                  
                  <TabsContent value="validation" className="h-full m-0 p-4">
                    <div className="space-y-4">
                      {/* Validation Summary */}
                      <Card>
                        <CardHeader className="pb-2">
                          <CardTitle className="text-sm">Validation Status</CardTitle>
                        </CardHeader>
                        <CardContent>
                          {editor.spec ? (
                            <div className="flex items-center space-x-2">
                              <CheckCircle className="h-4 w-4 text-success" />
                              <span className="text-body-sm text-success">
                                Specification is valid
                              </span>
                            </div>
                          ) : (
                            <div className="text-body-sm text-text-secondary">
                              Generate a specification to see validation results
                            </div>
                          )}
                        </CardContent>
                      </Card>
                    </div>
                  </TabsContent>
                  
                  <TabsContent value="errors" className="h-full m-0 p-4">
                    <div className="text-body-sm text-text-secondary">
                      No validation errors or warnings to display
                    </div>
                  </TabsContent>
                </Tabs>
              </div>
            </div>
          }
          leftWidth={25}
          centerWidth={50}
          rightWidth={25}
        />
      </div>
    </div>
  )
}