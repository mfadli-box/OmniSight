import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utility"

const statusBadgeVariants = cva(
  "inline-flex items-center gap-1.5 whitespace-nowrap font-medium",
  {
    variants: {
      status: {
        default: "",
        active: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
        inactive: "bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400",
        pending: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
        warning: "bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400",
        error: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
        info: "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
        success: "bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400",
        destructive: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
      },
      severity: {
        critical: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
        high: "bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400",
        medium: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
        low: "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
        info: "bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400",
      },
      size: {
        sm: "text-xs px-1.5 py-0.5",
        default: "text-xs px-2 py-0.5",
        lg: "text-sm px-2.5 py-1",
      },
    },
    defaultVariants: {
      status: "default",
      size: "default",
    },
  }
)

export interface StatusBadgeProps
  extends React.ComponentProps<"span">,
    VariantProps<typeof statusBadgeVariants> {
  dot?: boolean
}

function StatusBadge({ className, status, severity, size, dot = true, children, ...props }: StatusBadgeProps) {
  const variant = severity || status || "default"

  return (
    <span
      data-slot="status-badge"
      className={cn(statusBadgeVariants({ status: severity ? undefined : status, severity, size }), className)}
      {...props}
    >
      {dot && (
        <span
          className={cn(
            "size-1.5 rounded-full",
            severity === "critical" && "bg-red-500",
            severity === "high" && "bg-orange-500",
            severity === "medium" && "bg-yellow-500",
            severity === "low" && "bg-blue-500",
            severity === "info" && "bg-gray-500",
            status === "active" && "bg-green-500",
            status === "inactive" && "bg-gray-500",
            status === "pending" && "bg-yellow-500",
            status === "warning" && "bg-orange-500",
            status === "error" && "bg-red-500",
            status === "info" && "bg-blue-500",
            status === "success" && "bg-green-500",
            status === "destructive" && "bg-red-500",
            !status && !severity && "bg-current"
          )}
        />
      )}
      {children}
    </span>
  )
}

export { StatusBadge, statusBadgeVariants }
