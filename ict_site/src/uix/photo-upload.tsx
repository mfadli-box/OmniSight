"use client"

import * as React from "react"
import { cn } from "@/lib/utility"
import { Button } from "./button"
import { UploadIcon, CameraIcon, XIcon, ImageIcon } from "lucide-react"

export interface PhotoUploadFile {
  id: string
  file: File
  preview: string
  caption?: string
}

export interface PhotoUploadProps extends React.ComponentProps<"div"> {
  files?: PhotoUploadFile[]
  onFilesChange?: (files: PhotoUploadFile[]) => void
  maxFiles?: number
  maxSize?: number
  disabled?: boolean
  capture?: "user" | "environment"
}

function PhotoUpload({
  className,
  files: controlledFiles,
  onFilesChange,
  maxFiles = 10,
  maxSize = 10 * 1024 * 1024,
  disabled = false,
  capture,
  ...props
}: PhotoUploadProps) {
  const [internalFiles, setInternalFiles] = React.useState<PhotoUploadFile[]>([])
  const fileInputRef = React.useRef<HTMLInputElement>(null)
  const cameraInputRef = React.useRef<HTMLInputElement>(null)

  const files = controlledFiles ?? internalFiles

  const addFiles = React.useCallback(
    (newFiles: FileList | File[]) => {
      const fileArray = Array.from(newFiles).filter((file) => {
        if (!file.type.startsWith("image/")) return false
        if (file.size > maxSize) return false
        if (files.length + 1 > maxFiles) return false
        return true
      })

      const newEntries: PhotoUploadFile[] = fileArray.map((file) => ({
        id: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
        file,
        preview: URL.createObjectURL(file),
      }))

      const updated = [...files, ...newEntries]
      setInternalFiles(updated)
      onFilesChange?.(updated)
    },
    [files, maxFiles, maxSize, onFilesChange]
  )

  const removeFile = React.useCallback(
    (id: string) => {
      const removed = files.find((f) => f.id === id)
      if (removed) URL.revokeObjectURL(removed.preview)
      const updated = files.filter((f) => f.id !== id)
      setInternalFiles(updated)
      onFilesChange?.(updated)
    },
    [files, onFilesChange]
  )

  const updateCaption = React.useCallback(
    (id: string, caption: string) => {
      const updated = files.map((f) => (f.id === id ? { ...f, caption } : f))
      setInternalFiles(updated)
      onFilesChange?.(updated)
    },
    [files, onFilesChange]
  )

  return (
    <div data-slot="photo-upload" className={cn("space-y-3", className)} {...props}>
      <div className="flex gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => fileInputRef.current?.click()}
          disabled={disabled || files.length >= maxFiles}
        >
          <UploadIcon className="size-4 mr-1.5" />
          Gallery
        </Button>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => cameraInputRef.current?.click()}
          disabled={disabled || files.length >= maxFiles}
        >
          <CameraIcon className="size-4 mr-1.5" />
          Camera
        </Button>
        <span className="text-xs text-muted-foreground self-center">
          {files.length}/{maxFiles}
        </span>
      </div>

      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        multiple={maxFiles > 1}
        className="hidden"
        onChange={(e) => {
          if (e.target.files) addFiles(e.target.files)
          e.target.value = ""
        }}
      />
      <input
        ref={cameraInputRef}
        type="file"
        accept="image/*"
        capture={capture}
        multiple={maxFiles > 1}
        className="hidden"
        onChange={(e) => {
          if (e.target.files) addFiles(e.target.files)
          e.target.value = ""
        }}
      />

      {files.length > 0 && (
        <div data-slot="photo-upload-grid" className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
          {files.map((file) => (
            <div
              key={file.id}
              data-slot="photo-upload-item"
              className="group relative overflow-hidden rounded-md border bg-muted"
            >
              <img
                src={file.preview}
                alt={file.caption || "Upload preview"}
                className="aspect-square w-full object-cover"
              />
              <button
                type="button"
                onClick={() => removeFile(file.id)}
                className="absolute top-1 right-1 size-6 rounded-full bg-black/60 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <XIcon className="size-3" />
              </button>
              <input
                type="text"
                placeholder="Caption..."
                value={file.caption || ""}
                onChange={(e) => updateCaption(file.id, e.target.value)}
                className="w-full border-t px-2 py-1 text-xs bg-background"
              />
            </div>
          ))}
        </div>
      )}

      {files.length === 0 && (
        <div className="flex flex-col items-center justify-center rounded-lg border-2 border-dashed p-8 text-center text-muted-foreground">
          <ImageIcon className="size-8 mb-2" />
          <p className="text-sm">No photos selected</p>
          <p className="text-xs">Click Gallery or Camera to add photos</p>
        </div>
      )}
    </div>
  )
}

export { PhotoUpload }
