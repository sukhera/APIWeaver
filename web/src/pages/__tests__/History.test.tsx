import React from 'react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { ThemeProvider } from '@/components/theme-provider'
import History from '../History'

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

// Mock localStorage
const mockLocalStorage = {
  store: new Map<string, string>(),
  getItem: vi.fn((key: string) => mockLocalStorage.store.get(key) || null),
  setItem: vi.fn((key: string, value: string) => mockLocalStorage.store.set(key, value)),
  removeItem: vi.fn((key: string) => mockLocalStorage.store.delete(key)),
  clear: vi.fn(() => mockLocalStorage.store.clear()),
}

Object.defineProperty(window, 'localStorage', {
  value: mockLocalStorage,
})

describe('History Page - Simple Tests', () => {
  beforeEach(() => {
    mockLocalStorage.clear()
    vi.clearAllMocks()
  })

  it('renders without crashing', () => {
    expect(() => {
      render(
        <TestWrapper>
          <History />
        </TestWrapper>
      )
    }).not.toThrow()
  })

  it('displays the page title', () => {
    render(
      <TestWrapper>
        <History />
      </TestWrapper>
    )

    expect(screen.getByText('Conversion History')).toBeInTheDocument()
  })

  it('displays empty state when no history exists', () => {
    mockLocalStorage.getItem.mockReturnValue('[]')
    
    render(
      <TestWrapper>
        <History />
      </TestWrapper>
    )

    expect(screen.getByText('No conversion history')).toBeInTheDocument()
  })

  it('displays header controls', () => {
    render(
      <TestWrapper>
        <History />
      </TestWrapper>
    )

    expect(screen.getByRole('button', { name: /refresh/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /clear all/i })).toBeInTheDocument()
  })

  it('displays search input', () => {
    render(
      <TestWrapper>
        <History />
      </TestWrapper>
    )

    expect(screen.getByPlaceholderText('Search conversions...')).toBeInTheDocument()
  })

  it('handles localStorage errors gracefully', () => {
    // Mock localStorage to throw error only for specific key, not theme
    mockLocalStorage.getItem.mockImplementation((key) => {
      if (key === 'apiweaver-history') {
        throw new Error('localStorage error')
      }
      if (key === 'test-theme') {
        return 'light' // Return valid theme
      }
      return null
    })

    // The component should handle errors internally and not crash
    render(
      <TestWrapper>
        <History />
      </TestWrapper>
    )

    // Should still render the basic page structure
    expect(screen.getByText('Conversion History')).toBeInTheDocument()
  })

  it('creates mock data when localStorage is empty', () => {
    mockLocalStorage.getItem.mockReturnValue(null)
    
    render(
      <TestWrapper>
        <History />
      </TestWrapper>
    )

    // Should have called setItem to save mock data
    expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
      'apiweaver-history',
      expect.stringContaining('Task API Generation')
    )
  })

  it('loads existing history from localStorage', () => {
    const mockHistory = [
      {
        id: '1',
        timestamp: '2024-01-15T10:30:00Z',
        inputType: 'markdown',
        inputContent: 'Test API content',
        outputContent: 'openapi: 3.1.0\ninfo:\n  title: Test API',
        outputFormat: 'yaml',
        success: true,
        operation: 'generate',
        title: 'Test_Generation',
        inputSize: 100,
        outputSize: 200,
        processingTime: 300
      }
    ]
    
    mockLocalStorage.getItem.mockImplementation((key) => {
      if (key === 'apiweaver-history') {
        return JSON.stringify(mockHistory)
      }
      if (key === 'test-theme') {
        return 'light'
      }
      return null
    })
    
    render(
      <TestWrapper>
        <History />
      </TestWrapper>
    )

    expect(screen.getByText('Test_Generation')).toBeInTheDocument()
    expect(screen.getByText('Success')).toBeInTheDocument()
  })

  it('displays operation-specific content', () => {
    const mockHistory = [
      {
        id: '1',
        timestamp: '2024-01-15T10:30:00Z',
        inputType: 'markdown',
        inputContent: 'Test_API_content',
        outputContent: 'openapi: 3.1.0\ninfo:\n  title: Test',
        outputFormat: 'yaml',
        success: true,
        operation: 'generate',
        title: 'Test_Generation',
        inputSize: 100,
        outputSize: 200,
        processingTime: 300
      },
      {
        id: '2',
        timestamp: '2024-01-14T10:30:00Z',
        inputType: 'openapi',
        inputContent: 'invalid_spec_content',
        outputFormat: 'yaml',
        success: false,
        errors: ['Validation_error'],
        operation: 'validate',
        title: 'Test_Validation',
        inputSize: 50,
        processingTime: 150
      }
    ]
    
    mockLocalStorage.getItem.mockImplementation((key) => {
      if (key === 'apiweaver-history') {
        return JSON.stringify(mockHistory)
      }
      if (key === 'test-theme') {
        return 'light'
      }
      return null
    })
    
    render(
      <TestWrapper>
        <History />
      </TestWrapper>
    )

    expect(screen.getByText('Test_Generation')).toBeInTheDocument()
    expect(screen.getByText('Test_Validation')).toBeInTheDocument()
    expect(screen.getByText('Success')).toBeInTheDocument()
    expect(screen.getByText('Error')).toBeInTheDocument()
  })
})
