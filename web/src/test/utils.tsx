import { ReactElement } from 'react'
import { render, RenderOptions } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { ThemeProvider } from '@/components/theme-provider'

// Create a test query client
const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        staleTime: 0,
        gcTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  })

interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  queryClient?: QueryClient
  initialRoute?: string
}

// Custom render function with providers
export function renderWithProviders(
  ui: ReactElement,
  {
    queryClient = createTestQueryClient(),
    initialRoute = '/',
    ...renderOptions
  }: CustomRenderOptions = {}
) {
  // Set initial route if specified
  if (initialRoute !== '/') {
    window.history.pushState({}, 'Test page', initialRoute)
  }

  function Wrapper({ children }: { children: React.ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <ThemeProvider defaultTheme="light" storageKey="test-theme">
          <BrowserRouter>{children}</BrowserRouter>
        </ThemeProvider>
      </QueryClientProvider>
    )
  }

  return {
    ...render(ui, { wrapper: Wrapper, ...renderOptions }),
    queryClient,
  }
}

// Helper to create a file for testing
export function createTestFile(
  name: string = 'test.md',
  content: string = 'Test content',
  type: string = 'text/markdown'
): File {
  return new File([content], name, { type })
}

// Mock API responses
export const mockApiResponses = {
  generate: {
    spec: 'openapi: 3.1.0\ninfo:\n  title: Test API\n  version: 1.0.0',
    format: 'yaml' as const,
    errors: [],
    warnings: [],
  },
  validate: {
    valid: true,
    errors: [],
    warnings: [],
    summary: {
      totalErrors: 0,
      totalWarnings: 0,
      endpoints: 1,
      schemas: 0,
      parameters: 0,
    },
  },
  health: {
    status: 'healthy' as const,
    version: '1.0.0',
    uptime: 3600,
    timestamp: new Date().toISOString(),
  },
}

// Helper to wait for async operations
export const waitFor = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

// Re-export testing utilities
// eslint-disable-next-line react-refresh/only-export-components
export * from '@testing-library/react'
export { userEvent } from '@testing-library/user-event'