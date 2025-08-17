import '@testing-library/jest-dom'
import { beforeAll, vi } from 'vitest'
import React from 'react'

// Mock ResizeObserver
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Mock IntersectionObserver
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock URL.createObjectURL
global.URL.createObjectURL = vi.fn()
global.URL.revokeObjectURL = vi.fn()

// Mock hasPointerCapture for Radix UI compatibility
Element.prototype.hasPointerCapture = vi.fn().mockReturnValue(false)
Element.prototype.setPointerCapture = vi.fn()
Element.prototype.releasePointerCapture = vi.fn()
Element.prototype.scrollIntoView = vi.fn()

// Mock fetch globally
global.fetch = vi.fn()

// Mock localStorage
const localStorageMock = {
  store: new Map<string, string>(),
  getItem: (key: string) => localStorageMock.store.get(key) || null,
  setItem: (key: string, value: string) => localStorageMock.store.set(key, value),
  removeItem: (key: string) => localStorageMock.store.delete(key),
  clear: () => localStorageMock.store.clear(),
}

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
})

// Mock clipboard API
Object.defineProperty(navigator, 'clipboard', {
  value: {
    writeText: vi.fn().mockResolvedValue(undefined),
    readText: vi.fn().mockResolvedValue(''),
  },
  configurable: true,
})

// Mock FileReader
global.FileReader = class FileReader {
  readAsText = vi.fn()
  readAsDataURL = vi.fn()
  addEventListener = vi.fn()
  removeEventListener = vi.fn()
  dispatchEvent = vi.fn()
  abort = vi.fn()
  result: string | ArrayBuffer | null = null
  error: DOMException | null = null
  readyState: number = 0
  static readonly EMPTY = 0
  static readonly LOADING = 1
  static readonly DONE = 2
  
  constructor() {
    // Simulate successful file read
    setTimeout(() => {
      this.readyState = 2 // DONE
      if (this.onload) {
        this.onload({ target: this } as ProgressEvent<FileReader>)
      }
    }, 0)
  }
  
  onload: ((this: FileReader, ev: ProgressEvent<FileReader>) => any) | null = null
  onerror: ((this: FileReader, ev: ProgressEvent<FileReader>) => any) | null = null
} as unknown as typeof FileReader

// Mock Monaco Editor
vi.mock('@monaco-editor/react', () => ({
  default: vi.fn(({ onChange, value }) => {
    return React.createElement('textarea', {
      'data-testid': 'monaco-editor',
      value: value,
      onChange: (e) => onChange?.((e.target as HTMLTextAreaElement).value)
    })
  }),
}))

// Mock react-resizable-panels
vi.mock('react-resizable-panels', () => ({
  PanelGroup: ({ children, ...props }: Record<string, unknown>) => 
    React.createElement('div', { 'data-testid': 'resizable-panel-group', ...props }, children as React.ReactNode),
  Panel: ({ children, ...props }: Record<string, unknown>) => 
    React.createElement('div', { 'data-testid': 'resizable-panel', ...props }, children as React.ReactNode),
  PanelResizeHandle: (props: Record<string, unknown>) => 
    React.createElement('div', { 'data-testid': 'resizable-handle', ...props }),
}))

beforeAll(() => {
  // Setup global test environment
  global.console = {
    ...console,
    // Suppress console.log in tests unless needed
    log: vi.fn(),
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }

  // Ensure document.body exists
  if (!document.body) {
    document.body = document.createElement('body')
    document.documentElement.appendChild(document.body)
  }
})