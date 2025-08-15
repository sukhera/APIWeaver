import { describe, it, expect, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useAppStore } from '../useAppStore'
import type { GenerateResponse, ValidateResponse } from '@/types/api'

// Mock localStorage
const mockLocalStorage = {
  store: new Map<string, string>(),
  getItem: (key: string) => mockLocalStorage.store.get(key) || null,
  setItem: (key: string, value: string) => mockLocalStorage.store.set(key, value),
  removeItem: (key: string) => mockLocalStorage.store.delete(key),
  clear: () => mockLocalStorage.store.clear(),
}

Object.defineProperty(window, 'localStorage', {
  value: mockLocalStorage,
})

describe('useAppStore', () => {
  beforeEach(() => {
    // Clear localStorage and reset store
    mockLocalStorage.clear()
    useAppStore.persist.clearStorage()
  })

  describe('Editor State', () => {
    it('initializes with default editor state', () => {
      const { result } = renderHook(() => useAppStore())
      
      expect(result.current.editor).toEqual({
        markdown: '',
        spec: '',
        format: 'yaml',
        activeTab: 'editor',
        isModified: false,
        lastSaved: null,
      })
    })

    it('updates markdown content', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.updateMarkdown('# Test API')
      })
      
      expect(result.current.editor.markdown).toBe('# Test API')
      expect(result.current.editor.isModified).toBe(true)
    })

    it('updates spec content', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.updateSpec('openapi: 3.1.0\ninfo:\n  title: Test API')
      })
      
      expect(result.current.editor.spec).toBe('openapi: 3.1.0\ninfo:\n  title: Test API')
    })

    it('changes format', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.setFormat('json')
      })
      
      expect(result.current.editor.format).toBe('json')
    })

    it('changes active tab', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.setActiveTab('preview')
      })
      
      expect(result.current.editor.activeTab).toBe('preview')
    })

    it('marks as modified and saved', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.markAsModified()
      })
      
      expect(result.current.editor.isModified).toBe(true)
      
      act(() => {
        result.current.markAsSaved()
      })
      
      expect(result.current.editor.isModified).toBe(false)
      expect(result.current.editor.lastSaved).toBeTruthy()
    })
  })

  describe('Workspace State', () => {
    it('initializes with default workspace state', () => {
      const { result } = renderHook(() => useAppStore())
      
      expect(result.current.workspace).toEqual({
        showLeftPanel: true,
        showRightPanel: true,
        leftPanelWidth: 25,
        rightPanelWidth: 25,
        currentFile: null,
        recentFiles: [],
      })
    })

    it('toggles panels', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.toggleLeftPanel()
      })
      
      expect(result.current.workspace.showLeftPanel).toBe(false)
      
      act(() => {
        result.current.toggleRightPanel()
      })
      
      expect(result.current.workspace.showRightPanel).toBe(false)
    })

    it('sets panel widths', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.setLeftPanelWidth(30)
        result.current.setRightPanelWidth(35)
      })
      
      expect(result.current.workspace.leftPanelWidth).toBe(30)
      expect(result.current.workspace.rightPanelWidth).toBe(35)
    })

    it('manages current file', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.setCurrentFile('test.md')
      })
      
      expect(result.current.workspace.currentFile).toBe('test.md')
    })

    it('manages recent files', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.addRecentFile('file1.md')
        result.current.addRecentFile('file2.md')
        result.current.addRecentFile('file1.md') // Should move to top
      })
      
      expect(result.current.workspace.recentFiles).toEqual(['file1.md', 'file2.md'])
    })

    it('limits recent files to 10', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        // Add 12 files
        for (let i = 1; i <= 12; i++) {
          result.current.addRecentFile(`file${i}.md`)
        }
      })
      
      expect(result.current.workspace.recentFiles).toHaveLength(10)
      expect(result.current.workspace.recentFiles[0]).toBe('file12.md') // Most recent first
    })
  })

  describe('Generation State', () => {
    it('initializes with default generation state', () => {
      const { result } = renderHook(() => useAppStore())
      
      expect(result.current.generation).toEqual({
        lastGeneration: null,
        lastValidation: null,
        lastAmendment: null,
        lastDiff: null,
        isGenerating: false,
        isValidating: false,
        isAmending: false,
      })
    })

    it('manages generation results', () => {
      const { result } = renderHook(() => useAppStore())
      
      const mockGeneration: GenerateResponse = {
        spec: 'openapi: 3.1.0',
        format: 'yaml',
        errors: [],
        warnings: [],
      }
      
      const mockValidation: ValidateResponse = {
        valid: true,
        errors: [],
        warnings: [],
      }
      
      act(() => {
        result.current.setGenerationResult(mockGeneration)
        result.current.setValidationResult(mockValidation)
      })
      
      expect(result.current.generation.lastGeneration).toEqual(mockGeneration)
      expect(result.current.generation.lastValidation).toEqual(mockValidation)
    })

    it('manages loading states', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.setGenerating(true)
        result.current.setValidating(true)
        result.current.setAmending(true)
      })
      
      expect(result.current.generation.isGenerating).toBe(true)
      expect(result.current.generation.isValidating).toBe(true)
      expect(result.current.generation.isAmending).toBe(true)
    })
  })

  describe('Settings State', () => {
    it('initializes with default settings', () => {
      const { result } = renderHook(() => useAppStore())
      
      expect(result.current.settings).toEqual({
        autoSave: true,
        autoValidate: true,
        showLineNumbers: true,
        wordWrap: true,
        fontSize: 14,
        tabSize: 2,
        theme: 'system',
      })
    })

    it('updates settings', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.updateSettings({
          fontSize: 16,
          tabSize: 4,
          autoSave: false,
        })
      })
      
      expect(result.current.settings.fontSize).toBe(16)
      expect(result.current.settings.tabSize).toBe(4)
      expect(result.current.settings.autoSave).toBe(false)
      // Other settings should remain unchanged
      expect(result.current.settings.showLineNumbers).toBe(true)
    })

    it('resets settings', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.updateSettings({ fontSize: 20 })
        result.current.resetSettings()
      })
      
      expect(result.current.settings.fontSize).toBe(14) // Back to default
    })
  })

  describe('Reset Functions', () => {
    it('resets editor state', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.updateMarkdown('# Test')
        result.current.setFormat('json')
        result.current.markAsModified()
        result.current.resetEditor()
      })
      
      expect(result.current.editor).toEqual({
        markdown: '',
        spec: '',
        format: 'yaml',
        activeTab: 'editor',
        isModified: false,
        lastSaved: null,
      })
    })

    it('resets workspace state', () => {
      const { result } = renderHook(() => useAppStore())
      
      act(() => {
        result.current.toggleLeftPanel()
        result.current.setLeftPanelWidth(50)
        result.current.addRecentFile('test.md')
        result.current.resetWorkspace()
      })
      
      expect(result.current.workspace).toEqual({
        showLeftPanel: true,
        showRightPanel: true,
        leftPanelWidth: 25,
        rightPanelWidth: 25,
        currentFile: null,
        recentFiles: [],
      })
    })

    it('resets generation state', () => {
      const { result } = renderHook(() => useAppStore())
      
      const mockGeneration: GenerateResponse = {
        spec: 'test',
        format: 'yaml',
        errors: [],
        warnings: [],
      }
      
      act(() => {
        result.current.setGenerationResult(mockGeneration)
        result.current.setGenerating(true)
        result.current.resetGeneration()
      })
      
      expect(result.current.generation).toEqual({
        lastGeneration: null,
        lastValidation: null,
        lastAmendment: null,
        lastDiff: null,
        isGenerating: false,
        isValidating: false,
        isAmending: false,
      })
    })
  })

  describe('Selector Hooks', () => {
    it('provides individual state selectors', () => {
      const { result: editorResult } = renderHook(() => useAppStore((state) => state.editor))
      const { result: workspaceResult } = renderHook(() => useAppStore((state) => state.workspace))
      const { result: generationResult } = renderHook(() => useAppStore((state) => state.generation))
      const { result: settingsResult } = renderHook(() => useAppStore((state) => state.settings))
      
      expect(editorResult.current).toBeDefined()
      expect(workspaceResult.current).toBeDefined()
      expect(generationResult.current).toBeDefined()
      expect(settingsResult.current).toBeDefined()
    })
  })
})