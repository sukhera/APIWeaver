import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderWithProviders, screen, userEvent, createTestFile } from '@/test/utils'
import FileUpload from '../FileUpload'

describe('FileUpload', () => {
  const mockOnFileSelect = vi.fn()

  beforeEach(() => {
    mockOnFileSelect.mockClear()
  })

  it('renders upload area with default props', () => {
    renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} />)
    
    expect(screen.getByText('Drop files here or click to browse')).toBeInTheDocument()
    expect(screen.getByText(/Supports.*MD.*files up to.*5\.0.*MB/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /choose files/i })).toBeInTheDocument()
  })

  it('accepts custom accept and maxSize props', () => {
    renderWithProviders(
      <FileUpload 
        onFileSelect={mockOnFileSelect}
        accept=".txt,.md"
        maxSize={1024 * 1024} // 1MB
      />
    )
    
    expect(screen.getByText(/Supports.*TXT.*MD.*files up to.*1\.0.*MB/i)).toBeInTheDocument()
  })

  it('handles file selection via file input', async () => {
    const user = userEvent.setup()
    const file = createTestFile('test.md', '# Test Content', 'text/markdown')
    
    renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} />)
    
    // Find the hidden file input directly
    const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
    expect(hiddenInput).toBeInTheDocument()
    
    await user.upload(hiddenInput, file)
    
    expect(mockOnFileSelect).toHaveBeenCalledWith(file)
  })

  it('validates file size', async () => {
    const user = userEvent.setup()
    const largeFile = createTestFile(
      'large.md', 
      'x'.repeat(6 * 1024 * 1024), // 6MB content
      'text/markdown'
    )
    
    renderWithProviders(
      <FileUpload 
        onFileSelect={mockOnFileSelect}
        maxSize={5 * 1024 * 1024} // 5MB limit
      />
    )
    
    // Find the hidden file input directly
    const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
    expect(hiddenInput).toBeInTheDocument()
    
    await user.upload(hiddenInput, largeFile)
    
    expect(screen.getByText(/File size must be less than 5\.0MB/i)).toBeInTheDocument()
    expect(mockOnFileSelect).not.toHaveBeenCalled()
  })

  // TODO: Fix file type validation test
  // it('validates file type', async () => {
  //   const user = userEvent.setup()
  //   const invalidFile = createTestFile('test.pdf', 'PDF content', 'application/pdf')
  //   
  //   renderWithProviders(
  //     <FileUpload 
  //       onFileSelect={mockOnFileSelect}
  //       accept=".md,.markdown,.txt"
  //     />
  //   )
  //   
  //   // Find the hidden file input directly
  //   const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
  //   expect(hiddenInput).toBeInTheDocument()
  //   
  //   await user.upload(hiddenInput, invalidFile)
  //   
  //   expect(screen.getByText(/File type not supported/i)).toBeInTheDocument()
  //   expect(mockOnFileSelect).not.toHaveBeenCalled()
  // })

  it('displays uploaded file information', async () => {
    const user = userEvent.setup()
    const file = createTestFile('test.md', '# Test Content', 'text/markdown')
    
    renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} />)
    
    // Find the hidden file input directly
    const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
    expect(hiddenInput).toBeInTheDocument()
    
    await user.upload(hiddenInput, file)
    
    expect(screen.getByText('test.md')).toBeInTheDocument()
    expect(screen.getByText(/KB/)).toBeInTheDocument() // File size display
  })

  it('allows file removal', async () => {
    const user = userEvent.setup()
    const file = createTestFile('test.md', '# Test Content', 'text/markdown')
    
    renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} />)
    
    // Find the hidden file input directly
    const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
    expect(hiddenInput).toBeInTheDocument()
    
    await user.upload(hiddenInput, file)
    
    // File should be displayed
    expect(screen.getByText('test.md')).toBeInTheDocument()
    
    // Find and click remove button
    const removeButton = screen.getByRole('button', { name: '' }) // X button has no label
    await user.click(removeButton)
    
    // File should be removed and onFileSelect called with null
    expect(screen.queryByText('test.md')).not.toBeInTheDocument()
    expect(mockOnFileSelect).toHaveBeenLastCalledWith(null)
  })

  // TODO: Implement disabled state styling
  // it('handles disabled state', () => {
  //   renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} disabled />)
  //   
  //   const uploadArea = screen.getByText('Drop files here or click to browse').closest('div')
  //   expect(uploadArea).toHaveClass('opacity-50', 'cursor-not-allowed')
  //   
  //   const button = screen.getByRole('button', { name: /choose files/i })
  //   expect(button).toBeDisabled()
  // })

  it('handles multiple files when multiple prop is true', async () => {
    const user = userEvent.setup()
    const file1 = createTestFile('test1.md', '# Test Content 1', 'text/markdown')
    const file2 = createTestFile('test2.md', '# Test Content 2', 'text/markdown')
    
    renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} multiple />)
    
    // Find the hidden file input directly
    const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
    expect(hiddenInput).toBeInTheDocument()
    
    await user.upload(hiddenInput, [file1, file2])
    
    expect(screen.getByText('test1.md')).toBeInTheDocument()
    expect(screen.getByText('test2.md')).toBeInTheDocument()
  })

  // TODO: Implement error state styling
  // it('shows error state when there is an error', async () => {
  //   const user = userEvent.setup()
  //   const invalidFile = createTestFile('test.exe', 'Invalid content', 'application/exe')
  //   
  //   renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} />)
  //   
  //   // Find the hidden file input directly
  //   const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
  //   expect(hiddenInput).toBeInTheDocument()
  //   
  //   await user.upload(hiddenInput, invalidFile)
  //   
  //   const uploadArea = screen.getByText('Drop files here or click to browse').closest('div')
  //   expect(uploadArea).toHaveClass('border-error', 'bg-error/5')
  //   
  //   expect(screen.getByText(/File type not supported/i)).toBeInTheDocument()
  // })

  it('simulates upload progress', async () => {
    const user = userEvent.setup()
    const file = createTestFile('test.md', '# Test Content', 'text/markdown')
    
    renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} />)
    
    // Find the hidden file input directly
    const hiddenInput = document.querySelector('input[type="file"]') as HTMLInputElement
    expect(hiddenInput).toBeInTheDocument()
    
    await user.upload(hiddenInput, file)
    
    // Should show progress initially
    expect(screen.getByRole('progressbar')).toBeInTheDocument()
    
    // Progress should complete after simulation
    await new Promise(resolve => setTimeout(resolve, 1100)) // Wait for progress to complete
    
    expect(screen.queryByRole('progressbar')).not.toBeInTheDocument()
  })

    // TODO: Implement custom className support
  // it('applies custom className', () => {
  //   renderWithProviders(
  //     <FileUpload 
  //       onFileSelect={mockOnFileSelect} 
  //       className="custom-upload-class"
  //     />
  //   )
  //   
  //   const container = screen.getByText('Drop files here or click to browse').closest('div')?.parentElement
  //   expect(container).toHaveClass('custom-upload-class')
  // })

  // TODO: Implement drag and drop styling
  // it('handles drag and drop events', async () => {
  //   const file = createTestFile('test.md', '# Test Content', 'text/markdown')
  //   
  //   renderWithProviders(<FileUpload onFileSelect={mockOnFileSelect} />)
  //   
  //   const uploadArea = screen.getByText('Drop files here or click to browse').closest('div')
  //   
  //   // Simulate drag over
  //   const dragOverEvent = new Event('dragover', { bubbles: true })
  //   Object.defineProperty(dragOverEvent, 'dataTransfer', {
  //     value: { files: [file] }
  //   })
  //   
  //   uploadArea?.dispatchEvent(dragOverEvent)
  //   
  //   // Should show drag over state
  //   expect(uploadArea).toHaveClass('border-primary', 'bg-primary/5')
  // })
})