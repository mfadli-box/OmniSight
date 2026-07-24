"use client"

import * as React from "react"
import { cn } from "@/lib/utility"
import { CheckIcon, CircleDashedIcon, XIcon, ClockIcon } from "lucide-react"

export interface ApprovalStep {
  id: string
  label: string
  description?: string
  status: "pending" | "approved" | "rejected" | "waiting"
  approver?: string
  timestamp?: string
}

export interface ApprovalStepsProps extends React.ComponentProps<"div"> {
  steps: ApprovalStep[]
  orientation?: "vertical" | "horizontal"
}

function ApprovalSteps({ className, steps, orientation = "vertical", ...props }: ApprovalStepsProps) {
  return (
    <div
      data-slot="approval-steps"
      className={cn(
        "flex",
        orientation === "vertical" && "flex-col",
        orientation === "horizontal" && "flex-row items-start",
        className
      )}
      {...props}
    >
      {steps.map((step, index) => (
        <div
          key={step.id}
          className={cn(
            "flex gap-3",
            orientation === "vertical" && "flex-row",
            orientation === "horizontal" && "flex-col items-center flex-1",
            index < steps.length - 1 && orientation === "vertical" && "pb-4",
            index < steps.length - 1 && orientation === "horizontal" && "pr-4 last:pr-0"
          )}
        >
          <div className="flex flex-col items-center">
            <div
              className={cn(
                "flex size-8 shrink-0 items-center justify-center rounded-full border-2",
                step.status === "approved" && "border-green-500 bg-green-50 text-green-600 dark:bg-green-950 dark:text-green-400",
                step.status === "rejected" && "border-red-500 bg-red-50 text-red-600 dark:bg-red-950 dark:text-red-400",
                step.status === "pending" && "border-blue-500 bg-blue-50 text-blue-600 dark:bg-blue-950 dark:text-blue-400",
                step.status === "waiting" && "border-gray-300 bg-gray-50 text-gray-400 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-500"
              )}
            >
              {step.status === "approved" && <CheckIcon className="size-4" />}
              {step.status === "rejected" && <XIcon className="size-4" />}
              {step.status === "pending" && <CircleDashedIcon className="size-4" />}
              {step.status === "waiting" && <ClockIcon className="size-4" />}
            </div>
            {index < steps.length - 1 && orientation === "vertical" && (
              <div
                className={cn(
                  "w-0.5 flex-1 my-1",
                  step.status === "approved" ? "bg-green-500" : "bg-gray-200 dark:bg-gray-700"
                )}
              />
            )}
            {index < steps.length - 1 && orientation === "horizontal" && (
              <div
                className={cn(
                  "h-0.5 flex-1 mt-4 mx-1",
                  step.status === "approved" ? "bg-green-500" : "bg-gray-200 dark:bg-gray-700"
                )}
              />
            )}
          </div>
          <div className={cn("flex-1 min-w-0", orientation === "horizontal" && "text-center mt-1")}>
            <p
              className={cn(
                "text-sm font-medium",
                step.status === "approved" && "text-green-700 dark:text-green-400",
                step.status === "rejected" && "text-red-700 dark:text-red-400",
                step.status === "pending" && "text-blue-700 dark:text-blue-400",
                step.status === "waiting" && "text-muted-foreground"
              )}
            >
              {step.label}
            </p>
            {step.description && (
              <p className="text-xs text-muted-foreground mt-0.5">{step.description}</p>
            )}
            {step.approver && (
              <p className="text-xs text-muted-foreground mt-0.5">{step.approver}</p>
            )}
            {step.timestamp && (
              <p className="text-xs text-muted-foreground">{step.timestamp}</p>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}

export { ApprovalSteps }
