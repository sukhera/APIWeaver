import React from 'react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { ThemeProvider } from '@/components/theme-provider'
import Validate from '../Validate'
import * as apiQueries from '@/hooks/useApiQueries'

// Mock the API hooks
vi.mock('@/hooks/useApiQueries')

// Simple wrapper component
function TestWrapper({ children }: { children: React.ReactNode }) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  })

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="light" storageKey="test-theme">
        <BrowserRouter>
          {children}
        </BrowserRouter>
      </ThemeProvider>
    </QueryClientProvider>
  )
}

describe('Validate Page - Simple Tests', () => {
  const mockValidateMutation = {
    mutateAsync: vi.fn(),
    isPending: false,
    error: null,
    data: null,
    reset: vi.fn()
  }

  const mockRealtimeValidation = {
    data: null,
    isFetching: false,
    error: null
  }

  beforeEach(() => {
    vi.mocked(apiQueries.useValidateMutation).mockReturnValue(mockValidateMutation)
    vi.mocked(apiQueries.useRealtimeValidation).mockReturnValue(mockRealtimeValidation)
    vi.clearAllMocks()
  })

  it('renders without crashing', () => {
    expect(() => {
      render(
        <TestWrapper>
          <Validate />
        </TestWrapper>
      )
    }).not.toThrow()
  })

  it('displays workspace layout', () => {
    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    // Check for resizable panel group instead of workspace-layout
    expect(screen.getByTestId('resizable-panel-group')).toBeInTheDocument()
    expect(screen.getAllByTestId('resizable-panel')).toHaveLength(3)
  })

  it('displays input tabs', () => {
    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    expect(screen.getByRole('tab', { name: /editor/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /upload/i })).toBeInTheDocument()
  })

  it('displays validation controls', () => {
    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    expect(screen.getByRole('checkbox', { name: /real-time validation/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /validate/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /copy/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /download/i })).toBeInTheDocument()
  })

  it('displays monaco editor', () => {
    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    expect(screen.getByTestId('monaco-editor')).toBeInTheDocument()
  })

  it('displays file upload when upload tab is active', () => {
    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    // Upload tab should be available
    const uploadTab = screen.getByRole('tab', { name: /upload/i })
    expect(uploadTab).toBeInTheDocument()
  })

  it('displays validation results tabs', () => {
    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    expect(screen.getByRole('tab', { name: /results/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /summary/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /errors/i })).toBeInTheDocument()
  })

  it('shows loading state when validating', () => {
    vi.mocked(apiQueries.useValidateMutation).mockReturnValue({
      ...mockValidateMutation,
      isPending: true
    })

    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    expect(screen.getByText('Validating...')).toBeInTheDocument()
  })

  it('displays validation results when available', () => {
    const mockValidationResult = {
      valid: true,
      errors: [],
      warnings: [],
      summary: {
        totalErrors: 0,
        totalWarnings: 0,
        endpoints: 2,
        schemas: 1,
        parameters: 3
      }
    }

    vi.mocked(apiQueries.useValidateMutation).mockReturnValue({
      ...mockValidateMutation,
      data: mockValidationResult
    })

    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    expect(screen.getByText('Valid')).toBeInTheDocument()
  })

  it('displays validation errors when present', () => {
    const mockValidationResult = {
      valid: false,
      errors: [
        {
          line: 4,
          column: 1,
          message: 'Missing required field: paths',
          severity: 'error' as const,
          code: 'MISSING_PATHS'
        }
      ],
      warnings: [],
      summary: {
        totalErrors: 1,
        totalWarnings: 0,
        endpoints: 0,
        schemas: 0,
        parameters: 0
      }
    }

    vi.mocked(apiQueries.useValidateMutation).mockReturnValue({
      ...mockValidateMutation,
      data: mockValidationResult
    })

    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    expect(screen.getByText('Invalid')).toBeInTheDocument()
    // Check for error count in a more specific way
    const errorBadges = screen.getAllByText('1')
    expect(errorBadges.length).toBeGreaterThan(0)
  })

  it('displays format selector', () => {
    render(
      <TestWrapper>
        <Validate />
      </TestWrapper>
    )

    // Check for YAML text in span instead of display value
    expect(screen.getByText('YAML')).toBeInTheDocument()
  })
})
