"use client";

import { Controller, useForm } from "react-hook-form";
import { Suspense, useCallback, useEffect, useMemo, useState } from "react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { storageKey, parseSession } from "@/lib/utility";
import { clientApi, ClientApiError } from "@/lib/client-api";
import { toast } from "sonner";
import { getErrorMessage } from "@/lib/error-message";
import DataTable, { Column, ActionConfig } from "@/uix/datatable";
import DataDialog from "@/uix/datadialog";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
} from "@/uix/input-group";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/uix/select";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldSet,
} from "@/uix/field";
import { Checkbox } from "@/uix/checkbox";
import { Input } from "@/uix/input";

interface ModuleListItem {
  id: string;
  parent_id: string;
  code: string;
  name: string;
  path: string;
  is_page: boolean;
  is_active: boolean;
  created_at: string;
}

const createSchema = z.object({
  parent_id: z.string().optional(),
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  path: z.string().min(1, "Path is required"),
  is_page: z.boolean().optional(),
  is_active: z.boolean().optional(),
});

const updateSchema = z.object({
  parent_id: z.string().optional(),
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  path: z.string().min(1, "Path is required"),
  is_page: z.boolean().optional(),
  is_active: z.boolean().optional(),
});

type CreateForm = z.infer<typeof createSchema>;
type UpdateForm = z.infer<typeof updateSchema>;

export default function Page() {
  const [selectedRow, setSelectedRow] = useState<ModuleListItem | null>(null);
  const [refresh, setRefresh] = useState(false);
  const [modal, setModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update" | "detail";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });
  const [allModules, setAllModules] = useState<ModuleListItem[]>([]);

  useEffect(() => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    clientApi<ModuleListItem[]>("/SM02", {
      params: { page: "1", size: "100", sort_by: "code", sort_order: "asc" },
    })
      .then((res) => {
        setAllModules((res as any)?.data ?? []);
      })
      .catch(() => {
        setAllModules([]);
      });
  }, []);

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
        data: ModuleListItem[];
        meta: { total: number; page: number; size: number }
      }>("/SM02", {
        params: {
          page: String(params.page),
          size: String(params.size),
          sort_by: params.sort_by,
          sort_order: params.sort_order,
          search: params.search,
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

  const columns: Column<ModuleListItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Parent ID", accessor: "parent_id", sortable: false, hidden: true },
    { header: "Code", accessor: "code", sortable: true, hidden: false },
    { header: "Name", accessor: "name", sortable: true, hidden: false },
    { header: "Path", accessor: "path", sortable: false, hidden: true },
    {
      header: "Page",
      accessor: "is_page",
      sortable: false,
      hidden: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-muted-foreground"}>
          {v ? "Yes" : "No"}
        </span>
      ),
    },
    {
      header: "Active",
      accessor: "is_active",
      sortable: false,
      hidden: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-red-600 font-medium"}>
          {v ? "Active" : "Inactive"}
        </span>
      ),
    },
    {
      header: "Created",
      accessor: "created_at",
      sortable: true,
      hidden: true,
      formatter: (v: string) => new Date(v).toLocaleString("id-ID"),
    },
  ], []);

  const closeModal = () => {
    setModal({ isOpen: false, mode: "create", title: "" });
    setSelectedRow(null);
    createForm.reset();
    updateForm.reset();
  };

  const openCreate = () => {
    createForm.reset({
      parent_id: "",
      code: "",
      name: "",
      path: "",
      is_page: false,
      is_active: true,
    });
    setModal({ isOpen: true, mode: "create", title: "Create Module" });
  };

  const openUpdate = (row: ModuleListItem) => {
    setSelectedRow(row);
    updateForm.reset({
      parent_id: row.parent_id ?? "",
      code: row.code,
      name: row.name,
      path: row.path,
      is_page: row.is_page,
      is_active: row.is_active,
    });
    setModal({ isOpen: true, mode: "update", title: "Update Module" });
  };

  const openDetail = (row: ModuleListItem) => {
    setSelectedRow(row);
    setModal({ isOpen: true, mode: "detail", title: "Module Details" });
  };

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: {
      parent_id: "",
      code: "",
      name: "",
      path: "",
      is_page: false,
      is_active: true,
    },
  });

  const updateForm = useForm<UpdateForm>({
    resolver: zodResolver(updateSchema),
    defaultValues: {
      parent_id: "",
      code: "",
      name: "",
      path: "",
      is_page: true,
      is_active: true,
    },
  });

  const createParentId = createForm.watch("parent_id");
  const updateParentId = updateForm.watch("parent_id");
  const isCreateHasParent = !!(createParentId && createParentId !== "");
  const isUpdateHasParent = !!(updateParentId && updateParentId !== "");

  const handleCreate = async (values: CreateForm) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    const hasParent = !!values.parent_id;
    try {
      await clientApi("/SM02", {
        method: "POST",
        body: {
          ...values,
          parent_id: values.parent_id || "",
          is_page: hasParent ? values.is_page : false,
        },
      });
      toast.success("Module created successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    }
  };

  const handleUpdate = async (values: UpdateForm) => {
    if (!selectedRow) return;
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    const hasParent = !!values.parent_id;
    try {
      await clientApi("/SM02/" + selectedRow.id, {
        method: "PUT",
        body: {
          ...values,
          parent_id: values.parent_id || "",
          is_page: hasParent ? values.is_page : false,
        },
      });
      toast.success("Module updated successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    }
  };

  const actions: ActionConfig<ModuleListItem> = {
    onSearch: () => setRefresh(!refresh),
    onCreate: openCreate,
    onUpdate: openUpdate,
    onDetail: openDetail,
  };

  const parentModules = allModules.filter(
    (m) => m.id !== selectedRow?.id && !m.parent_id && m.is_active,
  );

  return (
    <Suspense fallback={<p>Loading modules...</p>}>
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
          onSubmit={
            modal.mode === "create"
              ? createForm.handleSubmit(handleCreate)
              : modal.mode === "update"
                ? updateForm.handleSubmit(handleUpdate)
                : undefined
          }
        >
          {modal.mode === "detail" ? (
            <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
              <div className="grid gap-3 sm:grid-cols-2">
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Code</FieldLabel>
                  <Input value={selectedRow?.code ?? ""} disabled readOnly />
                </Field>
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Name</FieldLabel>
                  <Input value={selectedRow?.name ?? ""} disabled readOnly />
                </Field>
              </div>
              <div className="grid gap-3 sm:grid-cols-2">
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Path</FieldLabel>
                  <Input value={selectedRow?.path ?? ""} disabled readOnly />
                </Field>
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Parent</FieldLabel>
                  <Input
                    value={allModules.find((m) => m.id === selectedRow?.parent_id)?.name ?? "Root"}
                    disabled readOnly
                  />
                </Field>
              </div>
              <div className="grid gap-3 sm:grid-cols-2 px-2">
                <Field className="gap-1.5 pt-2">
                  <div className="flex items-center gap-2">
                    <Checkbox checked={selectedRow?.is_page ?? false} disabled />
                    <FieldLabel>Page</FieldLabel>
                  </div>
                </Field>
                <Field className="gap-1.5 pt-2">
                  <div className="flex items-center gap-2">
                    <Checkbox checked={selectedRow?.is_active ?? false} disabled />
                    <FieldLabel>Active</FieldLabel>
                  </div>
                </Field>
              </div>
              <Field className="gap-1.5 mt-[3px] px-2">
                <InputGroup>
                  <InputGroupAddon align="inline-start">
                    <InputGroupText>
                      Created at
                    </InputGroupText>
                  </InputGroupAddon>
                  <InputGroupInput value={selectedRow?.created_at ?? ""} disabled readOnly />
                </InputGroup>
              </Field>
            </FieldGroup>
          ) : modal.mode === "create" ? (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Controller
                  control={createForm.control}
                  name="parent_id"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Parent</FieldLabel>
                      <Select
                        value={field.value || ""}
                        onValueChange={(v) => {
                          const val = v === "__none__" ? "" : v;
                          field.onChange(val);
                          if (!val) createForm.setValue("is_page", false);
                        }}
                        items={[
                          { value: "__none__", label: "Root (No Parent)" },
                          ...parentModules.map((m) => ({
                            value: m.id,
                            label: `${m.code} - ${m.name}`,
                          })),
                        ]}
                      >
                        <SelectTrigger className="w-full">
                          <SelectValue placeholder="Select parent (optional)" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="__none__">Root (No Parent)</SelectItem>
                          {parentModules.map((m) => (
                            <SelectItem key={m.id} value={m.id}>{m.code} - {m.name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </Field>
                  )}
                />
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={createForm.control}
                    name="code"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="code">Code *</FieldLabel>
                        <Input {...field} id="code" placeholder="e.g. SM02"
                         maxLength={50} aria-invalid={fieldState.invalid} value={field.value ?? ""} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="name"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="name">Name *</FieldLabel>
                        <Input {...field} id="name" placeholder="Module Name"
                         maxLength={255} aria-invalid={fieldState.invalid} value={field.value ?? ""} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <Controller
                  control={createForm.control}
                  name="path"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="path">Path *</FieldLabel>
                      <Input {...field} id="path" placeholder="e.g. /board/pages/SM02"
                       maxLength={255} aria-invalid={fieldState.invalid} value={field.value ?? ""} />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <div className="grid gap-3 sm:grid-cols-2 px-2">
                  <Controller
                    control={createForm.control}
                    name="is_page"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox
                            checked={field.value}
                            onCheckedChange={field.onChange}
                            disabled={!isCreateHasParent}
                          />
                          <FieldLabel className={isCreateHasParent
                            ? "cursor-pointer"
                            : "text-muted-foreground"}>Page</FieldLabel>
                        </div>
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="is_active"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                          <FieldLabel className="cursor-pointer">Active</FieldLabel>
                        </div>
                      </Field>
                    )}
                  />
                </div>
              </FieldGroup>
            </FieldSet>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Controller
                  control={updateForm.control}
                  name="parent_id"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Parent</FieldLabel>
                      <Select
                        value={field.value || ""}
                        disabled
                        items={[
                          { value: "__none__", label: "Root (No Parent)" },
                          ...parentModules.map((m) => ({
                            value: m.id,
                            label: `${m.code} - ${m.name}`,
                          })),
                        ]}
                      >
                        <SelectTrigger className="w-full">
                          <SelectValue placeholder="Select parent (optional)" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="__none__">Root (No Parent)</SelectItem>
                          {parentModules.map((m) => (
                            <SelectItem key={m.id} value={m.id}>{m.code} - {m.name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </Field>
                  )}
                />
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={updateForm.control}
                    name="code"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="upd-code">Code *</FieldLabel>
                        <Input {...field} id="upd-code" placeholder="e.g. SM02"
                         maxLength={50} disabled aria-invalid={fieldState.invalid} value={field.value ?? ""} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                  <Controller
                    control={updateForm.control}
                    name="name"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="upd-name">Name *</FieldLabel>
                        <Input {...field} id="upd-name" placeholder="Module Name"
                         maxLength={255} aria-invalid={fieldState.invalid} value={field.value ?? ""} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <Controller
                  control={updateForm.control}
                  name="path"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="upd-path">Path *</FieldLabel>
                      <Input {...field} id="upd-path" placeholder="e.g. /board/pages/SM02"
                       maxLength={255} aria-invalid={fieldState.invalid} value={field.value ?? ""} />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <div className="grid gap-3 sm:grid-cols-2 px-2">
                  <Controller
                    control={updateForm.control}
                    name="is_page"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox
                            checked={field.value}
                            onCheckedChange={field.onChange}
                            disabled={!isUpdateHasParent}
                          />
                          <FieldLabel className={isUpdateHasParent
                            ? "cursor-pointer"
                            : "text-muted-foreground"}>Page</FieldLabel>
                        </div>
                      </Field>
                    )}
                  />
                  <Controller
                    control={updateForm.control}
                    name="is_active"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                          <FieldLabel className="cursor-pointer">Active</FieldLabel>
                        </div>
                      </Field>
                    )}
                  />
                </div>
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>
      </div>
    </Suspense>
  );
}
