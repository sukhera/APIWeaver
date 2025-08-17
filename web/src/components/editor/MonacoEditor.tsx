import { useRef, useEffect } from 'react'
import Editor, { OnMount, OnChange } from '@monaco-editor/react'
import type { editor } from 'monaco-editor'
import { useTheme } from '@/hooks/useTheme'
import { cn } from '@/lib/utils'

interface MonacoEditorProps {
  value: string
  onChange?: OnChange
  language?: string
  height?: string | number
  readOnly?: boolean
  className?: string
  options?: Record<string, unknown>
}

export default function MonacoEditor({
  value,
  onChange,
  language = 'markdown',
  height = '100%',
  readOnly = false,
  className,
  options = {}
}: MonacoEditorProps) {
  const { theme } = useTheme()
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null)

  const handleEditorDidMount: OnMount = (editor, monaco) => {
    editorRef.current = editor

    // Configure Monaco themes
    monaco.editor.defineTheme('apiweaver-light', {
      base: 'vs',
      inherit: true,
      rules: [],
      colors: {
        'editor.background': '#ffffff',
        'editor.foreground': '#1e293b',
        'editor.lineHighlightBackground': '#f8fafc',
        'editor.selectionBackground': '#3b82f6',
        'editorLineNumber.foreground': '#94a3b8',
        'editor.border': '#e2e8f0'
      }
    })

    monaco.editor.defineTheme('apiweaver-dark', {
      base: 'vs-dark',
      inherit: true,
      rules: [],
      colors: {
        'editor.background': '#0f172a',
        'editor.foreground': '#f1f5f9',
        'editor.lineHighlightBackground': '#1e293b',
        'editor.selectionBackground': '#3b82f6',
        'editorLineNumber.foreground': '#475569',
        'editor.border': '#334155'
      }
    })

    // Set initial theme
    const monacoTheme = theme === 'dark' ? 'apiweaver-dark' : 'apiweaver-light'
    monaco.editor.setTheme(monacoTheme)

    // Configure markdown language
    if (language === 'markdown') {
      monaco.languages.setLanguageConfiguration('markdown', {
        wordPattern: /(-?\d*\.\d\w*)|([^`~!@#%^&*()-=+[{}\\|;:'",.<>/?]+)/g,
        brackets: [
          ['{', '}'],
          ['[', ']'],
          ['(', ')']
        ],
        autoClosingPairs: [
          { open: '{', close: '}' },
          { open: '[', close: ']' },
          { open: '(', close: ')' },
          { open: '"', close: '"' },
          { open: "'", close: "'" },
          { open: '`', close: '`' }
        ]
      })
    }
  }

  // Update theme when it changes
  useEffect(() => {
    if (editorRef.current) {
      const monacoTheme = theme === 'dark' ? 'apiweaver-dark' : 'apiweaver-light'
      editorRef.current.updateOptions({ theme: monacoTheme })
    }
  }, [theme])

  const defaultOptions = {
    minimap: { enabled: false },
    fontSize: 14,
    lineHeight: 20,
    fontFamily: 'JetBrains Mono, Consolas, monospace',
    wordWrap: 'on' as const,
    lineNumbers: 'on' as const,
    scrollBeyondLastLine: false,
    automaticLayout: true,
    tabSize: 2,
    insertSpaces: true,
    folding: true,
    foldingStrategy: 'indentation' as const,
    renderLineHighlight: 'line' as const,
    selectOnLineNumbers: true,
    matchBrackets: 'always' as const,
    readOnly,
    ...options
  }

  return (
    <div className={cn('border border-border rounded-md overflow-hidden', className)}>
      <Editor
        height={height}
        language={language}
        value={value}
        onChange={onChange}
        onMount={handleEditorDidMount}
        options={defaultOptions}
        loading={
          <div className="flex items-center justify-center h-full bg-surface">
            <div className="text-text-secondary">Loading editor...</div>
          </div>
        }
      />
    </div>
  )
}