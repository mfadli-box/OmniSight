import * as React from "react"
import { cn } from "@/lib/utility"

export interface TimelineItem {
  id: string
  title: string
  description?: string
  timestamp?: string
  icon?: React.ReactNode
  status?: "default" | "active" | "success" | "error" | "warning"
}

export interface TimelineProps extends React.ComponentProps<"div"> {
  items: TimelineItem[]
  orientation?: "vertical" | "horizontal"
}

function TimelineItemComponent({
  item,
  isLast,
  orientation,
}: {
  item: TimelineItem
  isLast: boolean
  orientation: "vertical" | "horizontal"
}) {
  return (
    <div
      data-slot="timeline-item"
      className={cn(
        "flex gap-3",
        orientation === "vertical" && "flex-row",
        orientation === "horizontal" && "flex-col items-center flex-1",
        orientation === "vertical" && !isLast && "pb-4"
      )}
    >
      <div className="flex flex-col items-center">
        <div
          className={cn(
            "flex size-8 shrink-0 items-center justify-center rounded-full border-2",
            item.status === "active" && "border-blue-500 bg-blue-50 text-blue-600 dark:bg-blue-950 dark:text-blue-400",
            item.status === "success" && "border-green-500 bg-green-50 text-green-600 dark:bg-green-950 dark:text-green-400",
            item.status === "error" && "border-red-500 bg-red-50 text-red-600 dark:bg-red-950 dark:text-red-400",
            item.status === "warning" && "border-yellow-500 bg-yellow-50 text-yellow-600 dark:bg-yellow-950 dark:text-yellow-400",
            item.status === "default" && "border-gray-300 bg-gray-50 text-gray-600 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-400",
            !item.status && "border-gray-300 bg-gray-50 text-gray-600 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-400"
          )}
        >
          {item.icon || (
            <div
              className={cn(
                "size-2 rounded-full",
                item.status === "active" && "bg-blue-500",
                item.status === "success" && "bg-green-500",
                item.status === "error" && "bg-red-500",
                item.status === "warning" && "bg-yellow-500",
                item.status === "default" && "bg-gray-400",
                !item.status && "bg-gray-400"
              )}
            />
          )}
        </div>
        {!isLast && orientation === "vertical" && (
          <div className="w-0.5 flex-1 bg-gray-200 dark:bg-gray-700 my-1" />
        )}
      </div>
      <div className={cn("flex-1 min-w-0", orientation === "horizontal" && "text-center")}>
        <p className="text-sm font-medium">{item.title}</p>
        {item.description && (
          <p className="text-sm text-muted-foreground mt-0.5">{item.description}</p>
        )}
        {item.timestamp && (
          <p className="text-xs text-muted-foreground mt-1">{item.timestamp}</p>
        )}
      </div>
    </div>
  )
}

function Timeline({ className, items, orientation = "vertical", ...props }: TimelineProps) {
  return (
    <div
      data-slot="timeline"
      className={cn(
        "flex",
        orientation === "vertical" && "flex-col",
        orientation === "horizontal" && "flex-row items-start",
        className
      )}
      {...props}
    >
      {items.map((item, index) => (
        <TimelineItemComponent
          key={item.id}
          item={item}
          isLast={index === items.length - 1}
          orientation={orientation}
        />
      ))}
    </div>
  )
}

export { Timeline }
