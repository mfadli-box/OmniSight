"use client"

import * as React from "react"
import { cn } from "@/lib/utility"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/uix/select"

export interface LocationOption {
  id: string
  code: string
  name: string
  parentId?: string | null
  children?: LocationOption[]
}

export interface LocationPickerProps {
  value?: string
  onValueChange?: (value: string | null) => void
  locations: LocationOption[]
  placeholder?: string
  disabled?: boolean
  className?: string
}

function flattenTree(nodes: LocationOption[], level = 0): { option: LocationOption; level: number }[] {
  const result: { option: LocationOption; level: number }[] = []
  for (const node of nodes) {
    result.push({ option: node, level })
    if (node.children?.length) {
      result.push(...flattenTree(node.children, level + 1))
    }
  }
  return result
}

function LocationPicker({
  value,
  onValueChange,
  locations,
  placeholder = "Select location",
  disabled,
  className,
}: LocationPickerProps) {
  const flat = React.useMemo(() => flattenTree(locations), [locations])

  return (
    <Select
      value={value}
      onValueChange={(v) => onValueChange?.(v)}
      disabled={disabled}
    >
      <SelectTrigger className={cn("w-full", className)}>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {flat.map(({ option, level }) => (
          <SelectItem key={option.id} value={option.id}>
            <span style={{ paddingLeft: `${level * 16}px` }} className="inline-flex items-center gap-1.5">
              <span className="text-muted-foreground text-xs font-mono">{option.code}</span>
              <span>{option.name}</span>
            </span>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}

export { LocationPicker }
