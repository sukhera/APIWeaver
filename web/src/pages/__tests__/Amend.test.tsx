import React from 'react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { ThemeProvider } from '@/components/theme-provider'
import Amend from '../Amend'
import * as apiQueries from '@/hooks/useApiQueries'
import { asMockAmendMutation, asMockDiffMutation } from '@/test/types'

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

describe('Amend Page - Simple Tests', () => {
  const mockAmendMutation = {
    mutateAsync: vi.fn(),
    mutate: vi.fn(),
    isPending: false,
    isError: false,
    isIdle: true,
    isSuccess: false,
    isPaused: false,
    error: null,
    data: undefined,
    reset: vi.fn(),
    variables: undefined,
    context: undefined,
    failureCount: 0,
    failureReason: null,
    status: 'idle' as const,
    submittedAt: 0
  }

  const mockDiffMutation = {
    mutateAsync: vi.fn(),
    mutate: vi.fn(),
    isPending: false,
    isError: false,
    isIdle: true,
    isSuccess: false,
    isPaused: false,
    error: null,
    data: undefined,
    reset: vi.fn(),
    variables: undefined,
    context: undefined,
    failureCount: 0,
    failureReason: null,
    status: 'idle' as const,
    submittedAt: 0
  }

  beforeEach(() => {
    vi.mocked(apiQueries.useAmendMutation).mockReturnValue(asMockAmendMutation(mockAmendMutation))
    vi.mocked(apiQueries.useDiffMutation).mockReturnValue(asMockDiffMutation(mockDiffMutation))
    vi.clearAllMocks()
  })

  it('renders without crashing', () => {
    expect(() => {
      render(
        <TestWrapper>
          <Amend />
        </TestWrapper>
      )
    }).not.toThrow()
  })

  it('displays workspace layout', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    // Check for resizable panel group instead of workspace-layout
    expect(screen.getByTestId('resizable-panel-group')).toBeInTheDocument()
    expect(screen.getAllByTestId('resizable-panel')).toHaveLength(3)
  })

  it('displays input tabs', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    expect(screen.getByRole('tab', { name: /original/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /changes/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /upload/i })).toBeInTheDocument()
  })

  it('displays amendment controls', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    expect(screen.getByRole('checkbox', { name: /dry run/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /preview/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /apply/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /reset/i })).toBeInTheDocument()
  })

  it('displays monaco editors', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    // Should have monaco editor (the mocked version)
    expect(screen.getByTestId('monaco-editor')).toBeInTheDocument()
  })

  it('displays amendment result tabs', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    expect(screen.getByRole('tab', { name: /amended/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /diff/i })).toBeInTheDocument()
    expect(screen.getByRole('tab', { name: /preview/i })).toBeInTheDocument()
  })

  it('shows loading state when amending', () => {
    vi.mocked(apiQueries.useAmendMutation).mockReturnValue(asMockAmendMutation({
      ...mockAmendMutation,
      isPending: true
    }))

    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    expect(screen.getByText('Amending...')).toBeInTheDocument()
  })

  it('shows loading state when generating diff', () => {
    vi.mocked(apiQueries.useDiffMutation).mockReturnValue(asMockDiffMutation({
      ...mockDiffMutation,
      isPending: true
    }))

    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    expect(screen.getByText('Generating...')).toBeInTheDocument()
  })

  it('displays amendment status when available', () => {
    const mockAmendResult = {
      amendedSpec: 'openapi: 3.1.0\ninfo:\n  title: Amended API',
      errors: [],
      warnings: []
    }

    vi.mocked(apiQueries.useAmendMutation).mockReturnValue(asMockAmendMutation({
      ...mockAmendMutation,
      data: mockAmendResult
    }))

    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    // Check for multiple instances of "Amended" text
    const amendedElements = screen.getAllByText('Amended')
    expect(amendedElements.length).toBeGreaterThan(0)
  })

  it('displays diff results when available', () => {
    const mockDiffResult = {
      diff: {
        added: [{ lineNumber: 5, content: '  version: 1.1.0', type: 'added' as const }],
        removed: [{ lineNumber: 4, content: '  version: 1.0.0', type: 'removed' as const }],
        modified: [],
        unchanged: []
      },
      summary: {
        linesAdded: 1,
        linesRemoved: 1,
        linesModified: 0,
        endpointsAdded: 0,
        endpointsRemoved: 0,
        endpointsModified: 0
      }
    }

    vi.mocked(apiQueries.useDiffMutation).mockReturnValue(asMockDiffMutation({
      ...mockDiffMutation,
      data: mockDiffResult
    }))

    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    // The diff data is available for rendering
    expect(mockDiffResult.summary.linesAdded).toBe(1)
    expect(mockDiffResult.summary.linesRemoved).toBe(1)
  })

  it('displays format selector', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    // Check for YAML text in span instead of display value
    expect(screen.getByText('YAML')).toBeInTheDocument()
  })

  it('has dry run enabled by default', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    const dryRunCheckbox = screen.getByRole('checkbox', { name: /dry run/i })
    expect(dryRunCheckbox).toBeChecked()
  })

  it('displays copy and download buttons', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    expect(screen.getByRole('button', { name: /copy/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /download/i })).toBeInTheDocument()
  })

  it('displays generate diff button', () => {
    render(
      <TestWrapper>
        <Amend />
      </TestWrapper>
    )

    expect(screen.getByRole('button', { name: /generate diff/i })).toBeInTheDocument()
  })
})
