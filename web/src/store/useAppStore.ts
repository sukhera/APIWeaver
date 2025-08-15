import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { GenerateResponse, ValidateResponse, AmendResponse, DiffResponse } from '@/types/api'

// Editor State
interface EditorState {
  markdown: string
  spec: string
  format: 'yaml' | 'json'
  activeTab: 'editor' | 'preview' | 'diff'
  isModified: boolean
  lastSaved: string | null
}

// Workspace State
interface WorkspaceState {
  showLeftPanel: boolean
  showRightPanel: boolean
  leftPanelWidth: number
  rightPanelWidth: number
  currentFile: string | null
  recentFiles: string[]
}

// Generation State
interface GenerationState {
  lastGeneration: GenerateResponse | null
  lastValidation: ValidateResponse | null
  lastAmendment: AmendResponse | null
  lastDiff: DiffResponse | null
  isGenerating: boolean
  isValidating: boolean
  isAmending: boolean
}

// Settings State
interface SettingsState {
  autoSave: boolean
  autoValidate: boolean
  showLineNumbers: boolean
  wordWrap: boolean
  fontSize: number
  tabSize: number
  theme: 'light' | 'dark' | 'system'
}

// Combined App State
interface AppState {
  editor: EditorState
  workspace: WorkspaceState
  generation: GenerationState
  settings: SettingsState
}

// Actions
interface AppActions {
  // Editor Actions
  updateMarkdown: (markdown: string) => void
  updateSpec: (spec: string) => void
  setFormat: (format: 'yaml' | 'json') => void
  setActiveTab: (tab: 'editor' | 'preview' | 'diff') => void
  markAsModified: () => void
  markAsSaved: () => void
  
  // Workspace Actions
  toggleLeftPanel: () => void
  toggleRightPanel: () => void
  setLeftPanelWidth: (width: number) => void
  setRightPanelWidth: (width: number) => void
  setCurrentFile: (file: string | null) => void
  addRecentFile: (file: string) => void
  
  // Generation Actions
  setGenerationResult: (result: GenerateResponse) => void
  setValidationResult: (result: ValidateResponse) => void
  setAmendmentResult: (result: AmendResponse) => void
  setDiffResult: (result: DiffResponse) => void
  setGenerating: (isGenerating: boolean) => void
  setValidating: (isValidating: boolean) => void
  setAmending: (isAmending: boolean) => void
  
  // Settings Actions
  updateSettings: (settings: Partial<SettingsState>) => void
  resetSettings: () => void
  
  // Utility Actions
  resetEditor: () => void
  resetWorkspace: () => void
  resetGeneration: () => void
}

const defaultEditorState: EditorState = {
  markdown: '',
  spec: '',
  format: 'yaml',
  activeTab: 'editor',
  isModified: false,
  lastSaved: null,
}

const defaultWorkspaceState: WorkspaceState = {
  showLeftPanel: true,
  showRightPanel: true,
  leftPanelWidth: 25,
  rightPanelWidth: 25,
  currentFile: null,
  recentFiles: [],
}

const defaultGenerationState: GenerationState = {
  lastGeneration: null,
  lastValidation: null,
  lastAmendment: null,
  lastDiff: null,
  isGenerating: false,
  isValidating: false,
  isAmending: false,
}

const defaultSettingsState: SettingsState = {
  autoSave: true,
  autoValidate: true,
  showLineNumbers: true,
  wordWrap: true,
  fontSize: 14,
  tabSize: 2,
  theme: 'system',
}

export const useAppStore = create<AppState & AppActions>()(
  persist(
    (set) => ({
      // Initial State
      editor: defaultEditorState,
      workspace: defaultWorkspaceState,
      generation: defaultGenerationState,
      settings: defaultSettingsState,

      // Editor Actions
      updateMarkdown: (markdown) =>
        set((state) => ({
          editor: {
            ...state.editor,
            markdown,
            isModified: true,
          },
        })),

      updateSpec: (spec) =>
        set((state) => ({
          editor: {
            ...state.editor,
            spec,
          },
        })),

      setFormat: (format) =>
        set((state) => ({
          editor: {
            ...state.editor,
            format,
          },
        })),

      setActiveTab: (activeTab) =>
        set((state) => ({
          editor: {
            ...state.editor,
            activeTab,
          },
        })),

      markAsModified: () =>
        set((state) => ({
          editor: {
            ...state.editor,
            isModified: true,
          },
        })),

      markAsSaved: () =>
        set((state) => ({
          editor: {
            ...state.editor,
            isModified: false,
            lastSaved: new Date().toISOString(),
          },
        })),

      // Workspace Actions
      toggleLeftPanel: () =>
        set((state) => ({
          workspace: {
            ...state.workspace,
            showLeftPanel: !state.workspace.showLeftPanel,
          },
        })),

      toggleRightPanel: () =>
        set((state) => ({
          workspace: {
            ...state.workspace,
            showRightPanel: !state.workspace.showRightPanel,
          },
        })),

      setLeftPanelWidth: (leftPanelWidth) =>
        set((state) => ({
          workspace: {
            ...state.workspace,
            leftPanelWidth,
          },
        })),

      setRightPanelWidth: (rightPanelWidth) =>
        set((state) => ({
          workspace: {
            ...state.workspace,
            rightPanelWidth,
          },
        })),

      setCurrentFile: (currentFile) =>
        set((state) => ({
          workspace: {
            ...state.workspace,
            currentFile,
          },
        })),

      addRecentFile: (file) =>
        set((state) => {
          const recentFiles = [
            file,
            ...state.workspace.recentFiles.filter((f) => f !== file),
          ].slice(0, 10) // Keep only last 10 files

          return {
            workspace: {
              ...state.workspace,
              recentFiles,
            },
          }
        }),

      // Generation Actions
      setGenerationResult: (lastGeneration) =>
        set((state) => ({
          generation: {
            ...state.generation,
            lastGeneration,
          },
        })),

      setValidationResult: (lastValidation) =>
        set((state) => ({
          generation: {
            ...state.generation,
            lastValidation,
          },
        })),

      setAmendmentResult: (lastAmendment) =>
        set((state) => ({
          generation: {
            ...state.generation,
            lastAmendment,
          },
        })),

      setDiffResult: (lastDiff) =>
        set((state) => ({
          generation: {
            ...state.generation,
            lastDiff,
          },
        })),

      setGenerating: (isGenerating) =>
        set((state) => ({
          generation: {
            ...state.generation,
            isGenerating,
          },
        })),

      setValidating: (isValidating) =>
        set((state) => ({
          generation: {
            ...state.generation,
            isValidating,
          },
        })),

      setAmending: (isAmending) =>
        set((state) => ({
          generation: {
            ...state.generation,
            isAmending,
          },
        })),

      // Settings Actions
      updateSettings: (newSettings) =>
        set((state) => ({
          settings: {
            ...state.settings,
            ...newSettings,
          },
        })),

      resetSettings: () =>
        set(() => ({
          settings: defaultSettingsState,
        })),

      // Utility Actions
      resetEditor: () =>
        set(() => ({
          editor: defaultEditorState,
        })),

      resetWorkspace: () =>
        set(() => ({
          workspace: defaultWorkspaceState,
        })),

      resetGeneration: () =>
        set(() => ({
          generation: defaultGenerationState,
        })),
    }),
    {
      name: 'apiweaver-app-store',
      partialize: (state) => ({
        workspace: state.workspace,
        settings: state.settings,
        // Don't persist editor content or generation results
      }),
    }
  )
)

// Convenient selector hooks
export const useEditorState = () => useAppStore((state) => state.editor)
export const useWorkspaceState = () => useAppStore((state) => state.workspace)
export const useGenerationState = () => useAppStore((state) => state.generation)
export const useSettingsState = () => useAppStore((state) => state.settings)