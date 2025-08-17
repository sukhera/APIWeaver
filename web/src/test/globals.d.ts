/// <reference types="@testing-library/jest-dom" />

declare global {
  const ResizeObserver: typeof ResizeObserver
  const IntersectionObserver: typeof IntersectionObserver
  const fetch: typeof fetch
  const FileReader: typeof FileReader
  
  namespace globalThis {
    const ResizeObserver: typeof ResizeObserver
    const IntersectionObserver: typeof IntersectionObserver
    const fetch: typeof fetch
    const FileReader: typeof FileReader
  }
}

// Extend Element interface for test mocks
declare global {
  interface Element {
    hasPointerCapture: (pointerId: number) => boolean
    setPointerCapture: (pointerId: number) => void
    releasePointerCapture: (pointerId: number) => void
    scrollIntoView: () => void
  }
}

export {}