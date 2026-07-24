"use client";

import { Suspense, useCallback, useMemo, useState } from "react";
import { storageKey, parseSession } from "@/lib/utility";
import { clientApi, ClientApiError } from "@/lib/client-api";
import { getErrorMessage } from "@/lib/error-message";
import { toast } from "sonner";
import DataTable, { Column, ActionConfig } from "@/uix/datatable";
import DataDialog from "@/uix/datadialog";
import { Field, FieldGroup, FieldLabel } from "@/uix/field";
import { Input } from "@/uix/input";

interface SessionListItem {
  id: string;
  user_id: string;
  username: string;
  fullname: string;
  company_name: string;
  ip_address: string;
  user_agent: string;
  token_preview: string;
  created_at: string;
  expires_at: string;
  is_active: boolean;
}

interface SessionDetailItem {
  id: string;
  user_id: string;
  username: string;
  fullname: string;
  email: string;
  phone: string;
  role: string;
  is_admin: boolean;
  is_hris: boolean;
  is_active: boolean;
  company_id: string;
  company_name: string;
  token: string;
  ip_address: string;
  user_agent: string;
  created_at: string;
  expires_at: string;
  session_active: boolean;
}

export default function Page() {
  const [refresh, setRefresh] = useState(false);
  const [detailModal, setDetailModal] = useState<{
    isOpen: boolean;
    data: SessionDetailItem | null;
  }>({ isOpen: false, data: null });

  const loadData = useCallback(async (params: {
    search: string;
    page: number;
    size: number;
    sort_by: string;
    sort_order: "asc" | "desc";
  }) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
    try {
      const data = await clientApi<{
        data: SessionListItem[];
        meta: { total: number; page: number; size: number }
      }>("/SM05", {
        params: {
          search: params.search,
          page: String(params.page),
          size: String(params.size),
          sort_by: params.sort_by,
          sort_order: params.sort_order,
        },
      });
      return {
        data: data?.data ?? [],
        meta: data?.meta ?? { total: 0, page: 1, size: 10 },
      };
    } catch {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
  }, []);

  const openDetail = async (row: SessionListItem) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      const res = await clientApi<{ data: SessionDetailItem }>("/SM05/" + row.id);
      setDetailModal({ isOpen: true, data: (res as any)?.data ?? null });
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(
        err instanceof ClientApiError ? getErrorMessage(err.code) : "Failed to load session detail.",
      );
    }
  };

  const handleRevoke = async (row: SessionListItem) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM05/" + row.id, { method: "DELETE" });
      toast.success("Session revoked successfully.");
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(
        err instanceof ClientApiError ? err.message : "Failed to revoke session.",
      );
    }
  };

  const columns: Column<SessionListItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "User ID", accessor: "user_id", sortable: false, hidden: true },
    { header: "Username", accessor: "username", sortable: true, hidden: false },
    { header: "Fullname", accessor: "fullname", sortable: true, hidden: false },
    { header: "Company", accessor: "company_name", sortable: false, hidden: false },
    { header: "IP Address", accessor: "ip_address", sortable: true, hidden: false },
    {
      header: "Token",
      accessor: "token_preview",
      sortable: false,
      hidden: false,
      formatter: (v: string) => (
        <code className="text-xs bg-muted px-1.5 py-0.5 rounded font-mono">{v}</code>
      ),
    },
    { header: "User Agent", accessor: "user_agent", sortable: false, hidden: true },
    {
      header: "Created",
      accessor: "created_at",
      sortable: true,
      hidden: false,
      formatter: (v: string) => new Date(v).toLocaleString("id-ID"),
    },
    {
      header: "Expires",
      accessor: "expires_at",
      sortable: true,
      hidden: false,
      formatter: (v: string) => new Date(v).toLocaleString("id-ID"),
    },
    {
      header: "Status",
      accessor: "is_active",
      sortable: false,
      hidden: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-red-600 font-medium"}>
          {v ? "Active" : "Expired"}
        </span>
      ),
    },
  ], []);

  const actions: ActionConfig<SessionListItem> = {
    onSearch: () => setRefresh(!refresh),
    hideCreate: true,
    hideUpdate: true,
    onDetail: openDetail,
    onDelete: handleRevoke,
  };

  return (
    <Suspense fallback={<p>Loading sessions...</p>}>
      <div className="p-3 max-w-8xl mx-auto space-y-3">
        <DataTable
          fetchData={loadData}
          columns={columns}
          actions={actions}
          hideSearch={false}
          hideSelect={true}
          hidePaging={false}
          hideSort={false}
          hideColumnToggle={true}
          refreshTrigger={refresh}
        />

        <DataDialog
          isOpen={detailModal.isOpen}
          mode="detail"
          title="Session Detail"
          onClose={() => setDetailModal({ isOpen: false, data: null })}
        >
          {detailModal.data && (
            <FieldGroup className="p-3 max-h-[70vh] overflow-y-auto space-y-4">
              <div>
                <p className="text-sm font-medium text-muted-foreground mb-2">User Information</p>
                <div className="grid gap-3 sm:grid-cols-2">
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Username</FieldLabel>
                    <Input value={detailModal.data.username} disabled readOnly />
                  </Field>
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Fullname</FieldLabel>
                    <Input value={detailModal.data.fullname} disabled readOnly />
                  </Field>
                </div>
                <div className="grid gap-3 sm:grid-cols-2 mt-3">
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Email</FieldLabel>
                    <Input value={detailModal.data.email} disabled readOnly />
                  </Field>
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Phone</FieldLabel>
                    <Input value={detailModal.data.phone || "—"} disabled readOnly />
                  </Field>
                </div>
                <div className="grid gap-3 sm:grid-cols-2 mt-3">
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Role</FieldLabel>
                    <Input value={detailModal.data.role} disabled readOnly />
                  </Field>
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Company</FieldLabel>
                    <Input value={detailModal.data.company_name || "—"} disabled readOnly />
                  </Field>
                </div>
                <div className="grid gap-3 sm:grid-cols-3 mt-3">
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Admin</FieldLabel>
                    <Input value={detailModal.data.is_admin ? "Yes" : "No"} disabled readOnly />
                  </Field>
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>HRIS</FieldLabel>
                    <Input value={detailModal.data.is_hris ? "Yes" : "No"} disabled readOnly />
                  </Field>
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Account Active</FieldLabel>
                    <Input value={detailModal.data.is_active ? "Yes" : "No"} disabled readOnly />
                  </Field>
                </div>
              </div>

              <div className="border-t pt-3">
                <p className="text-sm font-medium text-muted-foreground mb-2">Session Information</p>
                <div className="grid gap-3 sm:grid-cols-2">
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Token</FieldLabel>
                    <Input value={detailModal.data.token} disabled readOnly className="font-mono text-xs" />
                  </Field>
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Status</FieldLabel>
                    <Input
                      value={detailModal.data.session_active ? "Active" : "Expired"}
                      disabled readOnly
                      className={detailModal.data.session_active ? "text-emerald-600" : "text-red-600"}
                    />
                  </Field>
                </div>
                <div className="grid gap-3 sm:grid-cols-2 mt-3">
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>IP Address</FieldLabel>
                    <Input value={detailModal.data.ip_address || "—"} disabled readOnly />
                  </Field>
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>User Agent</FieldLabel>
                    <Input value={detailModal.data.user_agent || "—"} disabled readOnly />
                  </Field>
                </div>
                <div className="grid gap-3 sm:grid-cols-2 mt-3">
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Created At</FieldLabel>
                    <Input value={detailModal.data.created_at} disabled readOnly />
                  </Field>
                  <Field className="gap-1.5 px-2">
                    <FieldLabel>Expires At</FieldLabel>
                    <Input value={detailModal.data.expires_at} disabled readOnly />
                  </Field>
                </div>
              </div>
            </FieldGroup>
          )}
        </DataDialog>
      </div>
    </Suspense>
  );
}
