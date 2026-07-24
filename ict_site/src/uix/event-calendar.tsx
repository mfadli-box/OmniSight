"use client"

import * as React from "react"
import { cn } from "@/lib/utility"
import { Button } from "./button"
import { ChevronLeftIcon, ChevronRightIcon } from "lucide-react"

export interface CalendarEvent {
  id: string
  title: string
  date: string
  color?: "default" | "blue" | "green" | "yellow" | "red" | "purple"
  description?: string
}

export interface EventCalendarProps extends React.ComponentProps<"div"> {
  events?: CalendarEvent[]
  onDateClick?: (date: Date) => void
  onEventClick?: (event: CalendarEvent) => void
  selectedDate?: Date
}

const eventColors = {
  default: "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
  blue: "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
  green: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
  yellow: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  red: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
  purple: "bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400",
}

const WEEKDAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]

function getDaysInMonth(year: number, month: number) {
  return new Date(year, month + 1, 0).getDate()
}

function getFirstDayOfMonth(year: number, month: number) {
  return new Date(year, month, 1).getDay()
}

function formatDateKey(date: Date): string {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`
}

function EventCalendar({
  className,
  events = [],
  onDateClick,
  onEventClick,
  selectedDate,
  ...props
}: EventCalendarProps) {
  const [currentMonth, setCurrentMonth] = React.useState(() => {
    const now = selectedDate || new Date()
    return { year: now.getFullYear(), month: now.getMonth() }
  })

  const daysInMonth = getDaysInMonth(currentMonth.year, currentMonth.month)
  const firstDay = getFirstDayOfMonth(currentMonth.year, currentMonth.month)
  const today = new Date()

  const eventsByDate = React.useMemo(() => {
    const map = new Map<string, CalendarEvent[]>()
    for (const event of events) {
      const key = event.date.split("T")[0]
      const existing = map.get(key) || []
      existing.push(event)
      map.set(key, existing)
    }
    return map
  }, [events])

  const prevMonth = () => {
    setCurrentMonth((prev) => {
      if (prev.month === 0) return { year: prev.year - 1, month: 11 }
      return { ...prev, month: prev.month - 1 }
    })
  }

  const nextMonth = () => {
    setCurrentMonth((prev) => {
      if (prev.month === 11) return { year: prev.year + 1, month: 0 }
      return { ...prev, month: prev.month + 1 }
    })
  }

  const monthName = new Date(currentMonth.year, currentMonth.month).toLocaleString("default", {
    month: "long",
    year: "numeric",
  })

  const cells: (number | null)[] = []
  for (let i = 0; i < firstDay; i++) cells.push(null)
  for (let i = 1; i <= daysInMonth; i++) cells.push(i)

  return (
    <div data-slot="event-calendar" className={cn("rounded-lg border", className)} {...props}>
      <div className="flex items-center justify-between px-4 py-3 border-b">
        <h3 className="text-sm font-semibold">{monthName}</h3>
        <div className="flex gap-1">
          <Button type="button" variant="ghost" size="icon-xs" onClick={prevMonth}>
            <ChevronLeftIcon className="size-4" />
          </Button>
          <Button type="button" variant="ghost" size="icon-xs" onClick={nextMonth}>
            <ChevronRightIcon className="size-4" />
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-7">
        {WEEKDAYS.map((day) => (
          <div
            key={day}
            className="px-1 py-2 text-center text-xs font-medium text-muted-foreground border-b"
          >
            {day}
          </div>
        ))}
      </div>
      <div className="grid grid-cols-7">
        {cells.map((day, index) => {
          if (day === null) {
            return <div key={`empty-${index}`} className="min-h-16 border-b border-r bg-muted/30" />
          }

          const dateKey = formatDateKey(new Date(currentMonth.year, currentMonth.month, day))
          const dayEvents = eventsByDate.get(dateKey) || []
          const isToday =
            today.getFullYear() === currentMonth.year &&
            today.getMonth() === currentMonth.month &&
            today.getDate() === day
          const isSelected =
            selectedDate &&
            selectedDate.getFullYear() === currentMonth.year &&
            selectedDate.getMonth() === currentMonth.month &&
            selectedDate.getDate() === day

          return (
            <div
              key={day}
              className={cn(
                "min-h-16 border-b border-r p-1 cursor-pointer hover:bg-muted/50 transition-colors",
                isSelected && "bg-primary/10",
                (dayEvents.length > 0 || isSelected) && "ring-1 ring-primary/30"
              )}
              onClick={() => onDateClick?.(new Date(currentMonth.year, currentMonth.month, day))}
            >
              <div
                className={cn(
                  "flex items-center justify-center size-6 rounded-full text-xs",
                  isToday && "bg-primary text-primary-foreground font-bold"
                )}
              >
                {day}
              </div>
              <div className="mt-0.5 space-y-0.5">
                {dayEvents.slice(0, 3).map((event) => (
                  <div
                    key={event.id}
                    className={cn(
                      "truncate rounded px-1 py-0.5 text-[10px] font-medium cursor-pointer",
                      eventColors[event.color || "default"]
                    )}
                    onClick={(e) => {
                      e.stopPropagation()
                      onEventClick?.(event)
                    }}
                    title={event.title}
                  >
                    {event.title}
                  </div>
                ))}
                {dayEvents.length > 3 && (
                  <div className="text-[10px] text-muted-foreground text-center">
                    +{dayEvents.length - 3} more
                  </div>
                )}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}

export { EventCalendar }
