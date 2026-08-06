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
import { Checkbox } from "@/uix/checkbox";
import { Input } from "@/uix/input";
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
import SearchSelect from "@/uix/search-select";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldSet,
} from "@/uix/field";
import {
  Tabs,
  TabsList,
  TabsTrigger,
  TabsContent,
} from "@/uix/tabs";
import { Alert, AlertTitle } from "@/uix/alert";

interface CompanyListItem {
  id: string;
  code: string;
  name: string;
  vat_id: string;
  reg_no: string;
  address: string;
  valuta: string;
  hris_link: string;
  is_active: boolean;
  created_at: string;
}

interface CompanyModuleItem {
  id: string;
  module_id: string;
  code: string;
  name: string;
  path: string;
  is_active: boolean;
}

interface ModuleOption {
  id: string;
  code: string;
  name: string;
  parent_id: string;
  is_active: boolean;
}

interface AreaItem {
  id: string;
  code: string;
  name: string;
  description: string;
  is_active: boolean;
  created_at: string;
}

const createSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  vat_id: z.string().optional(),
  reg_no: z.string().optional(),
  address: z.string().optional(),
  valuta: z.string().optional(),
  hris_link: z.string().optional(),
  is_active: z.boolean().optional(),
});

const updateSchema = z.object({
  name: z.string().min(1, "Name is required"),
  vat_id: z.string().optional(),
  reg_no: z.string().optional(),
  address: z.string().optional(),
  valuta: z.string().optional(),
  hris_link: z.string().optional(),
  is_active: z.boolean().optional(),
});

const moduleCreateSchema = z.object({
  module_id: z.string().min(1, "Module is required"),
});

const moduleUpdateSchema = z.object({
  is_active: z.boolean(),
});

const areaCreateSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  is_active: z.boolean().optional(),
});

const areaUpdateSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  is_active: z.boolean(),
});

type CreateForm = z.infer<typeof createSchema>;
type UpdateForm = z.infer<typeof updateSchema>;
type ModuleCreateForm = z.infer<typeof moduleCreateSchema>;
type ModuleUpdateForm = z.infer<typeof moduleUpdateSchema>;
type AreaCreateForm = z.infer<typeof areaCreateSchema>;
type AreaUpdateForm = z.infer<typeof areaUpdateSchema>;

const valutaOptions = ["IDR", "USD", "EUR", "SGD", "JPY", "AUD"];

export default function Page() {
  const [selectedRow, setSelectedRow] = useState<CompanyListItem | null>(null);
  const [refresh, setRefresh] = useState(false);
  const [modal, setModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update" | "detail";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });

  const [moduleModal, setModuleModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });
  const [moduleRefresh, setModuleRefresh] = useState(false);
  const [allModules, setAllModules] = useState<ModuleOption[]>([]);
  const [assignedModuleIds, setAssignedModuleIds] = useState<Set<string>>(new Set());
  const [selectedModule, setSelectedModule] = useState<CompanyModuleItem | null>(null);

  const [areaModal, setAreaModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });
  const [areaRefresh, setAreaRefresh] = useState(false);
  const [selectedArea, setSelectedArea] = useState<AreaItem | null>(null);
  const [areaTotal, setAreaTotal] = useState(0);

  const [companyLoading, setCompanyLoading] = useState(false);
  const [moduleLoading, setModuleLoading] = useState(false);
  const [areaLoading, setAreaLoading] = useState(false);

  useEffect(() => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    clientApi<ModuleOption[]>("/SM03/module")
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
        data: CompanyListItem[];
        meta: { total: number; page: number; size: number }
      }>("/SM03", {
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

  const columns: Column<CompanyListItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Code", accessor: "code", sortable: true, hidden: false },
    { header: "Name", accessor: "name", sortable: true, hidden: false },
    { header: "VAT ID", accessor: "vat_id", sortable: false, hidden: true },
    { header: "Reg No", accessor: "reg_no", sortable: false, hidden: true },
    { header: "Valuta", accessor: "valuta", sortable: false, hidden: false },
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
      code: "",
      name: "",
      vat_id: "",
      reg_no: "",
      address: "",
      valuta: "IDR",
      hris_link: "",
      is_active: true,
    });
    setModal({ isOpen: true, mode: "create", title: "Create Company" });
  };

  const openUpdate = (row: CompanyListItem) => {
    setSelectedRow(row);
    updateForm.reset({
      name: row.name,
      vat_id: row.vat_id ?? "",
      reg_no: row.reg_no ?? "",
      address: row.address ?? "",
      valuta: row.valuta ?? "IDR",
      hris_link: row.hris_link ?? "",
      is_active: row.is_active,
    });
    setModal({ isOpen: true, mode: "update", title: "Update Company" });
  };

  const openDetail = (row: CompanyListItem) => {
    setSelectedRow(row);
    setModal({ isOpen: true, mode: "detail", title: "Company Details" });
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (session) {
      clientApi("/SM03/" + row.id + "/module/select")
        .then((res: any) => {
          setAssignedModuleIds(new Set((res?.data ?? []).map((m: CompanyModuleItem) => m.module_id)));
        })
        .catch(() => {
          setAssignedModuleIds(new Set());
        });
    }
  };

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: {
      code: "",
      name: "",
      vat_id: "",
      reg_no: "",
      address: "",
      valuta: "IDR",
      hris_link: "",
      is_active: true,
    },
  });

  const updateForm = useForm<UpdateForm>({
    resolver: zodResolver(updateSchema),
    defaultValues: {
      name: "",
      vat_id: "",
      reg_no: "",
      address: "",
      valuta: "IDR",
      hris_link: "",
      is_active: true,
    },
  });

  const handleCreate = async (values: CreateForm) => {
    setCompanyLoading(true);
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM03", {
        method: "POST",
        body: {
          ...values,
          vat_id: values.vat_id || "",
          reg_no: values.reg_no || "",
          address: values.address || "",
          valuta: values.valuta || "IDR",
          hris_link: values.hris_link || "",
        },
      });
      toast.success("Company created successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    } finally {
      setCompanyLoading(false);
    }
  };

  const handleUpdate = async (values: UpdateForm) => {
    setCompanyLoading(true);
    if (!selectedRow) return;
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM03/" + selectedRow.id, {
        method: "PUT",
        body: {
          ...values,
          vat_id: values.vat_id || "",
          reg_no: values.reg_no || "",
          address: values.address || "",
          valuta: values.valuta || "IDR",
          hris_link: values.hris_link || "",
        },
      });
      toast.success("Company updated successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    } finally {
      setCompanyLoading(false);
    }
  };

  const actions: ActionConfig<CompanyListItem> = {
    onSearch: () => setRefresh(!refresh),
    onCreate: openCreate,
    onUpdate: openUpdate,
    onDetail: openDetail,
  };

  const loadModuleData = useCallback(async (params: {
    search: string;
    page: number;
    size: number;
    sort_by: string;
    sort_order: "asc" | "desc";
  }) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session || !selectedRow) {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
    try {
      const res = await clientApi<{
        data: CompanyModuleItem[];
        meta: { total: number; page: number; size: number };
      }>("/SM03/" + selectedRow.id + "/module", {
        params: {
          search: params.search,
          page: String(params.page),
          size: String(params.size),
          sort_by: params.sort_by,
          sort_order: params.sort_order,
        },
      });
      const data = (res as any)?.data ?? [];
      setAssignedModuleIds(new Set(data.map((m: CompanyModuleItem) => m.module_id)));
      return {
        data,
        meta: (res as any)?.meta ?? { total: 0, page: 1, size: 10 },
      };
    } catch {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
  }, [selectedRow]);

  const moduleColumns: Column<CompanyModuleItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Module ID", accessor: "module_id", sortable: false, hidden: true },
    { header: "Code", accessor: "code", sortable: true, hidden: false },
    { header: "Name", accessor: "name", sortable: true, hidden: false },
    { header: "Path", accessor: "path", sortable: false, hidden: true },
    {
      header: "Active",
      accessor: "is_active",
      sortable: false,
      hidden: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-red-600 font-medium"}>
          {v ? "Yes" : "No"}
        </span>
      ),
    },
  ], []);

  const areaColumns: Column<AreaItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Code", accessor: "code", sortable: true, hidden: false },
    { header: "Name", accessor: "name", sortable: true, hidden: false },
    { header: "Description", accessor: "description", sortable: false, hidden: true },
    {
      header: "Active",
      accessor: "is_active",
      sortable: false,
      hidden: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-red-600 font-medium"}>
          {v ? "Yes" : "No"}
        </span>
      ),
    },
    {
      header: "Created",
      accessor: "created_at",
      sortable: false,
      hidden: true,
      formatter: (v: string) => new Date(v).toLocaleString("id-ID"),
    },
  ], []);

  const moduleCreateForm = useForm<ModuleCreateForm>({
    resolver: zodResolver(moduleCreateSchema),
    defaultValues: { module_id: "" },
  });

  const moduleUpdateForm = useForm<ModuleUpdateForm>({
    resolver: zodResolver(moduleUpdateSchema),
    defaultValues: { is_active: true },
  });

  const openModuleCreate = async () => {
    moduleCreateForm.reset({ module_id: "" });
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (session && selectedRow) {
      try {
        const res = await clientApi<{
          data: CompanyModuleItem[]
        }>("/SM03/" + selectedRow.id + "/module/select");
        setAssignedModuleIds(new Set((res as any)?.data?.map((m: CompanyModuleItem) => m.module_id) ?? []));
      } catch {
        setAssignedModuleIds(new Set());
      }
    }
    setModuleModal({ isOpen: true, mode: "create", title: "Assign Module" });
  };

  const openModuleUpdate = (row: CompanyModuleItem) => {
    setSelectedModule(row);
    moduleUpdateForm.reset({ is_active: row.is_active });
    setModuleModal({ isOpen: true, mode: "update", title: "Update Module" });
  };

  const handleModuleCreate = async (values: ModuleCreateForm) => {
    setModuleLoading(true);
    if (!selectedRow) return;
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM03/" + selectedRow.id + "/module", {
        method: "POST",
        body: { module_id: values.module_id },
      });
      toast.success("Module assigned successfully.");
      setModuleModal({ isOpen: false, mode: "create", title: "" });
      setModuleRefresh(!moduleRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    } finally {
      setModuleLoading(false);
    }
  };

  const handleModuleUpdate = async (values: ModuleUpdateForm) => {
    setModuleLoading(true);
    if (!selectedRow || !selectedModule) return;
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM03/" + selectedRow.id + "/module/" + selectedModule.module_id, {
        method: "PUT",
        body: { is_active: values.is_active },
      });
      toast.success("Module updated successfully.");
      setModuleModal({ isOpen: false, mode: "update", title: "" });
      setModuleRefresh(!moduleRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    } finally {
      setModuleLoading(false);
    }
  };

  const moduleActions: ActionConfig<CompanyModuleItem> = {
    onSearch: () => setModuleRefresh(!moduleRefresh),
    onCreate: openModuleCreate,
    onUpdate: openModuleUpdate,
    hideDetail: true,
  };

  const availableModules = allModules.filter(
    (m) => !assignedModuleIds.has(m.id) && m.parent_id !== "" && m.is_active,
  );

  const areaCreateForm = useForm<AreaCreateForm>({
    resolver: zodResolver(areaCreateSchema),
    defaultValues: { code: "", name: "", description: "", is_active: true },
  });

  const areaUpdateForm = useForm<AreaUpdateForm>({
    resolver: zodResolver(areaUpdateSchema),
    defaultValues: { code: "", name: "", description: "", is_active: true },
  });

  const loadAreaData = useCallback(async (params: {
    search: string;
    page: number;
    size: number;
    sort_by: string;
    sort_order: "asc" | "desc";
  }) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session || !selectedRow) {
      setAreaTotal(0);
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
    try {
      const res = await clientApi<{
        data: AreaItem[];
        meta: { total: number; page: number; size: number };
      }>("/SM03/" + selectedRow.id + "/area", {
        params: {
          search: params.search,
          page: String(params.page),
          size: String(params.size),
          sort_by: params.sort_by,
          sort_order: params.sort_order,
        },
      });
      const data = (res as any)?.data ?? [];
      const meta = (res as any)?.meta ?? { total: 0, page: 1, size: 10 };
      setAreaTotal(meta.total);
      return {
        data,
        meta,
      };
    } catch {
      setAreaTotal(0);
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
  }, [selectedRow]);

  const openAreaCreate = () => {
    areaCreateForm.reset({ code: "", name: "", description: "", is_active: true });
    setAreaModal({ isOpen: true, mode: "create", title: "Create Area" });
  };

  const openAreaUpdate = (row: AreaItem) => {
    setSelectedArea(row);
    areaUpdateForm.reset({
      code: row.code,
      name: row.name,
      description: row.description ?? "",
      is_active: row.is_active,
    });
    setAreaModal({ isOpen: true, mode: "update", title: "Update Area" });
  };

  const handleAreaCreate = async (values: AreaCreateForm) => {
    setAreaLoading(true);
    if (!selectedRow) return;
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM03/" + selectedRow.id + "/area", {
        method: "POST",
        body: {
          ...values,
          description: values.description || "",
          is_active: values.is_active ?? true,
        },
      });
      toast.success("Area created successfully.");
      setAreaModal({ isOpen: false, mode: "create", title: "" });
      setAreaRefresh(!areaRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    } finally {
      setAreaLoading(false);
    }
  };

  const handleAreaUpdate = async (values: AreaUpdateForm) => {
    setAreaLoading(true);
    if (!selectedArea) return;
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM03/" + selectedRow?.id + "/area/" + selectedArea.id, {
        method: "PUT",
        body: {
          ...values,
          description: values.description || "",
          is_active: values.is_active ?? true,
        },
      });
      toast.success("Area updated successfully.");
      setAreaModal({ isOpen: false, mode: "update", title: "" });
      setAreaRefresh(!areaRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    } finally {
      setAreaLoading(false);
    }
  };

  const handleAreaDelete = async (row: AreaItem) => {
    if (!selectedRow) return;
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) return;
    try {
      await clientApi("/SM03/" + selectedRow.id + "/area/" + row.id, {
        method: "DELETE",
      });
      toast.success("Area deleted successfully.");
      setAreaRefresh(!areaRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status === 401) return;
      toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
    }
  };

  const areaActions: ActionConfig<AreaItem> = {
    onSearch: () => setAreaRefresh(!areaRefresh),
    onCreate: openAreaCreate,
    onUpdate: openAreaUpdate,
    onDelete: handleAreaDelete,
    hideDetail: true,
  };

  return (
    <Suspense fallback={<p>Loading companies...</p>}>
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
          loading={companyLoading}
          onSubmit={
            modal.mode === "create"
              ? createForm.handleSubmit(handleCreate)
              : modal.mode === "update"
                ? updateForm.handleSubmit(handleUpdate)
                : undefined
          }
        >
          {modal.mode === "detail" ? (
            <FieldGroup className="p-3">
              <Tabs defaultValue="information">
                <TabsList variant="line" className="mb-3">
                  <TabsTrigger value="information">Information</TabsTrigger>
                  <TabsTrigger value="modules">Modules</TabsTrigger>
                  <TabsTrigger value="area">Area</TabsTrigger>
                </TabsList>
                <TabsContent value="information">
                  <FieldGroup className="gap-3">
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
                        <FieldLabel>VAT ID</FieldLabel>
                        <Input value={selectedRow?.vat_id ?? ""} disabled readOnly />
                      </Field>
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Reg No</FieldLabel>
                        <Input value={selectedRow?.reg_no ?? ""} disabled readOnly />
                      </Field>
                    </div>
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Address</FieldLabel>
                      <Input value={selectedRow?.address ?? ""} disabled readOnly />
                    </Field>
                    <div className="grid gap-3 sm:grid-cols-2">
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Valuta</FieldLabel>
                        <Input value={selectedRow?.valuta ?? ""} disabled readOnly />
                      </Field>
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>HRIS Link</FieldLabel>
                        <Input value={selectedRow?.hris_link ?? ""} disabled readOnly />
                      </Field>
                    </div>
                    <Field className="gap-1.5 pt-2 px-2">
                      <div className="flex items-center gap-2">
                        <Checkbox checked={selectedRow?.is_active ?? false} disabled />
                        <FieldLabel>Active</FieldLabel>
                      </div>
                    </Field>
                    <Field className="gap-1.5 mt-[3px] px-2">
                      <InputGroup>
                        <InputGroupAddon align="inline-start">
                          <InputGroupText>Created at</InputGroupText>
                        </InputGroupAddon>
                        <InputGroupInput value={selectedRow?.created_at ?? ""} disabled readOnly />
                      </InputGroup>
                    </Field>
                  </FieldGroup>
                </TabsContent>
                <TabsContent value="modules">
                  <div className="space-y-3">
                    {assignedModuleIds.size > 0 && (
                      <Alert variant="info" className="mb-3">
                        <AlertTitle>{assignedModuleIds.size} module(s) registered to this company.</AlertTitle>
                      </Alert>
                    )}
                    <DataTable
                      fetchData={loadModuleData}
                      columns={moduleColumns}
                      actions={moduleActions}
                      hideSearch={false}
                      hideSelect={true}
                      hidePaging={false}
                      hideSort={false}
                      hideColumnToggle={true}
                      refreshTrigger={moduleRefresh}
                    />
                  </div>
                </TabsContent>
                <TabsContent value="area">
                  <div className="space-y-3">
                    {areaTotal > 0 && (
                      <Alert variant="info" className="mb-3">
                        <AlertTitle>{areaTotal} area(s) registered to this company.</AlertTitle>
                      </Alert>
                    )}
                    <DataTable
                      fetchData={loadAreaData}
                      columns={areaColumns}
                      actions={areaActions}
                      hideSearch={false}
                      hideSelect={true}
                      hidePaging={false}
                      hideSort={false}
                      hideColumnToggle={true}
                      refreshTrigger={areaRefresh}
                    />
                  </div>
                </TabsContent>
              </Tabs>
            </FieldGroup>
          ) : modal.mode === "create" ? (
            <FieldSet>
              <FieldGroup className="gap-3 p-3">
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={createForm.control}
                    name="code"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="code">Code *</FieldLabel>
                        <Input {...field} id="code" placeholder="e.g. CP01"
                         maxLength={50} aria-invalid={fieldState.invalid} />
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
                        <Input {...field} id="name" placeholder="Company Name"
                         maxLength={255} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={createForm.control}
                    name="vat_id"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="vat_id">VAT ID</FieldLabel>
                        <Input {...field} id="vat_id" placeholder="Optional" maxLength={100} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="reg_no"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="reg_no">Reg No</FieldLabel>
                        <Input {...field} id="reg_no" placeholder="Optional" maxLength={100} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                </div>
                <Controller
                  control={createForm.control}
                  name="address"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel htmlFor="address">Address</FieldLabel>
                      <Input {...field} id="address" placeholder="Optional" maxLength={500} value={field.value ?? ""} />
                    </Field>
                  )}
                />
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={createForm.control}
                    name="valuta"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Valuta</FieldLabel>
                        <Select
                          value={field.value || "IDR"}
                          onValueChange={field.onChange}
                          items={valutaOptions.map((v) => ({ value: v, label: v }))}
                        >
                          <SelectTrigger className="w-full">
                            <SelectValue placeholder="Select valuta" />
                          </SelectTrigger>
                          <SelectContent>
                            {valutaOptions.map((v) => (
                              <SelectItem key={v} value={v}>{v}</SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="hris_link"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="hris_link">HRIS Link</FieldLabel>
                        <Input {...field} id="hris_link" placeholder="Optional" maxLength={500} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                </div>
                <Field className="gap-1.5 pt-2 px-2">
                  <div className="flex items-center gap-2">
                    <Controller
                      control={createForm.control}
                      name="is_active"
                      render={({ field }) => (
                        <Checkbox
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      )}
                    />
                    <FieldLabel className="cursor-pointer">Active</FieldLabel>
                  </div>
                </Field>
              </FieldGroup>
            </FieldSet>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3">
                <div className="grid gap-3 sm:grid-cols-2">
                  <Field className="gap-1.5 pt-2 px-2">
                    <FieldLabel>Code</FieldLabel>
                    <Input value={selectedRow?.code ?? ""} disabled readOnly />
                  </Field>
                  <Controller
                    control={updateForm.control}
                    name="name"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="upd-name">Name *</FieldLabel>
                        <Input {...field} id="upd-name" placeholder="Company Name"
                         maxLength={255} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={updateForm.control}
                    name="vat_id"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="upd-vat_id">VAT ID</FieldLabel>
                        <Input {...field} id="upd-vat_id" placeholder="Optional" maxLength={100} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                  <Controller
                    control={updateForm.control}
                    name="reg_no"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="upd-reg_no">Reg No</FieldLabel>
                        <Input {...field} id="upd-reg_no" placeholder="Optional" maxLength={100} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                </div>
                <Controller
                  control={updateForm.control}
                  name="address"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel htmlFor="upd-address">Address</FieldLabel>
                      <Input {...field} id="upd-address" placeholder="Optional" maxLength={500} value={field.value ?? ""} />
                    </Field>
                  )}
                />
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={updateForm.control}
                    name="valuta"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Valuta</FieldLabel>
                        <Select
                          value={field.value || "IDR"}
                          onValueChange={field.onChange}
                          items={valutaOptions.map((v) => ({ value: v, label: v }))}
                        >
                          <SelectTrigger className="w-full">
                            <SelectValue placeholder="Select valuta" />
                          </SelectTrigger>
                          <SelectContent>
                            {valutaOptions.map((v) => (
                              <SelectItem key={v} value={v}>{v}</SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </Field>
                    )}
                  />
                  <Controller
                    control={updateForm.control}
                    name="hris_link"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="upd-hris_link">HRIS Link</FieldLabel>
                        <Input {...field} id="upd-hris_link" placeholder="Optional" maxLength={500} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                </div>
                <Field className="gap-1.5 pt-2 px-2">
                  <div className="flex items-center gap-2">
                    <Controller
                      control={updateForm.control}
                      name="is_active"
                      render={({ field }) => (
                        <Checkbox
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      )}
                    />
                    <FieldLabel className="cursor-pointer">Active</FieldLabel>
                  </div>
                </Field>
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>

        <DataDialog
          isOpen={moduleModal.isOpen}
          mode={moduleModal.mode}
          title={moduleModal.title}
          onClose={() => {
            setModuleModal({ isOpen: false, mode: "create", title: "" });
            setSelectedModule(null);
          }}
          loading={moduleLoading}
          onSubmit={
            moduleModal.mode === "create"
              ? moduleCreateForm.handleSubmit(handleModuleCreate)
              : moduleUpdateForm.handleSubmit(handleModuleUpdate)
          }
        >
          {moduleModal.mode === "create" ? (
            <FieldSet>
              <FieldGroup className="gap-3 p-3">
                <Controller
                  control={moduleCreateForm.control}
                  name="module_id"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Module *</FieldLabel>
                      <SearchSelect
                        items={availableModules.map((m) => ({
                          id: m.id,
                          name: `${m.code} - ${m.name}`,
                        }))}
                        value={field.value ?? ""}
                        onValueChange={field.onChange}
                        placeholder={availableModules.length === 0 ? "No available modules" : "Select module..."}
                        disabled={availableModules.length === 0}
                        error={fieldState.invalid}
                      />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
              </FieldGroup>
            </FieldSet>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3">
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Module</FieldLabel>
                  <Input
                    value={selectedModule ? `${selectedModule.code} - ${selectedModule.name}` : ""}
                    disabled readOnly
                  />
                </Field>
                <Controller
                  control={moduleUpdateForm.control}
                  name="is_active"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
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
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>

        <DataDialog
          isOpen={areaModal.isOpen}
          mode={areaModal.mode}
          title={areaModal.title}
          onClose={() => {
            setAreaModal({ isOpen: false, mode: "create", title: "" });
            setSelectedArea(null);
          }}
          loading={areaLoading}
          onSubmit={
            areaModal.mode === "create"
              ? areaCreateForm.handleSubmit(handleAreaCreate)
              : areaUpdateForm.handleSubmit(handleAreaUpdate)
          }
        >
          {areaModal.mode === "create" ? (
            <FieldSet>
              <FieldGroup className="gap-3 p-3">
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={areaCreateForm.control}
                    name="code"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="type-code">Code *</FieldLabel>
                        <Input {...field} id="type-code" placeholder="e.g. A1" maxLength={50} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                  <Controller
                    control={areaCreateForm.control}
                    name="name"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="type-name">Name *</FieldLabel>
                        <Input {...field} id="type-name" placeholder="Type Name" maxLength={255} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <Controller
                  control={areaCreateForm.control}
                  name="description"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel htmlFor="type-desc">Description</FieldLabel>
                      <Input {...field} id="type-desc" placeholder="Optional" maxLength={500} value={field.value ?? ""} />
                    </Field>
                  )}
                />
                <Field className="gap-1.5 pt-2 px-2">
                  <div className="flex items-center gap-2">
                    <Controller
                      control={areaCreateForm.control}
                      name="is_active"
                      render={({ field }) => (
                        <Checkbox checked={field.value ?? true} onCheckedChange={field.onChange} />
                      )}
                    />
                    <FieldLabel className="cursor-pointer">Active</FieldLabel>
                  </div>
                </Field>
              </FieldGroup>
            </FieldSet>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3">
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={areaUpdateForm.control}
                    name="code"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="upd-type-code">Code *</FieldLabel>
                        <Input {...field} id="upd-type-code" placeholder="e.g. A1" maxLength={50} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                  <Controller
                    control={areaUpdateForm.control}
                    name="name"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="upd-type-name">Name *</FieldLabel>
                        <Input {...field} id="upd-type-name" placeholder="Type Name" maxLength={255} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <Controller
                  control={areaUpdateForm.control}
                  name="description"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel htmlFor="upd-type-desc">Description</FieldLabel>
                      <Input {...field} id="upd-type-desc" placeholder="Optional" maxLength={500} value={field.value ?? ""} />
                    </Field>
                  )}
                />
                <Field className="gap-1.5 pt-2 px-2">
                  <div className="flex items-center gap-2">
                    <Controller
                      control={areaUpdateForm.control}
                      name="is_active"
                      render={({ field }) => (
                        <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                      )}
                    />
                    <FieldLabel className="cursor-pointer">Active</FieldLabel>
                  </div>
                </Field>
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>
      </div>
    </Suspense>
  );
}
