"use client"

import * as React from "react"
import { cn } from "@/lib/utility"
import { ChevronRightIcon, ChevronDownIcon } from "lucide-react"

export interface TreeNode {
  id: string
  label: string
  icon?: React.ReactNode
  badge?: React.ReactNode
  children?: TreeNode[]
  disabled?: boolean
  selected?: boolean
}

export interface TreeViewProps extends Omit<React.ComponentProps<"div">, "onSelect"> {
  data: TreeNode[]
  defaultExpanded?: string[]
  expanded?: string[]
  onExpand?: (id: string, expanded: boolean) => void
  selected?: string[]
  onSelect?: (id: string) => void
  size?: "sm" | "default"
}

function TreeNodeComponent({
  node,
  level,
  size,
  expandedSet,
  onToggle,
  selected,
  onSelect,
}: {
  node: TreeNode
  level: number
  size: "sm" | "default"
  expandedSet: Set<string>
  onToggle: (id: string) => void
  selected: Set<string>
  onSelect: (id: string) => void
}) {
  const hasChildren = node.children && node.children.length > 0
  const isExpanded = expandedSet.has(node.id)
  const isSelected = selected.has(node.id)
  const indent = level * 16

  return (
    <div data-slot="tree-node">
      <div
        data-slot="tree-node-content"
        className={cn(
          "flex items-center gap-1.5 rounded-md px-2 py-1 cursor-pointer hover:bg-accent hover:text-accent-foreground",
          isSelected && "bg-accent text-accent-foreground",
          node.disabled && "opacity-50 pointer-events-none",
          size === "sm" && "text-xs py-0.5"
        )}
        style={{ paddingLeft: `${indent + 8}px` }}
        onClick={() => {
          if (hasChildren) onToggle(node.id)
          onSelect(node.id)
        }}
      >
        {hasChildren ? (
          <button
            type="button"
            className="shrink-0 size-4 flex items-center justify-center rounded hover:bg-muted"
            onClick={(e) => {
              e.stopPropagation()
              onToggle(node.id)
            }}
          >
            {isExpanded ? (
              <ChevronDownIcon className="size-3" />
            ) : (
              <ChevronRightIcon className="size-3" />
            )}
          </button>
        ) : (
          <span className="shrink-0 size-4" />
        )}
        {node.icon && <span className="shrink-0">{node.icon}</span>}
        <span className="flex-1 truncate">{node.label}</span>
        {node.badge && <span className="shrink-0">{node.badge}</span>}
      </div>
      {hasChildren && isExpanded && (
        <div data-slot="tree-node-children">
          {node.children!.map((child) => (
            <TreeNodeComponent
              key={child.id}
              node={child}
              level={level + 1}
              size={size}
              expandedSet={expandedSet}
              onToggle={onToggle}
              selected={selected}
              onSelect={onSelect}
            />
          ))}
        </div>
      )}
    </div>
  )
}

function TreeView({
  className,
  data,
  defaultExpanded = [],
  expanded: controlledExpanded,
  onExpand,
  selected: controlledSelected,
  onSelect,
  size = "default",
  ...props
}: TreeViewProps) {
  const [internalExpanded, setInternalExpanded] = React.useState<Set<string>>(
    new Set(defaultExpanded)
  )
  const [internalSelected, setInternalSelected] = React.useState<Set<string>>(new Set())

  const expandedSet = controlledExpanded ? new Set(controlledExpanded) : internalExpanded
  const selectedSet = controlledSelected ? new Set(controlledSelected) : internalSelected

  const handleToggle = React.useCallback(
    (id: string) => {
      if (controlledExpanded) {
        onExpand?.(id, !expandedSet.has(id))
      } else {
        setInternalExpanded((prev) => {
          const next = new Set(prev)
          const isExpanded = next.has(id)
          if (isExpanded) {
            next.delete(id)
          } else {
            next.add(id)
          }
          onExpand?.(id, !isExpanded)
          return next
        })
      }
    },
    [controlledExpanded, expandedSet, onExpand]
  )

  const handleSelect = React.useCallback(
    (id: string) => {
      if (controlledSelected) {
        onSelect?.(id)
      } else {
        setInternalSelected(new Set([id]))
        onSelect?.(id)
      }
    },
    [controlledSelected, onSelect]
  )

  return (
    <div data-slot="tree-view" className={cn("py-1", className)} {...props}>
      {data.map((node) => (
        <TreeNodeComponent
          key={node.id}
          node={node}
          level={0}
          size={size}
          expandedSet={expandedSet}
          onToggle={handleToggle}
          selected={selectedSet}
          onSelect={handleSelect}
        />
      ))}
    </div>
  )
}

export { TreeView }
