"use client"

import * as React from "react"
import { cn } from "@/lib/utility"
import { ScrollArea } from "./scroll-area"

export interface KanbanCard {
  id: string
  title: string
  description?: string
  badge?: React.ReactNode
  assignee?: React.ReactNode
  dueDate?: string
  priority?: "low" | "medium" | "high"
}

export interface KanbanColumn {
  id: string
  title: string
  color?: "default" | "blue" | "green" | "yellow" | "red"
  cards: KanbanCard[]
}

export interface KanbanBoardProps extends React.ComponentProps<"div"> {
  columns: KanbanColumn[]
  onCardClick?: (card: KanbanCard, column: KanbanColumn) => void
  onCardMove?: (cardId: string, fromColumn: string, toColumn: string) => void
  renderCard?: (card: KanbanCard, column: KanbanColumn) => React.ReactNode
}

const columnColors = {
  default: "bg-muted",
  blue: "bg-blue-50 dark:bg-blue-950/30",
  green: "bg-green-50 dark:bg-green-950/30",
  yellow: "bg-yellow-50 dark:bg-yellow-950/30",
  red: "bg-red-50 dark:bg-red-950/30",
}

const headerColors = {
  default: "bg-muted/50",
  blue: "bg-blue-100 dark:bg-blue-900/50",
  green: "bg-green-100 dark:bg-green-900/50",
  yellow: "bg-yellow-100 dark:bg-yellow-900/50",
  red: "bg-red-100 dark:bg-red-900/50",
}

function KanbanBoard({
  className,
  columns,
  onCardClick,
  onCardMove,
  renderCard,
  ...props
}: KanbanBoardProps) {
  const [draggedCard, setDraggedCard] = React.useState<{ cardId: string; columnId: string } | null>(null)
  const [dragOverColumn, setDragOverColumn] = React.useState<string | null>(null)

  const handleDragStart = (cardId: string, columnId: string) => {
    setDraggedCard({ cardId, columnId })
  }

  const handleDragOver = (e: React.DragEvent, columnId: string) => {
    e.preventDefault()
    setDragOverColumn(columnId)
  }

  const handleDragLeave = () => {
    setDragOverColumn(null)
  }

  const handleDrop = (e: React.DragEvent, toColumnId: string) => {
    e.preventDefault()
    if (draggedCard && draggedCard.columnId !== toColumnId) {
      onCardMove?.(draggedCard.cardId, draggedCard.columnId, toColumnId)
    }
    setDraggedCard(null)
    setDragOverColumn(null)
  }

  return (
    <div
      data-slot="kanban-board"
      className={cn("flex gap-4 overflow-x-auto pb-4", className)}
      {...props}
    >
      {columns.map((column) => (
        <div
          key={column.id}
          data-slot="kanban-column"
          className={cn(
            "flex flex-col w-72 shrink-0 rounded-lg",
            columnColors[column.color || "default"],
            dragOverColumn === column.id && "ring-2 ring-primary/50"
          )}
          onDragOver={(e) => handleDragOver(e, column.id)}
          onDragLeave={handleDragLeave}
          onDrop={(e) => handleDrop(e, column.id)}
        >
          <div
            className={cn(
              "flex items-center justify-between rounded-t-lg px-3 py-2",
              headerColors[column.color || "default"]
            )}
          >
            <h3 className="text-sm font-semibold">{column.title}</h3>
            <span className="text-xs text-muted-foreground bg-background/50 rounded-full px-2 py-0.5">
              {column.cards.length}
            </span>
          </div>
          <ScrollArea className="flex-1 max-h-[calc(100vh-200px)]">
            <div className="flex flex-col gap-2 p-2">
              {column.cards.map((card) => (
                <div
                  key={card.id}
                  data-slot="kanban-card"
                  draggable
                  onDragStart={() => handleDragStart(card.id, column.id)}
                  onClick={() => onCardClick?.(card, column)}
                  className={cn(
                    "rounded-md border bg-card p-3 shadow-sm cursor-pointer hover:shadow-md transition-shadow",
                    draggedCard?.cardId === card.id && "opacity-50"
                  )}
                >
                  {renderCard ? (
                    renderCard(card, column)
                  ) : (
                    <>
                      <div className="flex items-start justify-between gap-2">
                        <p className="text-sm font-medium">{card.title}</p>
                        {card.badge}
                      </div>
                      {card.description && (
                        <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                          {card.description}
                        </p>
                      )}
                      <div className="flex items-center justify-between mt-2">
                        {card.assignee}
                        {card.dueDate && (
                          <span className="text-xs text-muted-foreground">{card.dueDate}</span>
                        )}
                      </div>
                    </>
                  )}
                </div>
              ))}
              {column.cards.length === 0 && (
                <div className="flex items-center justify-center py-8 text-sm text-muted-foreground">
                  No cards
                </div>
              )}
            </div>
          </ScrollArea>
        </div>
      ))}
    </div>
  )
}

export { KanbanBoard }
