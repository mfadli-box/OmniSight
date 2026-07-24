"use client"

import * as React from "react"
import { cn } from "@/lib/utility"
import { ChevronRightIcon, ChevronDownIcon, MapPinIcon } from "lucide-react"
import { LocationBadge, type locationTypeVariants } from "@/uix/location-badge"

export interface LocationNode {
  id: string
  code: string
  name: string
  type?: NonNullable<React.ComponentProps<typeof LocationBadge>["type"]>
  contactCount?: number
  children?: LocationNode[]
}

export interface LocationTreeProps extends Omit<React.ComponentProps<"div">, "onSelect"> {
  data: LocationNode[]
  defaultExpanded?: string[]
  expanded?: string[]
  onExpand?: (id: string, expanded: boolean) => void
  selected?: string[]
  onSelect?: (id: string) => void
  onEdit?: (id: string) => void
  size?: "sm" | "default"
}

function LocationTreeNode({
  node,
  level,
  size,
  expandedSet,
  onToggle,
  selected,
  onSelect,
  onEdit,
}: {
  node: LocationNode
  level: number
  size: "sm" | "default"
  expandedSet: Set<string>
  onToggle: (id: string) => void
  selected: Set<string>
  onSelect: (id: string) => void
  onEdit?: (id: string) => void
}) {
  const hasChildren = node.children && node.children.length > 0
  const isExpanded = expandedSet.has(node.id)
  const isSelected = selected.has(node.id)
  const indent = level * 20

  return (
    <div data-slot="location-tree-node">
      <div
        data-slot="location-tree-node-content"
        className={cn(
          "flex items-center gap-2 rounded-md px-2 py-1.5 cursor-pointer hover:bg-accent hover:text-accent-foreground group",
          isSelected && "bg-accent text-accent-foreground",
          size === "sm" && "text-xs py-1"
        )}
        style={{ paddingLeft: `${indent + 8}px` }}
        onClick={() => {
          if (hasChildren) onToggle(node.id)
          onSelect(node.id)
        }}
        onDoubleClick={() => onEdit?.(node.id)}
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
          <span className="shrink-0 size-4 flex items-center justify-center">
            <MapPinIcon className="size-3 text-muted-foreground" />
          </span>
        )}
        <span className="font-mono text-xs text-muted-foreground">{node.code}</span>
        <span className="flex-1 truncate font-medium">{node.name}</span>
        {node.type && (
          <LocationBadge type={node.type} size="sm" showIcon={false}>
            {node.type}
          </LocationBadge>
        )}
        {node.contactCount !== undefined && (
          <span className="text-xs text-muted-foreground">{node.contactCount} contacts</span>
        )}
      </div>
      {hasChildren && isExpanded && (
        <div data-slot="location-tree-children">
          {node.children!.map((child) => (
            <LocationTreeNode
              key={child.id}
              node={child}
              level={level + 1}
              size={size}
              expandedSet={expandedSet}
              onToggle={onToggle}
              selected={selected}
              onSelect={onSelect}
              onEdit={onEdit}
            />
          ))}
        </div>
      )}
    </div>
  )
}

function LocationTree({
  className,
  data,
  defaultExpanded = [],
  expanded: controlledExpanded,
  onExpand,
  selected: controlledSelected,
  onSelect,
  onEdit,
  size = "default",
  ...props
}: LocationTreeProps) {
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
    <div data-slot="location-tree" className={cn("py-1", className)} {...props}>
      {data.map((node) => (
        <LocationTreeNode
          key={node.id}
          node={node}
          level={0}
          size={size}
          expandedSet={expandedSet}
          onToggle={handleToggle}
          selected={selectedSet}
          onSelect={handleSelect}
          onEdit={onEdit}
        />
      ))}
    </div>
  )
}

export { LocationTree }
