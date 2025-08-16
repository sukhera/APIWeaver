import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { ThemeProvider } from '@/components/theme-provider'
import { Toaster } from 'sonner'

// Import pages
import Home from '@/pages/Home'
import Generate from '@/pages/Generate'
import Validate from '@/pages/Validate'
import Amend from '@/pages/Amend'
import History from '@/pages/History'

// Import layout
import AppLayout from '@/components/layout/AppLayout'

// Import styles
import '@/styles/globals.css'

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, error: unknown) => {
        // Don't retry on 4xx errors
        const apiError = error as { status?: number }
        if (apiError?.status && apiError.status >= 400 && apiError.status < 500) {
          return false
        }
        return failureCount < 3
      },
      staleTime: 5 * 60 * 1000, // 5 minutes
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: (failureCount, error: unknown) => {
        // Don't retry mutations on client errors
        const apiError = error as { status?: number }
        if (apiError?.status && apiError.status >= 400 && apiError.status < 500) {
          return false
        }
        return failureCount < 2
      },
    },
  },
})

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="system" storageKey="apiweaver-ui-theme">
        <Router>
          <div className="min-h-screen bg-background font-sans antialiased">
            <AppLayout>
              <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/generate" element={<Generate />} />
                <Route path="/validate" element={<Validate />} />
                <Route path="/amend" element={<Amend />} />
                <Route path="/history" element={<History />} />
              </Routes>
            </AppLayout>
            <Toaster richColors position="top-right" />
          </div>
        </Router>
      </ThemeProvider>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  )
}

export default App