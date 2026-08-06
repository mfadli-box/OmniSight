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
import { Checkbox } from "@/uix/checkbox";

interface XX99Item {
  id: string;
  company_id: string;
  code: string;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

interface XX99Form {
  company_id: string;
  code: string;
  name: string;
  is_active: boolean;
}

export default function Page() {
  const [refresh, setRefresh] = useState(false);
  const [loading, setLoading] = useState(false);
  const [selectedRow, setSelectedRow] = useState<XX99Item | null>(null);
  const [form, setForm] = useState<XX99Form>({
    company_id: "",
    code: "",
    name: "",
    is_active: true,
  });
  const [modal, setModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update" | "detail";
    title: string;
  }>({
    isOpen: false,
    mode: "create",
    title: "",
  });

  const closeModal = () => {
    setModal({ isOpen: false, mode: "create", title: "" });
    setSelectedRow(null);
    setForm({ company_id: "", code: "", name: "", is_active: true });
  };

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
        data: XX99Item[];
        meta: { total: number; page: number; size: number };
      }>("/XX99", {
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

  const openCreate = () => {
    setSelectedRow(null);
    setForm({ company_id: "", code: "", name: "", is_active: true });
    setModal({ isOpen: true, mode: "create", title: "Create XX99" });
  };

  const openUpdate = (row: XX99Item) => {
    setSelectedRow(row);
    setForm({
      company_id: row.company_id ?? "",
      code: row.code ?? "",
      name: row.name ?? "",
      is_active: row.is_active,
    });
    setModal({ isOpen: true, mode: "update", title: "Update XX99" });
  };

  const openDetail = (row: XX99Item) => {
    setSelectedRow(row);
    setModal({ isOpen: true, mode: "detail", title: "XX99 Details" });
  };

  const handleDelete = async (row: XX99Item) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi(`/XX99/${row.id}`, { method: "DELETE" });
      toast.success("XX99 deleted successfully.");
      setRefresh((prev) => !prev);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(
        err instanceof ClientApiError
          ? getErrorMessage(err.code)
          : "Failed to delete XX99.",
      );
    }
  };

  const handleSubmit = async () => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    if (!form.company_id || !form.code || !form.name) {
      toast.error("Company ID, code, dan name wajib diisi.");
      return;
    }

    setLoading(true);
    try {
      if (modal.mode === "create") {
        await clientApi("/XX99", { method: "POST", body: form });
        toast.success("XX99 created successfully.");
      } else if (modal.mode === "update" && selectedRow) {
        await clientApi(`/XX99/${selectedRow.id}`, { method: "PUT", body: form });
        toast.success("XX99 updated successfully.");
      }
      closeModal();
      setRefresh((prev) => !prev);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(
        err instanceof ClientApiError
          ? getErrorMessage(err.code)
          : "Failed to submit XX99.",
      );
    } finally {
      setLoading(false);
    }
  };

  const columns: Column<XX99Item>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Company", accessor: "company_id", sortable: true, hidden: false },
    { header: "Code", accessor: "code", sortable: true, hidden: false },
    { header: "Name", accessor: "name", sortable: true, hidden: false },
    {
      header: "Active",
      accessor: "is_active",
      sortable: false,
      hidden: false,
      formatter: (value: boolean) => (
        <span className={value ? "text-emerald-600 font-medium" : "text-red-600 font-medium"}>
          {value ? "Active" : "Inactive"}
        </span>
      ),
    },
    {
      header: "Created",
      accessor: "created_at",
      sortable: true,
      hidden: false,
      formatter: (value: string) => (value ? new Date(value).toLocaleString("id-ID") : "-"),
    },
  ], []);

  const actions: ActionConfig<XX99Item> = {
    onSearch: () => setRefresh((prev) => !prev),
    onCreate: openCreate,
    onUpdate: openUpdate,
    onDetail: openDetail,
    onDelete: handleDelete,
  };

  return (
    <Suspense fallback={<p>Loading XX99...</p>}>
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
          isOpen={modal.isOpen}
          mode={modal.mode}
          title={modal.title}
          onClose={closeModal}
          onSubmit={handleSubmit}
          loading={loading}
        >
          {modal.mode === "detail" && selectedRow ? (
            <FieldGroup className="p-2">
              {([
                ["ID", selectedRow.id],
                ["Company ID", selectedRow.company_id],
                ["Code", selectedRow.code],
                ["Name", selectedRow.name],
                ["Status", selectedRow.is_active ? "Active" : "Inactive"],
                ["Created", selectedRow.created_at ? new Date(selectedRow.created_at).toLocaleString("id-ID") : "-"],
                ["Updated", selectedRow.updated_at ? new Date(selectedRow.updated_at).toLocaleString("id-ID") : "-"],
              ] as [string, string][]).map(([label, value]) => (
                <div
                  key={label}
                  className="flex flex-col gap-1 border-b p-1 last:border-b-0 sm:flex-row sm:items-start sm:gap-3"
                >
                  <FieldLabel className="w-32 shrink-0 text-sm text-muted-foreground">{label}</FieldLabel>
                  <span className="text-sm break-all">{value || "-"}</span>
                </div>
              ))}
            </FieldGroup>
          ) : (
            <FieldGroup className="p-3 grid gap-3">
              <Field>
                <FieldLabel htmlFor="xx99-company-id">Company ID</FieldLabel>
                <Input
                  id="xx99-company-id"
                  value={form.company_id}
                  onChange={(e) => setForm((prev) => ({ ...prev, company_id: e.target.value }))}
                  placeholder="company uuid"
                />
              </Field>
              <Field>
                <FieldLabel htmlFor="xx99-code">Code</FieldLabel>
                <Input
                  id="xx99-code"
                  value={form.code}
                  onChange={(e) => setForm((prev) => ({ ...prev, code: e.target.value }))}
                  placeholder="XX99-001"
                />
              </Field>
              <Field>
                <FieldLabel htmlFor="xx99-name">Name</FieldLabel>
                <Input
                  id="xx99-name"
                  value={form.name}
                  onChange={(e) => setForm((prev) => ({ ...prev, name: e.target.value }))}
                  placeholder="Example Name"
                />
              </Field>
              <Field className="flex-row items-center gap-2 px-2">
                <Checkbox
                  checked={form.is_active}
                  onCheckedChange={(checked) =>
                    setForm((prev) => ({ ...prev, is_active: Boolean(checked) }))
                  }
                />
                <FieldLabel className="mb-0">Is Active</FieldLabel>
              </Field>
            </FieldGroup>
          )}
        </DataDialog>
      </div>
    </Suspense>
  );
}