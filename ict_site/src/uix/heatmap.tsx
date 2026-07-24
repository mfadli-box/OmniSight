import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utility"

export interface HeatmapCell {
  id: string
  label?: string
  value: number
  disabled?: boolean
}

export interface HeatmapProps extends React.ComponentProps<"div"> {
  rows: string[]
  columns: string[]
  data: (number | null)[][]
  cellLabels?: string[][]
  colorScale?: "default" | "severity" | "risk"
  size?: "sm" | "default" | "lg"
  showLabels?: boolean
  showValues?: boolean
  onCellClick?: (row: number, col: number, value: number | null) => void
}

const colorScales = {
  default: [
    "bg-green-100 dark:bg-green-900/30",
    "bg-green-200 dark:bg-green-800/40",
    "bg-yellow-200 dark:bg-yellow-800/40",
    "bg-orange-200 dark:bg-orange-800/40",
    "bg-red-200 dark:bg-red-800/40",
  ],
  severity: [
    "bg-blue-100 dark:bg-blue-900/30",
    "bg-green-100 dark:bg-green-900/30",
    "bg-yellow-100 dark:bg-yellow-900/30",
    "bg-orange-100 dark:bg-orange-900/30",
    "bg-red-100 dark:bg-red-900/30",
  ],
  risk: [
    "bg-green-100 dark:bg-green-900/30",
    "bg-yellow-100 dark:bg-yellow-900/30",
    "bg-orange-100 dark:bg-orange-900/30",
    "bg-red-100 dark:bg-red-900/30",
    "bg-red-200 dark:bg-red-800/40",
  ],
}

const sizeClasses = {
  sm: "size-8 text-xs",
  default: "size-10 text-xs",
  lg: "size-12 text-sm",
}

function getColorIndex(value: number | null, max: number, colorScale: "default" | "severity" | "risk"): number {
  if (value === null) return -1
  if (value === 0) return 0
  const ratio = value / max
  const scale = colorScales[colorScale]
  const index = Math.min(Math.floor(ratio * scale.length), scale.length - 1)
  return Math.max(0, index)
}

function Heatmap({
  className,
  rows,
  columns,
  data,
  colorScale = "default",
  size = "default",
  showLabels = true,
  showValues = true,
  onCellClick,
  ...props
}: HeatmapProps) {
  const maxValue = React.useMemo(() => {
    let max = 0
    for (const row of data) {
      for (const cell of row) {
        if (cell !== null && cell > max) max = cell
      }
    }
    return max
  }, [data])

  const scale = colorScales[colorScale]

  return (
    <div data-slot="heatmap" className={cn("overflow-x-auto", className)} {...props}>
      <div className="inline-flex flex-col gap-0.5">
        {showLabels && (
          <div className="flex items-end gap-0.5 ml-auto" style={{ marginLeft: `${rows.length > 0 ? 80 : 0}px` }}>
            {columns.map((col, i) => (
              <div
                key={i}
                className={cn("flex items-end justify-center font-medium text-muted-foreground", sizeClasses[size])}
                style={{ writingMode: "vertical-rl", transform: "rotate(180deg)" }}
                title={col}
              >
                <span className="truncate">{col}</span>
              </div>
            ))}
          </div>
        )}
        {rows.map((row, rowIndex) => (
          <div key={rowIndex} className="flex items-center gap-0.5">
            {showLabels && (
              <div
                className={cn(
                  "flex items-center justify-end pr-2 font-medium text-muted-foreground truncate",
                  size === "sm" && "w-16 text-xs",
                  size === "default" && "w-20 text-xs",
                  size === "lg" && "w-24 text-sm"
                )}
                title={row}
              >
                {row}
              </div>
            )}
            {columns.map((_, colIndex) => {
              const value = data[rowIndex]?.[colIndex] ?? null
              const colorIdx = getColorIndex(value, maxValue, colorScale)
              return (
                <div
                  key={colIndex}
                  className={cn(
                    "flex items-center justify-center rounded-sm border border-transparent font-medium transition-colors",
                    sizeClasses[size],
                    colorIdx >= 0 ? scale[colorIdx] : "bg-muted",
                    onCellClick && "cursor-pointer hover:ring-2 hover:ring-primary/50"
                  )}
                  onClick={() => onCellClick?.(rowIndex, colIndex, value)}
                  title={value !== null ? `${row} × ${columns[colIndex]}: ${value}` : `${row} × ${columns[colIndex]}: N/A`}
                >
                  {showValues && value !== null && value}
                </div>
              )
            })}
          </div>
        ))}
      </div>
    </div>
  )
}

export { Heatmap, colorScales }
