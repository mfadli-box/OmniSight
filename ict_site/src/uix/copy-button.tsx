"use client"

import * as React from "react"
import { cn } from "@/lib/utility"
import { CheckIcon, CopyIcon } from "lucide-react"
import { Button } from "@/uix/button"

export interface CopyButtonProps extends React.ComponentProps<typeof Button> {
  text: string
  copiedText?: string
  copiedDuration?: number
}

function CopyButton({
  text,
  copiedText = "Copied",
  copiedDuration = 2000,
  variant = "ghost",
  size = "icon-sm",
  className,
  ...props
}: CopyButtonProps) {
  const [copied, setCopied] = React.useState(false)

  const handleCopy = React.useCallback(async () => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), copiedDuration)
    } catch {
      // clipboard API not available
    }
  }, [text, copiedDuration])

  return (
    <Button
      variant={variant}
      size={size}
      onClick={handleCopy}
      className={cn("shrink-0", className)}
      {...props}
    >
      {copied ? <CheckIcon className="size-3.5 text-green-600" /> : <CopyIcon className="size-3.5" />}
    </Button>
  )
}

export { CopyButton }
