"use client";

import { Controller, useFieldArray, useForm } from "react-hook-form";
import { Suspense, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { storageKey, parseSession } from "@/lib/utility";
import { clientApi, ClientApiError } from "@/lib/client-api";
import { toast } from "sonner";
import { getErrorMessage } from "@/lib/error-message";
import { Plus, Trash2, Search } from "lucide-react";
import DataTable, { Column, ActionConfig } from "@/uix/datatable";
import DataDialog from "@/uix/datadialog";
import { Button } from "@/uix/button";
import { Input } from "@/uix/input";
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

interface ApprovalSign {
  user_id: string;
}

interface ApprovalStep {
  id: string;
  step: number;
  condition: string;
  signers: ApprovalSign[];
}

interface DocTypeItem {
  id: string;
  code: string;
  name: string;
  created_at: string;
  steps: ApprovalStep[];
}

interface UserItem {
  id: string;
  username: string;
  fullname: string;
}

const createSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  steps: z.array(z.object({
    step: z.number().min(1),
    condition: z.enum(["ALL_APPROVED", "ANY_APPROVED"]),
    user_ids: z.array(z.string()).min(1, "At least one signer is required"),
  })).min(1, "At least one approval step is required"),
});

type CreateForm = z.infer<typeof createSchema>;

export default function Page() {
  const [selectedRow, setSelectedRow] = useState<DocTypeItem | null>(null);
  const [refresh, setRefresh] = useState(false);
  const [modal, setModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update" | "detail";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });
  const [users, setUsers] = useState<UserItem[]>([]);
  const [signerSearch, setSignerSearch] = useState("");
  const [searchResults, setSearchResults] = useState<UserItem[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const searchTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const fetchUsers = useCallback(async (search: string) => {
    try {
      const res = await clientApi<{ data: UserItem[] }>("/SM04/user", {
        params: search ? { search } : {},
      });
      const list = (res as any)?.data ?? [];
      setUsers(list);
      if (search) setSearchResults(list);
    } catch {
      setUsers([]);
      setSearchResults([]);
    }
  }, []);

  useEffect(() => {
    fetchUsers("");
  }, [fetchUsers]);

  const handleSignerSearch = useCallback((value: string) => {
    setSignerSearch(value);
    if (searchTimerRef.current) clearTimeout(searchTimerRef.current);
    searchTimerRef.current = setTimeout(() => {
      setIsSearching(true);
      fetchUsers(value).finally(() => setIsSearching(false));
    }, 300);
  }, [fetchUsers]);

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
      const data = await clientApi<{ data: DocTypeItem[]; meta: { total: number; page: number; size: number } }>("/SM04", {
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

  const columns: Column<DocTypeItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Code", accessor: "code", sortable: true, hidden: false },
    { header: "Name", accessor: "name", sortable: true, hidden: false },
    {
      header: "Steps",
      accessor: "steps",
      sortable: false,
      hidden: false,
      formatter: (v: ApprovalStep[]) => (
        <span>{v?.length ?? 0}</span>
      ),
    },
    {
      header: "Signers",
      accessor: "steps",
      sortable: false,
      hidden: true,
      formatter: (v: ApprovalStep[]) => {
        const total = v?.reduce((acc, s) => acc + (s.signers?.length ?? 0), 0) ?? 0;
        return <span>{total}</span>;
      },
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
  };

  const openCreate = () => {
    createForm.reset({
      code: "",
      name: "",
      steps: [{ step: 1, condition: "ALL_APPROVED", user_ids: [] }],
    });
    setModal({ isOpen: true, mode: "create", title: "Create Signature Type" });
  };

  const openUpdate = (row: DocTypeItem) => {
    createForm.reset({
      code: row.code,
      name: row.name,
      steps: (row.steps ?? []).map((s) => ({
        step: s.step,
        condition: s.condition as "ALL_APPROVED" | "ANY_APPROVED",
        user_ids: (s.signers ?? []).map((sg) => sg.user_id),
      })),
    });
    setSelectedRow(row);
    setModal({ isOpen: true, mode: "update", title: "Update Signature Type" });
  };

  const openDetail = (row: DocTypeItem) => {
    setSelectedRow(row);
    setModal({ isOpen: true, mode: "detail", title: "Signature Type Details" });
  };

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: {
      code: "",
      name: "",
      steps: [{ step: 1, condition: "ALL_APPROVED", user_ids: [] }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: createForm.control,
    name: "steps",
  });

  const watchedSteps = createForm.watch("steps");
  const allSelectedUserIds = useMemo(() => {
    const set = new Set<string>();
    for (const step of watchedSteps) {
      for (const uid of step.user_ids) {
        set.add(uid);
      }
    }
    return set;
  }, [watchedSteps]);

  const handleCreate = async (values: CreateForm) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM04", {
        method: "POST",
        body: values,
      });
      toast.success("Signature type created successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    }
  };

  const handleUpdate = async (values: CreateForm) => {
    if (!selectedRow) return;
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM04/" + selectedRow.id, {
        method: "PUT",
        body: values,
      });
      toast.success("Signature type updated successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    }
  };

  const actions: ActionConfig<DocTypeItem> = {
    onSearch: () => setRefresh(!refresh),
    onCreate: openCreate,
    onDetail: openDetail,
    onUpdate: openUpdate,
    hideDelete: true,
  };

  const userMap = new Map(users.map((u) => [u.id, u]));

  return (
    <Suspense fallback={<p>Loading signature types...</p>}>
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
                ? createForm.handleSubmit(handleUpdate)
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
              <Field className="gap-1.5 pt-2 px-2">
                <FieldLabel>Created At</FieldLabel>
                <Input value={selectedRow?.created_at ?? ""} disabled readOnly />
              </Field>
              {selectedRow?.steps && selectedRow.steps.length > 0 && (
                <div className="px-2 space-y-3">
                  <FieldLabel>Approval Steps</FieldLabel>
                  {selectedRow.steps.map((step) => (
                    <div key={step.id} className="border rounded-md p-3 space-y-2">
                      <div className="grid gap-3 sm:grid-cols-2">
                        <Field className="gap-1.5">
                          <FieldLabel className="text-xs">Step</FieldLabel>
                          <Input value={String(step.step)} disabled readOnly />
                        </Field>
                        <Field className="gap-1.5">
                          <FieldLabel className="text-xs">Condition</FieldLabel>
                          <Input value={step.condition === "ALL_APPROVED"
                            ? "All Must Approve"
                            : "Any Can Approve"} disabled readOnly />
                        </Field>
                      </div>
                      <Field className="gap-1.5">
                        <FieldLabel className="text-xs">Signers</FieldLabel>
                        <div className="flex flex-wrap gap-1">
                          {step.signers.map((signer) => (
                            <span key={signer.user_id}
                              className="inline-flex items-center rounded-md bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground">
                              {userMap.get(signer.user_id)?.fullname || userMap.get(signer.user_id)?.username || signer.user_id}
                            </span>
                          ))}
                        </div>
                      </Field>
                    </div>
                  ))}
                </div>
              )}
            </FieldGroup>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={createForm.control}
                    name="code"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="code">Code *</FieldLabel>
                        <Input {...field} id="code" placeholder="e.g. PO, LR"
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
                        <Input {...field} id="name" placeholder="e.g. Purchase Order"
                          maxLength={255} aria-invalid={fieldState.invalid} value={field.value ?? ""} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>

                <div className="px-2 space-y-3">
                  <div className="flex items-center justify-between">
                    <FieldLabel>Approval Steps</FieldLabel>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => append({
                        step: fields.length + 1,
                        condition: "ALL_APPROVED",
                        user_ids: []
                      })}
                    >
                      <Plus className="size-3" /> Add Step
                    </Button>
                  </div>

                  {fields.map((field, index) => (
                    <div key={field.id} className="border rounded-md p-3 space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="text-sm font-medium">Step {index + 1}</span>
                        {fields.length > 1 && (
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              remove(index);
                              const currentSteps = createForm.getValues("steps");
                              currentSteps.forEach((_, i) => {
                                createForm.setValue(`steps.${i}.step`, i + 1);
                              });
                            }}
                          >
                            <Trash2 className="size-3" />
                          </Button>
                        )}
                      </div>
                      <div className="grid gap-3 sm:grid-cols-2">
                        <Controller
                          control={createForm.control}
                          name={`steps.${index}.condition`}
                          render={({ field: selField, fieldState }) => (
                            <Field className="gap-1.5" data-invalid={fieldState.invalid}>
                              <FieldLabel className="text-xs">Condition *</FieldLabel>
                              <Select
                                value={selField.value}
                                onValueChange={selField.onChange}
                                items={[
                                  { value: "ALL_APPROVED", label: "All Must Approve" },
                                  { value: "ANY_APPROVED", label: "Any Can Approve" },
                                ]}
                              >
                                <SelectTrigger className="w-full">
                                  <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectItem value="ALL_APPROVED">All Must Approve</SelectItem>
                                  <SelectItem value="ANY_APPROVED">Any Can Approve</SelectItem>
                                </SelectContent>
                              </Select>
                              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                            </Field>
                          )}
                        />
                        <Controller
                          control={createForm.control}
                          name={`steps.${index}.user_ids`}
                          render={({ field: arrField, fieldState }) => {
                            const displayUsers = signerSearch
                              ? searchResults.filter((u) => !allSelectedUserIds.has(u.id))
                              : users.filter((u) => !allSelectedUserIds.has(u.id));
                            return (
                            <Field className="gap-1.5" data-invalid={fieldState.invalid}>
                              <FieldLabel className="text-xs">Signers *</FieldLabel>
                              <div className="relative">
                                <Search className="absolute left-2 top-2.5 size-4 text-muted-foreground" />
                                <Input
                                  placeholder="Search signer by name..."
                                  value={signerSearch}
                                  onChange={(e) => handleSignerSearch(e.target.value)}
                                  className="pl-8"
                                />
                              </div>
                              {displayUsers.length > 0 && (
                                <div className="border rounded-md max-h-40 overflow-y-auto">
                                  {displayUsers.map((u) => (
                                    <button
                                      key={u.id}
                                      type="button"
                                      className="w-full text-left px-3 py-1.5 text-sm hover:bg-muted transition-colors"
                                      onClick={() => {
                                        if (!arrField.value.includes(u.id)) {
                                          arrField.onChange([...arrField.value, u.id]);
                                        }
                                        setSignerSearch("");
                                        setSearchResults([]);
                                      }}
                                    >
                                      {u.fullname || u.username}
                                    </button>
                                  ))}
                                </div>
                              )}
                              {signerSearch && displayUsers.length === 0 && !isSearching && (
                                <p className="text-xs text-muted-foreground px-1">No users found.</p>
                              )}
                              {arrField.value.length > 0 && (
                                <div className="flex flex-wrap gap-1 mt-1">
                                  {arrField.value.map((uid: string) => (
                                    <span key={uid}
                                      className="inline-flex items-center gap-1 rounded-md bg-muted px-2 py-0.5 text-xs font-medium text-muted-foreground">
                                      {userMap.get(uid)?.fullname || userMap.get(uid)?.username || uid}
                                      <button
                                        type="button"
                                        className="ml-0.5 hover:text-destructive"
                                        onClick={() => arrField.onChange(arrField.value.filter((v: string) => v !== uid))}
                                      >
                                        &times;
                                      </button>
                                    </span>
                                  ))}
                                </div>
                              )}
                              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                            </Field>
                          );}}
                        />
                      </div>
                    </div>
                  ))}

                  {createForm.formState.errors.steps?.message && (
                    <p className="text-sm text-destructive px-2">
                      {createForm.formState.errors.steps.message}
                    </p>
                  )}
                </div>
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>
      </div>
    </Suspense>
  );
}
