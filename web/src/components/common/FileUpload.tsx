import { useRef, useState, useCallback } from 'react'
import { Upload, X, FileText, AlertCircle } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'

interface FileUploadProps {
  onFileSelect: (file: File | null) => void
  accept?: string
  maxSize?: number // in bytes
  className?: string
  disabled?: boolean
  multiple?: boolean
}

interface UploadedFile {
  file: File
  progress: number
  error?: string
}

export default function FileUpload({
  onFileSelect,
  accept = '.md,.markdown,.txt',
  maxSize = 5 * 1024 * 1024, // 5MB default
  className,
  disabled = false,
  multiple = false
}: FileUploadProps) {
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [isDragOver, setIsDragOver] = useState(false)
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([])
  const [error, setError] = useState<string | null>(null)

  const validateFile = useCallback((file: File): string | null => {
    // Check file size
    if (file.size > maxSize) {
      return `File size must be less than ${(maxSize / 1024 / 1024).toFixed(1)}MB`
    }

    // Check file type
    const allowedTypes = accept.split(',').map(type => type.trim().toLowerCase())
    const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase()
    const mimeType = file.type.toLowerCase()

    const isValidExtension = allowedTypes.some(type => 
      type === fileExtension || (type.startsWith('.') && fileExtension === type)
    )
    const isValidMimeType = allowedTypes.some(type => 
      type === mimeType || mimeType.includes(type.replace('.', ''))
    )

    if (!isValidExtension && !isValidMimeType) {
      return `File type not supported. Allowed types: ${accept}`
    }

    return null
  }, [accept, maxSize])

  const handleFileSelect = useCallback((files: FileList | null) => {
    if (!files || files.length === 0) return

    const newFiles: UploadedFile[] = []
    let hasError = false

    Array.from(files).forEach(file => {
      const validationError = validateFile(file)
      if (validationError) {
        setError(validationError)
        hasError = true
        return
      }

      newFiles.push({
        file,
        progress: 0
      })
    })

    if (hasError) return

    setError(null)
    
    if (multiple) {
      setUploadedFiles(prev => [...prev, ...newFiles])
    } else {
      setUploadedFiles(newFiles)
      if (newFiles.length > 0) {
        onFileSelect(newFiles[0].file)
      }
    }

    // Simulate upload progress
    newFiles.forEach((uploadedFile) => {
      const interval = setInterval(() => {
        setUploadedFiles(prev => 
          prev.map(f => 
            f.file === uploadedFile.file 
              ? { ...f, progress: Math.min(f.progress + 10, 100) }
              : f
          )
        )
        
        if (uploadedFile.progress >= 90) {
          clearInterval(interval)
        }
      }, 100)
    })
  }, [validateFile, onFileSelect, multiple])

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (!disabled) {
      setIsDragOver(true)
    }
  }, [disabled])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragOver(false)
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragOver(false)
    
    if (disabled) return

    const files = e.dataTransfer.files
    handleFileSelect(files)
  }, [disabled, handleFileSelect])

  const removeFile = useCallback((fileToRemove: File) => {
    setUploadedFiles(prev => prev.filter(f => f.file !== fileToRemove))
    if (!multiple) {
      onFileSelect(null)
    }
  }, [multiple, onFileSelect])

  const openFileDialog = () => {
    if (!disabled && fileInputRef.current) {
      fileInputRef.current.click()
    }
  }

  return (
    <div className={cn('space-y-4', className)}>
      {/* Upload Area */}
      <div
        className={cn(
          'border-2 border-dashed rounded-lg p-8 text-center transition-colors cursor-pointer',
          isDragOver && !disabled
            ? 'border-primary bg-primary/5'
            : error
            ? 'border-error bg-error/5'
            : 'border-border bg-surface hover:bg-surface-elevated',
          disabled && 'opacity-50 cursor-not-allowed'
        )}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={openFileDialog}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept={accept}
          multiple={multiple}
          onChange={(e) => handleFileSelect(e.target.files)}
          className="hidden"
          disabled={disabled}
        />

        <div className="space-y-4">
          <div className="mx-auto w-12 h-12 rounded-full bg-muted flex items-center justify-center">
            <Upload className="h-6 w-6 text-text-secondary" />
          </div>
          
          <div className="space-y-2">
            <h3 className="text-heading-sm font-semibold text-text-primary">
              Drop files here or click to browse
            </h3>
            <p className="text-body-sm text-text-secondary">
              Supports: {accept.replace(/\./g, '').toUpperCase()} files up to{' '}
              {(maxSize / 1024 / 1024).toFixed(1)}MB
            </p>
          </div>

          <Button variant="outline" size="sm" disabled={disabled}>
            <Upload className="h-4 w-4 mr-2" />
            Choose Files
          </Button>
        </div>
      </div>

      {/* Error Display */}
      {error && (
        <div className="flex items-center space-x-2 p-3 bg-error/10 border border-error/20 rounded-md">
          <AlertCircle className="h-4 w-4 text-error" />
          <span className="text-body-sm text-error">{error}</span>
        </div>
      )}

      {/* Uploaded Files */}
      {uploadedFiles.length > 0 && (
        <div className="space-y-2">
          <h4 className="text-label-md font-medium text-text-primary">
            {multiple ? 'Uploaded Files' : 'Selected File'}
          </h4>
          {uploadedFiles.map((uploadedFile, index) => (
            <div
              key={`${uploadedFile.file.name}-${index}`}
              className="flex items-center space-x-3 p-3 border border-border rounded-md bg-surface"
            >
              <FileText className="h-5 w-5 text-text-secondary flex-shrink-0" />
              
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between">
                  <span className="text-body-sm font-medium text-text-primary truncate">
                    {uploadedFile.file.name}
                  </span>
                  <span className="text-body-xs text-text-tertiary ml-2">
                    {(uploadedFile.file.size / 1024).toFixed(1)} KB
                  </span>
                </div>
                
                {uploadedFile.progress < 100 && (
                  <div className="mt-2">
                    <Progress value={uploadedFile.progress} className="h-1" />
                  </div>
                )}
                
                {uploadedFile.error && (
                  <div className="mt-1 text-body-xs text-error">
                    {uploadedFile.error}
                  </div>
                )}
              </div>

              <Button
                variant="ghost"
                size="icon"
                onClick={(e) => {
                  e.stopPropagation()
                  removeFile(uploadedFile.file)
                }}
                className="h-8 w-8 text-text-tertiary hover:text-error"
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}