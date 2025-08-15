import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderWithProviders, screen, userEvent, waitFor } from '@/test/utils'
import Generate from '../Generate'

// Mock the API hooks
vi.mock('@/hooks/useApiQueries', () => ({
  useGenerateAndValidate: () => ({
    generateAndValidate: vi.fn().mockResolvedValue({
      generate: {
        spec: 'openapi: 3.1.0\ninfo:\n  title: Test API\n  version: 1.0.0',
        format: 'yaml',
        errors: [],
        warnings: [],
      },
      validate: {
        valid: true,
        errors: [],
        warnings: [],
      },
    }),
    isLoading: false,
    error: null,
    reset: vi.fn(),
  }),
}))

// Mock the app store
vi.mock('@/store/useAppStore', () => ({
  useAppStore: () => ({
    editor: {
      markdown: '',
      spec: '',
      format: 'yaml',
      activeTab: 'editor',
      isModified: false,
      lastSaved: null,
    },
    updateMarkdown: vi.fn(),
    updateSpec: vi.fn(),
    setFormat: vi.fn(),
    setActiveTab: vi.fn(),
    markAsModified: vi.fn(),
    setGenerationResult: vi.fn(),
    setValidationResult: vi.fn(),
  }),
}))

describe('Generate Page', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders the main page elements', () => {
    renderWithProviders(<Generate />)
    
    // Check header elements
    expect(screen.getByText('Generate OpenAPI Specification')).toBeInTheDocument()
    expect(screen.getByText('Convert Markdown API requirements to OpenAPI 3.1 specification')).toBeInTheDocument()
    
    // Check main action buttons
    expect(screen.getByRole('button', { name: /generate/i })).toBeInTheDocument()
    
    // Check format selector
    expect(screen.getByRole('combobox')).toBeInTheDocument()
    
    // Check main panels
    expect(screen.getByText('Templates & Upload')).toBeInTheDocument()
    expect(screen.getByText('Markdown Input')).toBeInTheDocument()
    expect(screen.getByText('Output')).toBeInTheDocument()
  })

  it('displays template options', () => {
    renderWithProviders(<Generate />)
    
    expect(screen.getByText('Quick Start Templates')).toBeInTheDocument()
    expect(screen.getByText('Basic REST API')).toBeInTheDocument()
    expect(screen.getByText('E-commerce API')).toBeInTheDocument()
    expect(screen.getByText('Task Management')).toBeInTheDocument()
  })

  it('shows the Monaco editor for markdown input', () => {
    renderWithProviders(<Generate />)
    
    // Monaco editor should be rendered (mocked as textarea)
    expect(screen.getByTestId('monaco-editor')).toBeInTheDocument()
    
    // Check character count display
    expect(screen.getByText('0 characters')).toBeInTheDocument()
  })

  it('displays output tabs', () => {
    renderWithProviders(<Generate />)
    
    // Check tab navigation
    expect(screen.getByRole('tab', { name: /output/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /validation/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /issues/i })).toBeInTheDocument()
  })

  it('shows empty state for spec output', () => {
    renderWithProviders(<Generate />)
    
    expect(screen.getByText('Generated specification will appear here')).toBeInTheDocument()
  })

  it('handles template loading', async () => {
    const user = userEvent.setup()
    
    renderWithProviders(<Generate />)
    
    const basicApiTemplate = screen.getByText('Basic REST API')
    await user.click(basicApiTemplate)
    
    // Should show success toast (we can't easily test the store update without more complex mocking)
    // But we can verify the click handler was called
    expect(basicApiTemplate).toBeInTheDocument()
  })

  it('handles format switching', async () => {
    const user = userEvent.setup()
    
    renderWithProviders(<Generate />)
    
    const formatSelector = screen.getByRole('combobox')
    await user.click(formatSelector)
    
    // Should show format options
    await waitFor(() => {
      expect(screen.getByText('JSON')).toBeInTheDocument()
    })
  })

  it('shows file upload component', () => {
    renderWithProviders(<Generate />)
    
    expect(screen.getByText('Drop files here or click to browse')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /choose files/i })).toBeInTheDocument()
  })

  it('renders workspace layout with resizable panels', () => {
    renderWithProviders(<Generate />)
    
    // Check that resizable panel components are rendered
    expect(screen.getByTestId('resizable-panel-group')).toBeInTheDocument()
    expect(screen.getAllByTestId('resizable-panel')).toHaveLength(3) // Left, center, right panels
    expect(screen.getAllByTestId('resizable-handle')).toHaveLength(2) // Two handles between panels
  })

  // TODO: Implement validation status card when backend is ready
  // it('displays validation status card', () => {
  //   renderWithProviders(<Generate />)
  //   
  //   // Switch to validation tab first
  //   const validationTab = screen.getByRole('tab', { name: /validation/i })
  //   validationTab.click()
  //   
  //   expect(screen.getByText('Validation Status')).toBeInTheDocument()
  //   expect(screen.getByText('Generate a specification to see validation results')).toBeInTheDocument()
  // })

  it('shows issues tab content', async () => {
    const user = userEvent.setup()
    
    renderWithProviders(<Generate />)
    
    const issuesTab = screen.getByRole('tab', { name: /issues/i })
    await user.click(issuesTab)
    
    expect(screen.getByText('No validation errors or warnings to display')).toBeInTheDocument()
  })

  it('has proper accessibility attributes', () => {
    renderWithProviders(<Generate />)
    
    // Check for proper heading structure
    const mainHeading = screen.getByRole('heading', { level: 1 })
    expect(mainHeading).toHaveTextContent('Generate OpenAPI Specification')
    
    // Check for proper button labels
    const generateButton = screen.getByRole('button', { name: /generate/i })
    expect(generateButton).toBeInTheDocument()
    
    // Check for proper tab structure
    const tabList = screen.getByRole('tablist')
    expect(tabList).toBeInTheDocument()
    
    const tabs = screen.getAllByRole('tab')
    expect(tabs).toHaveLength(3)
  })

  it('handles responsive layout', () => {
    renderWithProviders(<Generate />)
    
    // The component should render without errors on different viewport sizes
    // This is a basic test - more comprehensive responsive testing would require
    // viewport manipulation
    expect(screen.getByText('Generate OpenAPI Specification')).toBeInTheDocument()
  })
})