"use client";

import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { storageKey, parseSession, isSessionExpired, SessionData } from "@/lib/utility";
import { clientApi, ClientApiError } from "@/lib/client-api";
import { toast } from "sonner";
import { getErrorMessage } from "@/lib/error-message";
import { Button } from "@/uix/button";
import { Input } from "@/uix/input";
import { Field, FieldError, FieldGroup, FieldLabel, FieldSet } from "@/uix/field";

export default function Page() {
  const router = useRouter();
  const [session, setSession] = useState<SessionData | null>(null);
  useEffect(() => {
    const stored = parseSession(window.localStorage.getItem(storageKey));
    if (!stored || isSessionExpired(stored)) {
      window.localStorage.removeItem(storageKey);
      try {
        document.cookie = `${storageKey}=; path=/; max-age=0`;
      } catch { }
      router.replace("/login");
      return;
    }
    setSession(stored);
  }, [router]);
  const [isHRIS, setIsHRIS] = useState(false);
  useEffect(() => {
    setIsHRIS(Boolean(session?.user_profile.is_hris));
  }, [session]);
  const [passwordBusy, setPasswordBusy] = useState(false);
  const forceClientLogout = () => {
    window.localStorage.removeItem(storageKey);
    try {
      document.cookie = storageKey +"=; path=/; max-age=0";
    } catch (e) {}
    router.replace("/login");
  };
  type passwordInput = z.infer<typeof passwordSchema>;
  const passwordSchema = z
    .object({
      "current_password": z.string().min(1, "Current password is required"),
      "new_password": z
        .string()
        .min(6, "Password must be at least 6 characters long")
        .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
        .regex(/[a-z]/, "Password must contain at least one lowercase letter")
        .regex(/[0-9]/, "Password must contain at least one number"),
      "confirm_password": z.string().min(1, "Please confirm your new password"),
    })
    .superRefine(({ "new_password": newPassword, "confirm_password": confirmPassword }, ctx) => {
      if (newPassword !== confirmPassword) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: "The passwords do not match",
          path: ["confirm_password"],
        });
      }
    });
  const formPassword = useForm<passwordInput>({
    resolver: zodResolver(passwordSchema),
    defaultValues: { "current_password": "", "new_password": "", "confirm_password": "" },
  });
  const handlePassword = async (values: passwordInput) => {
    if (!session) return;
    setPasswordBusy(true);
    try {
      const data = await clientApi<{ message?: string }>("/SP02", {
        method: "PUT",
        body: values,
      });
      toast.success(data?.message || "Password updated successfully.");
      formPassword.reset();
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 403) {
        forceClientLogout();
        return;
      }
      toast.error(
        getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR")
      );
    } finally {
      setPasswordBusy(false);
    }
  };
  return (
    <>
    {!isHRIS ? (
      <main className="p-3 max-w-8xl mx-auto">
        <form noValidate className="space-y-4 mt-4" onSubmit={formPassword.handleSubmit(handlePassword)}>
          <FieldSet>
            <FieldGroup>
              <Controller
                control={formPassword.control}
                name="current_password"
                render={({ field, fieldState }) => (
                  <>
                    <Field className="grid grid-cols-1 items-start gap-3 sm:grid-cols-[10rem_minmax(0,1fr)]" data-invalid={fieldState.invalid}>
                      <FieldLabel className="pt-2" htmlFor="current_password">Current Password</FieldLabel>
                      <Input
                        {...field}
                        id="current_password"
                        type="password"
                        placeholder="Current Password"
                        aria-invalid={fieldState.invalid}
                        value={field.value ?? ""}
                      />
                    </Field>
                    {fieldState.invalid && <FieldError className="p-0 pl-5" errors={[fieldState.error]} />}
                  </>
                )}
              />
              <Controller
                control={formPassword.control}
                name="new_password"
                render={({ field, fieldState }) => (
                  <>
                    <Field className="grid grid-cols-1 items-start gap-3 sm:grid-cols-[10rem_minmax(0,1fr)]" data-invalid={fieldState.invalid}>
                      <FieldLabel className="pt-2" htmlFor="new_password">New Password</FieldLabel>
                      <Input
                        {...field}
                        id="new_password"
                        type="password"
                        placeholder="New Password"
                        aria-invalid={fieldState.invalid}
                        value={field.value ?? ""}
                      />
                    </Field>
                    {fieldState.invalid && <FieldError className="p-0 pl-5" errors={[fieldState.error]} />}
                  </>
                )}
              />
              <Controller
                control={formPassword.control}
                name="confirm_password"
                render={({ field, fieldState }) => (
                  <>
                    <Field className="grid grid-cols-1 items-start gap-3 sm:grid-cols-[10rem_minmax(0,1fr)]" data-invalid={fieldState.invalid}>
                      <FieldLabel className="pt-2" htmlFor="confirm_password">Confirm New Password</FieldLabel>
                      <Input
                        {...field}
                        id="confirm_password"
                        type="password"
                        placeholder="Confirm New Password"
                        aria-invalid={fieldState.invalid}
                        value={field.value ?? ""}
                      />
                    </Field>
                    {fieldState.invalid && <FieldError className="p-0 pl-5" errors={[fieldState.error]} />}
                  </>
                )}
              />
            </FieldGroup>
            <Button
              type="submit"
              disabled={passwordBusy}
            >
              {passwordBusy ? "Updating..." : "Change Password"}
            </Button>
          </FieldSet>
        </form>
      </main>
    ) : (
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div className="space-y-1">
            <p className="text-muted-foreground text-sm">
              You are using an HRIS account, you cannot change your password here.
            </p>
          </div>
        </div>
      </div>
    )}
    </>
  );
}
