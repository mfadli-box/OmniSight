"use client";

import { ReactNode } from "react";
import { Button } from "./button";
import {
  Dialog, DialogTitle, DialogDescription,
  DialogContent, DialogFooter, DialogHeader
} from "./dialog";

interface ActionConfig {
  label: string;
  type?: "button" | "submit";
  variant?: string;
  disabled?: boolean;
}
interface DataDialogProps {
  isOpen: boolean;
  mode: "create" | "update" | "delete" | "detail";
  title: string;
  description?: string;
  onClose: () => void;
  onSubmit?: () => void;
  children: ReactNode;
}
export default function DataDialog({
  isOpen,
  mode,
  title,
  description,
  onClose,
  onSubmit,
  children,
}: DataDialogProps) {
  if (!isOpen) return null;
  const ContentWrapper = "div";
  const colorCreate = "bg-green-500 text-white hover:bg-green-600";
  const colorUpdate = "bg-blue-500 text-white hover:bg-blue-600";
  const colorDelete = "bg-red-500 text-white hover:bg-red-600";
  const colorDetail = "bg-gray-500 text-white hover:bg-gray-600";
  const cancelAction: ActionConfig = {
    label: "Cancel",
    type: "button",
    variant: "bg-gray-200 text-gray-800 hover:bg-gray-300",
  };
  const submitAction: ActionConfig = {
    label: mode.charAt(0).toUpperCase() + mode.slice(1),
    type: "submit",
    variant: mode === "create" ? colorCreate
           : mode === "update" ? colorUpdate
           : mode === "delete" ? colorDelete : colorDetail,
  };
  const widthClass = "lg:max-w-5xl sm:max-w-2xl";
  return (
    <Dialog open={isOpen} >
      <DialogContent className={`${widthClass} scrollable`} showCloseButton={false} >
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          {description && (
            <DialogDescription>{description}</DialogDescription>
          )}
        </DialogHeader>
        <ContentWrapper onSubmit={onSubmit} className="min-w-0 overflow-hidden max-h-[85vh] overflow-y-auto">
          {children}
          <DialogFooter showCloseButton={false} className="p-2">
            {mode !== "detail" && submitAction && (
              <Button
                type={submitAction.type || "submit"}
                disabled={submitAction.disabled}
                onClick={onSubmit}
                className={`${submitAction.variant}`}
              >
                {submitAction.label}
              </Button>
            )}
            {cancelAction ? (
              <Button
                type={cancelAction.type || "button"}
                disabled={cancelAction.disabled}
                onClick={onClose}
                className={`${cancelAction.variant}`}
              >
                {cancelAction.label}
              </Button>
            ) : null}
          </DialogFooter>
        </ContentWrapper>
      </DialogContent>
    </Dialog>
  );
}
