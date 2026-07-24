"use client";

import { Controller, useForm, useWatch } from "react-hook-form";
import { Suspense, useCallback, useEffect, useMemo, useState } from "react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { storageKey, parseSession } from "@/lib/utility";
import { clientApi, ClientApiError } from "@/lib/client-api";
import { toast } from "sonner";
import { getErrorMessage } from "@/lib/error-message";
import { CircleCheckIcon } from "lucide-react";
import DataTable, { Column, ActionConfig } from "@/uix/datatable";
import DataDialog from "@/uix/datadialog";
import { Alert, AlertTitle } from "@/uix/alert";
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
import { Tabs,
  TabsList,
  TabsTrigger,
  TabsContent
} from "@/uix/tabs";
import { Checkbox } from "@/uix/checkbox";
import { Input } from "@/uix/input";

interface UserListItem {
  id: string;
  username: string;
  email: string;
  fullname: string;
  phone: string;
  company_id: string;
  role: string;
  job: string;
  is_admin: boolean;
  is_hris: boolean;
  is_active: boolean;
  created_at: string;
}

interface CompanyItem {
  id: string;
  code: string;
  name: string;
  hris_link: string;
}

interface UserCompanyListItem {
  id: string;
  user_id: string;
  company_id: string;
  company_name: string;
  is_active: boolean;
  created_at: string;
}

interface UserPrivilegeListItem {
  id: string;
  user_company_id: string;
  module_id: string;
  company_name: string;
  module_code: string;
  module_name: string;
  level: string;
  created_at: string;
}

interface CompanySelectItem {
  id: string;
  name: string;
}

interface ModuleSelectItem {
  id: string;
  parent_id: string;
  code: string;
  name: string;
  is_active: boolean;
}

interface UserLocationListItem {
  id: string;
  user_id: string;
  location_type_id: string;
  company_name: string;
  type_code: string;
  type_name: string;
  is_active: boolean;
  created_at: string;
}

interface LocationTypeSelectItem {
  id: string;
  code: string;
  name: string;
}

const createSchema = z.object({
  username: z.string().min(1, "Username is required"),
  email: z.string().min(1, "Email is required").email("Invalid email address"),
  password: z
    .string()
    .min(6, "Password must be at least 6 characters")
    .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
    .regex(/[a-z]/, "Password must contain at least one lowercase letter")
    .regex(/[0-9]/, "Password must contain at least one number"),
  fullname: z.string().min(1, "Fullname is required"),
  company_id: z.string().optional(),
  phone: z.string().optional(),
  role: z.string().optional(),
  job: z.string().optional(),
  is_admin: z.boolean().optional(),
  is_hris: z.boolean().optional(),
  is_active: z.boolean().optional(),
});

const updateSchema = z.object({
  email: z.string().min(1, "Email is required").email("Invalid email address"),
  fullname: z.string().min(1, "Fullname is required"),
  password: z.string().optional(),
  phone: z.string().optional(),
  role: z.string().optional(),
  job: z.string().optional(),
  is_admin: z.boolean().optional(),
  is_hris: z.boolean().optional(),
  is_active: z.boolean().optional(),
});

const companyCreateSchema = z.object({
  company_id: z.string().min(1, "Company is required"),
  is_active: z.boolean().optional(),
});

const companyUpdateSchema = z.object({
  is_active: z.boolean(),
});

const privilegeCreateSchema = z.object({
  user_company_id: z.string().min(1, "Company assignment is required"),
  module_id: z.string().min(1, "Module is required"),
  level: z.string().min(1, "Level is required"),
});

const privilegeUpdateSchema = z.object({
  level: z.string().min(1, "Level is required"),
});

const locationCreateSchema = z.object({
  company_id: z.string().min(1, "Company is required"),
  location_type_id: z.string().min(1, "Location type is required"),
  is_active: z.boolean().optional(),
});

const locationUpdateSchema = z.object({
  is_active: z.boolean(),
});

type CreateForm = z.infer<typeof createSchema>;
type UpdateForm = z.infer<typeof updateSchema>;
type CompanyCreateForm = z.infer<typeof companyCreateSchema>;
type CompanyUpdateForm = z.infer<typeof companyUpdateSchema>;
type PrivilegeCreateForm = z.infer<typeof privilegeCreateSchema>;
type PrivilegeUpdateForm = z.infer<typeof privilegeUpdateSchema>;
type LocationCreateForm = z.infer<typeof locationCreateSchema>;
type LocationUpdateForm = z.infer<typeof locationUpdateSchema>;

const levelOptions = [
  { value: "hide", label: "Hide" },
  { value: "view", label: "View" },
  { value: "book", label: "Book" },
  { value: "post", label: "Post" },
];

export default function Page() {
  const [selectedRow, setSelectedRow] = useState<UserListItem | null>(null);
  const [refresh, setRefresh] = useState(false);
  const [modal, setModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update" | "detail";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });
  const [companies, setCompanies] = useState<CompanyItem[]>([]);

  const [companyModal, setCompanyModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });
  const [companyRefresh, setCompanyRefresh] = useState(false);
  const [selectedUserCompany, setSelectedUserCompany] = useState<UserCompanyListItem | null>(null);
  const [allCompanies, setAllCompanies] = useState<CompanySelectItem[]>([]);
  const [assignedCompanyIds, setAssignedCompanyIds] = useState<Set<string>>(new Set());

  const [privilegeModal, setPrivilegeModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });
  const [privilegeRefresh, setPrivilegeRefresh] = useState(false);
  const [selectedUserPrivilege, setSelectedUserPrivilege] = useState<UserPrivilegeListItem | null>(null);
  const [allModules, setAllModules] = useState<ModuleSelectItem[]>([]);
  const [userCompanyList, setUserCompanyList] = useState<UserCompanyListItem[]>([]);

  const [locationModal, setLocationModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update";
    title: string;
  }>({ isOpen: false, mode: "create", title: "" });
  const [locationRefresh, setLocationRefresh] = useState(false);
  const [selectedUserLocation, setSelectedUserLocation] = useState<UserLocationListItem | null>(null);
  const [allLocationTypes, setAllLocationTypes] = useState<LocationTypeSelectItem[]>([]);

  useEffect(() => {
    clientApi<{ data: CompanyItem[] }>("/SM01/hris")
      .then((data) => setCompanies(data?.data ?? []))
      .catch(() => setCompanies([]));
    clientApi<{ data: CompanySelectItem[] }>("/SM01/company")
      .then((data) => setAllCompanies(data?.data ?? []))
      .catch(() => setAllCompanies([]));
    clientApi<{ data: ModuleSelectItem[] }>("/SM01/module")
      .then((data) => setAllModules(data?.data ?? []))
      .catch(() => setAllModules([]));
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
        data: UserListItem[];
        meta: { total: number; page: number; size: number }
      }>("/SM01", {
        params: params as unknown as Record<string, string>
      });
      return {
        data: data?.data ?? [],
        meta: data?.meta ?? { total: 0, page: 1, size: 10 },
      };
    } catch (err) {
      if (err instanceof ClientApiError) {
        return { data: [], meta: { total: 0, page: 1, size: 10 } };
      }
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
  }, []);

  const columns: Column<UserListItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Username", accessor: "username", sortable: true, hidden: false },
    { header: "Fullname", accessor: "fullname", sortable: true, hidden: false },
    { header: "Email", accessor: "email", sortable: true, hidden: true },
    { header: "Phone", accessor: "phone", sortable: false, hidden: true },
    { header: "Role", accessor: "role", sortable: false, hidden: true },
    { header: "Job", accessor: "job", sortable: false, hidden: true },
    {
      header: "Admin",
      accessor: "is_admin",
      sortable: false,
      hidden: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-muted-foreground"}>
          {v ? "Yes" : "No"}
        </span>
      ),
    },
    {
      header: "HRIS",
      accessor: "is_hris",
      sortable: false,
      hidden: true,
      formatter: (v: boolean) => (v ? "Yes" : "No"),
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
      username: "",
      email: "",
      password: "",
      fullname: "",
      company_id: "",
      phone: "",
      role: "staff",
      job: "",
      is_admin: false,
      is_hris: false,
      is_active: true,
    });
    setModal({ isOpen: true, mode: "create", title: "Create User" });
  };

  const openUpdate = (row: UserListItem) => {
    setSelectedRow(row);
    updateForm.reset({
      email: row.email,
      fullname: row.fullname,
      password: "",
      phone: row.phone ?? "",
      role: row.role ?? "staff",
      job: row.job ?? "",
      is_admin: row.is_admin,
      is_hris: row.is_hris,
      is_active: row.is_active,
    });
    setModal({ isOpen: true, mode: "update", title: "Update User" });
  };

  const openDetail = (row: UserListItem) => {
    setSelectedRow(row);
    setModal({ isOpen: true, mode: "detail", title: "User Details" });
    clientApi<{ data: UserCompanyListItem[] }>(`/SM01/${row.id}/company`)
      .then((data) => setUserCompanyList(data?.data ?? []))
      .catch(() => setUserCompanyList([]));
  };

  const createForm = useForm<CreateForm>({
    resolver: zodResolver(createSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
      fullname: "",
      company_id: "",
      phone: "",
      role: "staff",
      job: "",
      is_admin: false,
      is_hris: false,
      is_active: true,
    },
  });

  const updateForm = useForm<UpdateForm>({
    resolver: zodResolver(updateSchema),
    defaultValues: {
      email: "",
      fullname: "",
      password: "",
      phone: "",
      role: "staff",
      job: "",
      is_admin: false,
      is_hris: false,
      is_active: true,
    },
  });

  const watchedCompanyId = useWatch({ control: createForm.control, name: "company_id" });

  useEffect(() => {
    const hasCompany = !!watchedCompanyId && watchedCompanyId !== "__none__";
    createForm.setValue("is_hris", hasCompany);
  }, [watchedCompanyId, createForm]);

  const handleCreate = async (values: CreateForm) => {
    try {
      await clientApi("/SM01", { method: "POST", body: values });
      toast.success("User created successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const handleUpdate = async (values: UpdateForm) => {
    if (!selectedRow) return;
    try {
      await clientApi(`/SM01/${selectedRow.id}`, { method: "PUT", body: values });
      toast.success("User updated successfully.");
      closeModal();
      setRefresh(!refresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const actions: ActionConfig<UserListItem> = {
    onSearch: () => setRefresh(!refresh),
    onCreate: openCreate,
    onUpdate: openUpdate,
    onDetail: openDetail,
    hideDelete: true,
  };

  const companyColumns: Column<UserCompanyListItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Company ID", accessor: "company_id", sortable: false, hidden: true },
    { header: "Company", accessor: "company_name", sortable: true, hidden: false },
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

  const companyCreateForm = useForm<CompanyCreateForm>({
    resolver: zodResolver(companyCreateSchema),
    defaultValues: { company_id: "", is_active: true },
  });

  const companyUpdateForm = useForm<CompanyUpdateForm>({
    resolver: zodResolver(companyUpdateSchema),
    defaultValues: { is_active: true },
  });

  const loadCompanyData = useCallback(async (params: {
    search: string; page: number; size: number; sort_by: string; sort_order: "asc" | "desc";
  }) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session || !selectedRow) {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
    try {
      const data = await clientApi<{
        data: UserCompanyListItem[];
        meta: { total: number; page: number; size: number };
      }>(`/SM01/${selectedRow.id}/company`, {
        params: {
          search: params.search,
          page: String(params.page),
          size: String(params.size),
          sort_by: params.sort_by,
          sort_order: params.sort_order,
        },
      });
      const list = data?.data ?? [];
      setAssignedCompanyIds(new Set(list.map((c) => c.company_id)));
      return {
        data: list,
        meta: data?.meta ?? { total: 0, page: 1, size: 10 },
      };
    } catch (err) {
      if (err instanceof ClientApiError) {
        return { data: [], meta: { total: 0, page: 1, size: 10 } };
      }
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
  }, [selectedRow]);

  const openCompanyCreate = () => {
    companyCreateForm.reset({ company_id: "", is_active: true });
    setCompanyModal({ isOpen: true, mode: "create", title: "Assign Company" });
  };

  const openCompanyUpdate = (row: UserCompanyListItem) => {
    setSelectedUserCompany(row);
    companyUpdateForm.reset({ is_active: row.is_active });
    setCompanyModal({ isOpen: true, mode: "update", title: "Update Company" });
  };

  const handleCompanyCreate = async (values: CompanyCreateForm) => {
    if (!selectedRow) return;
    try {
      await clientApi(`/SM01/${selectedRow.id}/company/assign`, {
        method: "POST",
        body: { company_id: values.company_id, is_active: values.is_active ?? true },
      });
      toast.success("Company assigned successfully.");
      setCompanyModal({ isOpen: false, mode: "create", title: "" });
      setCompanyRefresh(!companyRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const handleCompanyUpdate = async (values: CompanyUpdateForm) => {
    if (!selectedRow || !selectedUserCompany) return;
    try {
      await clientApi(`/SM01/${selectedRow.id}/company/${selectedUserCompany.company_id}`, {
        method: "PUT",
        body: { is_active: values.is_active },
      });
      toast.success("Company updated successfully.");
      setCompanyModal({ isOpen: false, mode: "update", title: "" });
      setCompanyRefresh(!companyRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const companyActions: ActionConfig<UserCompanyListItem> = {
    onSearch: () => setCompanyRefresh(!companyRefresh),
    onCreate: openCompanyCreate,
    onUpdate: openCompanyUpdate,
    hideDetail: true,
  };

  const availableCompanies = useMemo(
    () => allCompanies.filter((c) => !assignedCompanyIds.has(c.id)),
    [allCompanies, assignedCompanyIds],
  );

  const availableModules = useMemo(
    () => allModules.filter((m) => m.parent_id !== "" && m.is_active),
    [allModules],
  );

  const privilegeColumns: Column<UserPrivilegeListItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Company", accessor: "company_name", sortable: true, hidden: false },
    { header: "Module", accessor: "module_code", sortable: true, hidden: false,
      formatter: (v: string, row: UserPrivilegeListItem) => `${v} - ${row.module_name}`,
    },
    { header: "Level", accessor: "level", sortable: true, hidden: false,
      formatter: (v: string) => (
        <span className={
          v === "post" ? "text-emerald-600 font-medium" :
          v === "book" ? "text-blue-600 font-medium" :
          v === "view" ? "text-amber-600 font-medium" :
          "text-muted-foreground"
        }>
          {v.charAt(0).toUpperCase() + v.slice(1)}
        </span>
      ),
    },
  ], []);

  const locationColumns: Column<UserLocationListItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Location Type ID", accessor: "location_type_id", sortable: false, hidden: true },
    { header: "Company", accessor: "company_name", sortable: true, hidden: false },
    { header: "Code", accessor: "type_code", sortable: true, hidden: false },
    { header: "Type", accessor: "type_name", sortable: true, hidden: false },
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

  const privilegeCreateForm = useForm<PrivilegeCreateForm>({
    resolver: zodResolver(privilegeCreateSchema),
    defaultValues: { user_company_id: "", module_id: "", level: "view" },
  });

  const privilegeUpdateForm = useForm<PrivilegeUpdateForm>({
    resolver: zodResolver(privilegeUpdateSchema),
    defaultValues: { level: "view" },
  });

  const loadPrivilegeData = useCallback(async (params: {
    search: string; page: number; size: number; sort_by: string; sort_order: "asc" | "desc";
  }) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session || !selectedRow) {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
    try {
      const [privData, compData] = await Promise.all([
        clientApi<{
          data: UserPrivilegeListItem[];
          meta: { total: number; page: number; size: number };
        }>(`/SM01/${selectedRow.id}/privilege`, {
          params: {
            search: params.search,
            page: String(params.page),
            size: String(params.size),
            sort_by: params.sort_by,
            sort_order: params.sort_order,
          },
        }),
        clientApi<{ data: UserCompanyListItem[] }>(`/SM01/${selectedRow.id}/company`),
      ]);
      setUserCompanyList(compData?.data ?? []);
      const data = privData?.data ?? [];
      return {
        data,
        meta: privData?.meta ?? { total: 0, page: 1, size: 10 },
      };
    } catch (err) {
      if (err instanceof ClientApiError) {
        return { data: [], meta: { total: 0, page: 1, size: 10 } };
      }
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
  }, [selectedRow]);

  const openPrivilegeCreate = () => {
    privilegeCreateForm.reset({ user_company_id: "", module_id: "", level: "view" });
    setPrivilegeModal({ isOpen: true, mode: "create", title: "Assign Privilege" });
  };

  const openPrivilegeUpdate = (row: UserPrivilegeListItem) => {
    setSelectedUserPrivilege(row);
    privilegeUpdateForm.reset({ level: row.level });
    setPrivilegeModal({ isOpen: true, mode: "update", title: "Update Privilege" });
  };

  const handlePrivilegeCreate = async (values: PrivilegeCreateForm) => {
    if (!selectedRow) return;
    try {
      await clientApi(`/SM01/${selectedRow.id}/privilege`, {
        method: "POST",
        body: { user_company_id: values.user_company_id, module_id: values.module_id, level: values.level },
      });
      toast.success("Privilege assigned successfully.");
      setPrivilegeModal({ isOpen: false, mode: "create", title: "" });
      setPrivilegeRefresh(!privilegeRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const handlePrivilegeUpdate = async (values: PrivilegeUpdateForm) => {
    if (!selectedRow || !selectedUserPrivilege) return;
    try {
      await clientApi(`/SM01/${selectedRow.id}/privilege/${selectedUserPrivilege.id}`, {
        method: "PUT",
        body: { level: values.level },
      });
      toast.success("Privilege updated successfully.");
      setPrivilegeModal({ isOpen: false, mode: "update", title: "" });
      setPrivilegeRefresh(!privilegeRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const privilegeActions: ActionConfig<UserPrivilegeListItem> = {
    onSearch: () => setPrivilegeRefresh(!privilegeRefresh),
    onCreate: openPrivilegeCreate,
    onUpdate: openPrivilegeUpdate,
    hideDetail: true,
  };

  const locationCreateForm = useForm<LocationCreateForm>({
    resolver: zodResolver(locationCreateSchema),
    defaultValues: { company_id: "", location_type_id: "", is_active: true },
  });

  const locationUpdateForm = useForm<LocationUpdateForm>({
    resolver: zodResolver(locationUpdateSchema),
    defaultValues: { is_active: true },
  });

  const watchedLocationCompanyId = useWatch({ control: locationCreateForm.control, name: "company_id" });

  useEffect(() => {
    if (!watchedLocationCompanyId || watchedLocationCompanyId === "") {
      setAllLocationTypes([]);
      locationCreateForm.setValue("location_type_id", "");
      return;
    }
    clientApi<{ data: LocationTypeSelectItem[] }>(`/SM01/type?company_id=${watchedLocationCompanyId}`)
      .then((data) => setAllLocationTypes(data?.data ?? []))
      .catch(() => setAllLocationTypes([]));
    locationCreateForm.setValue("location_type_id", "");
  }, [watchedLocationCompanyId]);

  const loadLocationData = useCallback(async (params: {
    search: string; page: number; size: number; sort_by: string; sort_order: "asc" | "desc";
  }) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session || !selectedRow) {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
    try {
      const data = await clientApi<{
        data: UserLocationListItem[];
        meta: { total: number; page: number; size: number };
      }>(`/SM01/${selectedRow.id}/location`, {
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
  }, [selectedRow]);

  const openLocationCreate = () => {
    locationCreateForm.reset({ company_id: "", location_type_id: "", is_active: true });
    setAllLocationTypes([]);
    setLocationModal({ isOpen: true, mode: "create", title: "Assign Location" });
  };

  const openLocationUpdate = (row: UserLocationListItem) => {
    setSelectedUserLocation(row);
    locationUpdateForm.reset({ is_active: row.is_active });
    setLocationModal({ isOpen: true, mode: "update", title: "Update Location" });
  };

  const handleLocationCreate = async (values: LocationCreateForm) => {
    if (!selectedRow) return;
    try {
      await clientApi(`/SM01/${selectedRow.id}/location`, {
        method: "POST",
        body: { location_type_id: values.location_type_id, is_active: values.is_active ?? true },
      });
      toast.success("Location assigned successfully.");
      setLocationModal({ isOpen: false, mode: "create", title: "" });
      setLocationRefresh(!locationRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const handleLocationUpdate = async (values: LocationUpdateForm) => {
    if (!selectedRow || !selectedUserLocation) return;
    try {
      await clientApi(`/SM01/${selectedRow.id}/location/${selectedUserLocation.id}`, {
        method: "PUT",
        body: { is_active: values.is_active },
      });
      toast.success("Location updated successfully.");
      setLocationModal({ isOpen: false, mode: "update", title: "" });
      setLocationRefresh(!locationRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const handleLocationDelete = async (row: UserLocationListItem) => {
    if (!selectedRow) return;
    try {
      await clientApi(`/SM01/${selectedRow.id}/location/${row.id}`, {
        method: "DELETE",
      });
      toast.success("Location removed successfully.");
      setLocationRefresh(!locationRefresh);
    } catch (err) {
      if (err instanceof ClientApiError && err.status !== 401) {
        toast.error(getErrorMessage(err instanceof ClientApiError ? err.code : "UNKNOWN_ERROR"));
      }
    }
  };

  const locationActions: ActionConfig<UserLocationListItem> = {
    onSearch: () => setLocationRefresh(!locationRefresh),
    onCreate: openLocationCreate,
    onUpdate: openLocationUpdate,
    onDelete: handleLocationDelete,
    hideDetail: true,
  };

  return (
    <Suspense fallback={<p>Loading users...</p>}>
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
            <FieldGroup className="p-3 max-h-[70vh] overflow-y-auto">
              <Tabs defaultValue="information">
                <TabsList variant="line" className="mb-3">
                  <TabsTrigger value="information">Information</TabsTrigger>
                  <TabsTrigger value="companies">Companies</TabsTrigger>
                  <TabsTrigger value="privilege">Privilege</TabsTrigger>
                  <TabsTrigger value="location">Location</TabsTrigger>
                </TabsList>
                <TabsContent value="information">
                  <FieldGroup className="gap-3">
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Company</FieldLabel>
                      <Input
                        value={companies.find((c) => c.id === selectedRow?.company_id)?.name ?? "Non-HRIS"}
                        disabled readOnly
                      />
                    </Field>
                    <div className="grid gap-3 sm:grid-cols-2">
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Username</FieldLabel>
                        <Input value={selectedRow?.username ?? ""} disabled readOnly />
                      </Field>
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Fullname</FieldLabel>
                        <Input value={selectedRow?.fullname ?? ""} disabled readOnly />
                      </Field>
                    </div>
                    <div className="grid gap-3 sm:grid-cols-2">
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Email</FieldLabel>
                        <Input value={selectedRow?.email ?? ""} disabled readOnly />
                      </Field>
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Password</FieldLabel>
                        <Input value="********" disabled readOnly />
                      </Field>
                    </div>
                    <div className="grid gap-3 sm:grid-cols-3">
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Phone</FieldLabel>
                        <Input value={selectedRow?.phone ?? ""} disabled readOnly />
                      </Field>
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Job</FieldLabel>
                        <Input value={selectedRow?.job ?? ""} disabled readOnly />
                      </Field>
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Role</FieldLabel>
                        <Input value={selectedRow?.role ?? ""} disabled readOnly />
                      </Field>
                    </div>
                    <div className="grid gap-3 sm:grid-cols-3 px-2">
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox checked={selectedRow?.is_admin ?? false} disabled />
                          <FieldLabel>Admin</FieldLabel>
                        </div>
                      </Field>
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox checked={selectedRow?.is_active ?? false} disabled />
                          <FieldLabel>Active</FieldLabel>
                        </div>
                      </Field>
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox checked={selectedRow?.is_hris ?? false} disabled />
                          <FieldLabel>HRIS</FieldLabel>
                        </div>
                      </Field>
                    </div>
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
                <TabsContent value="companies">
                  <div className="space-y-3">
                    {assignedCompanyIds.size > 0 && (
                      <Alert variant="complete">
                        <CircleCheckIcon className="size-4" />
                        <AlertTitle>{assignedCompanyIds.size} company(ies) assigned to this user.</AlertTitle>
                      </Alert>
                    )}
                    <DataTable
                      fetchData={loadCompanyData}
                      columns={companyColumns}
                      actions={companyActions}
                      hideSearch={false}
                      hideSelect={true}
                      hidePaging={false}
                      hideSort={false}
                      hideColumnToggle={true}
                      refreshTrigger={companyRefresh}
                    />
                  </div>
                </TabsContent>
                <TabsContent value="privilege">
                  <div className="space-y-3">
                    {userCompanyList.length > 0 && (
                      <Alert variant="complete">
                        <CircleCheckIcon className="size-4" />
                        <AlertTitle>{userCompanyList.length} privilege(s) configured for this user.</AlertTitle>
                      </Alert>
                    )}
                    <DataTable
                      fetchData={loadPrivilegeData}
                      columns={privilegeColumns}
                      actions={privilegeActions}
                      hideSearch={false}
                      hideSelect={true}
                      hidePaging={false}
                      hideSort={false}
                      hideColumnToggle={true}
                      refreshTrigger={privilegeRefresh}
                    />
                  </div>
                </TabsContent>
                <TabsContent value="location">
                  <div className="space-y-3">
                    <DataTable
                      fetchData={loadLocationData}
                      columns={locationColumns}
                      actions={locationActions}
                      hideSearch={false}
                      hideSelect={true}
                      hidePaging={false}
                      hideSort={false}
                      hideColumnToggle={true}
                      refreshTrigger={locationRefresh}
                    />
                  </div>
                </TabsContent>
              </Tabs>
            </FieldGroup>
          ) : modal.mode === "create" ? (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Controller
                  control={createForm.control}
                  name="company_id"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Company</FieldLabel>
                      <Select
                        value={field.value || ""}
                        onValueChange={(v) => field.onChange(v === "__none__" ? "" : v)}
                        items={[
                          { value: "__none__", label: "Non-HRIS" },
                          ...companies.map((c) => ({ value: c.id, label: c.name })),
                        ]}
                      >
                        <SelectTrigger className="w-full">
                          <SelectValue placeholder="Select company" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="__none__">Non-HRIS</SelectItem>
                          {companies.map((c) => (
                            <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </Field>
                  )}
                />
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={createForm.control}
                    name="username"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="username">Username *</FieldLabel>
                        <Input {...field} id="username" placeholder="Username"
                         maxLength={100} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="fullname"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="fullname">Fullname *</FieldLabel>
                        <Input {...field} id="fullname" placeholder="Fullname"
                         maxLength={255} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={createForm.control}
                    name="email"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="email">Email *</FieldLabel>
                        <Input {...field} id="email" type="email" placeholder="Email"
                         maxLength={255} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="password"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="password">Password *</FieldLabel>
                        <Input {...field} id="password" type="password" placeholder="Password"
                         aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <div className="grid gap-3 sm:grid-cols-3">
                  <Controller
                    control={createForm.control}
                    name="phone"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="phone">Phone</FieldLabel>
                        <Input {...field} id="phone" placeholder="Phone" maxLength={50} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="job"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="job">Job</FieldLabel>
                        <Input {...field} id="job" placeholder="Job" maxLength={100} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="role"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Role</FieldLabel>
                        <Select
                          value={field.value}
                          onValueChange={field.onChange}
                          items={[
                            { value: "staff", label: "Staff" },
                            { value: "manager", label: "Manager" },
                            { value: "admin", label: "Admin" },
                            { value: "viewer", label: "Viewer" },
                          ]}
                        >
                          <SelectTrigger className="w-full">
                            <SelectValue placeholder="Select role" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="staff">Staff</SelectItem>
                            <SelectItem value="manager">Manager</SelectItem>
                            <SelectItem value="admin">Admin</SelectItem>
                            <SelectItem value="viewer">Viewer</SelectItem>
                          </SelectContent>
                        </Select>
                      </Field>
                    )}
                  />
                </div>
                <div className="grid gap-3 sm:grid-cols-3 px-2">
                  <Controller
                    control={createForm.control}
                    name="is_admin"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                          <FieldLabel className="cursor-pointer">Admin</FieldLabel>
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
                          <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                          <FieldLabel className="cursor-pointer">Active</FieldLabel>
                        </div>
                      </Field>
                    )}
                  />
                  <Controller
                    control={createForm.control}
                    name="is_hris"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox checked={field.value} onCheckedChange={field.onChange} disabled />
                          <FieldLabel className="cursor-pointer">HRIS</FieldLabel>
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
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Company</FieldLabel>
                  <Select
                    value={selectedRow?.company_id || ""}
                    disabled
                    items={[
                      { value: "__none__", label: "Non-HRIS" },
                      ...companies.map((c) => ({ value: c.id, label: c.name })),
                    ]}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select company">
                        {companies.find((c) => c.id === selectedRow?.company_id)?.name ?? "Non-HRIS"}
                      </SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="__none__">Non-HRIS</SelectItem>
                      {companies.map((c) => (
                        <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </Field>
                <div className="grid gap-3 sm:grid-cols-2">
                  <Field className="gap-1.5 pt-2 px-2">
                    <FieldLabel>Username</FieldLabel>
                    <Input value={selectedRow?.username ?? ""} disabled readOnly />
                  </Field>
                  <Controller
                    control={updateForm.control}
                    name="fullname"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="upd-fullname">Fullname *</FieldLabel>
                        <Input {...field} id="upd-fullname" placeholder="Fullname"
                         maxLength={255} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                </div>
                <div className="grid gap-3 sm:grid-cols-2">
                  <Controller
                    control={updateForm.control}
                    name="email"
                    render={({ field, fieldState }) => (
                      <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                        <FieldLabel htmlFor="upd-email">Email *</FieldLabel>
                        <Input {...field} id="upd-email" type="email" placeholder="Email"
                         maxLength={255} aria-invalid={fieldState.invalid} />
                        {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                      </Field>
                    )}
                  />
                  <Controller
                    control={updateForm.control}
                    name="password"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="upd-password">Password</FieldLabel>
                        <Input {...field} id="upd-password" type="password"
                         placeholder="Leave empty to keep current" value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                </div>
                <div className="grid gap-3 sm:grid-cols-3">
                  <Controller
                    control={updateForm.control}
                    name="phone"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="upd-phone">Phone</FieldLabel>
                        <Input {...field} id="upd-phone" placeholder="Phone" maxLength={50} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                  <Controller
                    control={updateForm.control}
                    name="job"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel htmlFor="upd-job">Job</FieldLabel>
                        <Input {...field} id="upd-job" placeholder="Job" maxLength={100} value={field.value ?? ""} />
                      </Field>
                    )}
                  />
                  <Controller
                    control={updateForm.control}
                    name="role"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2 px-2">
                        <FieldLabel>Role</FieldLabel>
                        <Select
                          value={field.value}
                          onValueChange={field.onChange}
                          items={[
                            { value: "staff", label: "Staff" },
                            { value: "manager", label: "Manager" },
                            { value: "admin", label: "Admin" },
                            { value: "viewer", label: "Viewer" },
                          ]}
                        >
                          <SelectTrigger className="w-full">
                            <SelectValue placeholder="Select role" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="staff">Staff</SelectItem>
                            <SelectItem value="manager">Manager</SelectItem>
                            <SelectItem value="admin">Admin</SelectItem>
                            <SelectItem value="viewer">Viewer</SelectItem>
                          </SelectContent>
                        </Select>
                      </Field>
                    )}
                  />
                </div>
                <div className="grid gap-3 sm:grid-cols-3 px-2">
                  <Controller
                    control={updateForm.control}
                    name="is_admin"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                          <FieldLabel className="cursor-pointer">Admin</FieldLabel>
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
                          <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                          <FieldLabel className="cursor-pointer">Active</FieldLabel>
                        </div>
                      </Field>
                    )}
                  />
                  <Controller
                    control={updateForm.control}
                    name="is_hris"
                    render={({ field }) => (
                      <Field className="gap-1.5 pt-2">
                        <div className="flex items-center gap-2">
                          <Checkbox checked={field.value} onCheckedChange={field.onChange} disabled />
                          <FieldLabel className="cursor-pointer">HRIS</FieldLabel>
                        </div>
                      </Field>
                    )}
                  />
                </div>
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>

        <DataDialog
          isOpen={companyModal.isOpen}
          mode={companyModal.mode}
          title={companyModal.title}
          onClose={() => {
            setCompanyModal({ isOpen: false, mode: "create", title: "" });
            setSelectedUserCompany(null);
          }}
          onSubmit={
            companyModal.mode === "create"
              ? companyCreateForm.handleSubmit(handleCompanyCreate)
              : companyUpdateForm.handleSubmit(handleCompanyUpdate)
          }
        >
          {companyModal.mode === "create" ? (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Controller
                  control={companyCreateForm.control}
                  name="company_id"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Company *</FieldLabel>
                      <Select
                        value={field.value || ""}
                        onValueChange={field.onChange}
                        items={availableCompanies.map((c) => ({ value: c.id, label: c.name }))}
                      >
                        <SelectTrigger className="w-full" aria-invalid={fieldState.invalid}>
                          <SelectValue placeholder="Select company" />
                        </SelectTrigger>
                        <SelectContent>
                          {availableCompanies.length === 0 ? (
                            <SelectItem value="__none__" disabled>No available companies</SelectItem>
                          ) : (
                            availableCompanies.map((c) => (
                              <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                            ))
                          )}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <Controller
                  control={companyCreateForm.control}
                  name="is_active"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <div className="flex items-center gap-2">
                        <Checkbox checked={field.value ?? true} onCheckedChange={field.onChange} />
                        <FieldLabel className="cursor-pointer">Active</FieldLabel>
                      </div>
                    </Field>
                  )}
                />
              </FieldGroup>
            </FieldSet>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Company</FieldLabel>
                  <Input value={selectedUserCompany?.company_name ?? ""} disabled readOnly />
                </Field>
                <Controller
                  control={companyUpdateForm.control}
                  name="is_active"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <div className="flex items-center gap-2">
                        <Checkbox checked={field.value} onCheckedChange={field.onChange} />
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
          isOpen={privilegeModal.isOpen}
          mode={privilegeModal.mode}
          title={privilegeModal.title}
          onClose={() => {
            setPrivilegeModal({ isOpen: false, mode: "create", title: "" });
            setSelectedUserPrivilege(null);
          }}
          onSubmit={
            privilegeModal.mode === "create"
              ? privilegeCreateForm.handleSubmit(handlePrivilegeCreate)
              : privilegeUpdateForm.handleSubmit(handlePrivilegeUpdate)
          }
        >
          {privilegeModal.mode === "create" ? (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Controller
                  control={privilegeCreateForm.control}
                  name="user_company_id"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Company Assignment *</FieldLabel>
                      <Select
                        value={field.value || ""}
                        onValueChange={field.onChange}
                        items={userCompanyList.map((c) => ({ value: c.id, label: c.company_name }))}
                      >
                        <SelectTrigger className="w-full" aria-invalid={fieldState.invalid}>
                          <SelectValue placeholder="Select company assignment" />
                        </SelectTrigger>
                        <SelectContent>
                          {userCompanyList.length === 0 ? (
                            <SelectItem value="__none__" disabled>No company assignments</SelectItem>
                          ) : (
                            userCompanyList.map((c) => (
                              <SelectItem key={c.id} value={c.id}>{c.company_name}</SelectItem>
                            ))
                          )}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <Controller
                  control={privilegeCreateForm.control}
                  name="module_id"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Module *</FieldLabel>
                      <Select
                        value={field.value || ""}
                        onValueChange={field.onChange}
                        items={availableModules.map((m) => ({ value: m.id, label: `${m.code} - ${m.name}` }))}
                      >
                        <SelectTrigger className="w-full" aria-invalid={fieldState.invalid}>
                          <SelectValue placeholder="Select module" />
                        </SelectTrigger>
                        <SelectContent>
                          {availableModules.length === 0 ? (
                            <SelectItem value="__none__" disabled>No modules available</SelectItem>
                          ) : (
                            availableModules.map((m) => (
                              <SelectItem key={m.id} value={m.id}>{m.code} - {m.name}</SelectItem>
                            ))
                          )}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <Controller
                  control={privilegeCreateForm.control}
                  name="level"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Level *</FieldLabel>
                      <Select
                        value={field.value || "view"}
                        onValueChange={field.onChange}
                        items={levelOptions}
                      >
                        <SelectTrigger className="w-full" aria-invalid={fieldState.invalid}>
                          <SelectValue placeholder="Select level" />
                        </SelectTrigger>
                        <SelectContent>
                          {levelOptions.map((l) => (
                            <SelectItem key={l.value} value={l.value}>{l.label}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
              </FieldGroup>
            </FieldSet>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Company</FieldLabel>
                  <Input value={selectedUserPrivilege?.company_name ?? ""} disabled readOnly />
                </Field>
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Module</FieldLabel>
                  <Input value={selectedUserPrivilege
                    ? `${selectedUserPrivilege.module_code} - ${selectedUserPrivilege.module_name}`
                    : ""} disabled readOnly />
                </Field>
                <Controller
                  control={privilegeUpdateForm.control}
                  name="level"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Level *</FieldLabel>
                      <Select
                        value={field.value || "view"}
                        onValueChange={field.onChange}
                        items={levelOptions}
                      >
                        <SelectTrigger className="w-full" aria-invalid={fieldState.invalid}>
                          <SelectValue placeholder="Select level" />
                        </SelectTrigger>
                        <SelectContent>
                          {levelOptions.map((l) => (
                            <SelectItem key={l.value} value={l.value}>{l.label}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>

        <DataDialog
          isOpen={locationModal.isOpen}
          mode={locationModal.mode}
          title={locationModal.title}
          onClose={() => {
            setLocationModal({ isOpen: false, mode: "create", title: "" });
            setSelectedUserLocation(null);
          }}
          onSubmit={
            locationModal.mode === "create"
              ? locationCreateForm.handleSubmit(handleLocationCreate)
              : locationUpdateForm.handleSubmit(handleLocationUpdate)
          }
        >
          {locationModal.mode === "create" ? (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Controller
                  control={locationCreateForm.control}
                  name="company_id"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Company *</FieldLabel>
                      <Select
                        value={field.value || ""}
                        onValueChange={field.onChange}
                        items={allCompanies.map((c) => ({ value: c.id, label: c.name }))}
                      >
                        <SelectTrigger className="w-full" aria-invalid={fieldState.invalid}>
                          <SelectValue placeholder="Select company" />
                        </SelectTrigger>
                        <SelectContent>
                          {allCompanies.length === 0 ? (
                            <SelectItem value="__none__" disabled>No companies</SelectItem>
                          ) : (
                            allCompanies.map((c) => (
                              <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                            ))
                          )}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <Controller
                  control={locationCreateForm.control}
                  name="location_type_id"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Location Type *</FieldLabel>
                      <Select
                        value={field.value || ""}
                        onValueChange={field.onChange}
                        items={allLocationTypes.map((a) => ({ value: a.id, label: `${a.code} - ${a.name}` }))}
                      >
                        <SelectTrigger className="w-full" aria-invalid={fieldState.invalid}>
                          <SelectValue placeholder="Select location type" />
                        </SelectTrigger>
                        <SelectContent>
                          {allLocationTypes.length === 0 ? (
                            <SelectItem value="__none__" disabled>
                              {!watchedLocationCompanyId ? "Select company first" : "No available types"}
                            </SelectItem>
                          ) : (
                            allLocationTypes.map((a) => (
                              <SelectItem key={a.id} value={a.id}>{a.code} - {a.name}</SelectItem>
                            ))
                          )}
                        </SelectContent>
                      </Select>
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <Controller
                  control={locationCreateForm.control}
                  name="is_active"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <div className="flex items-center gap-2">
                        <Checkbox checked={field.value ?? true} onCheckedChange={field.onChange} />
                        <FieldLabel className="cursor-pointer">Active</FieldLabel>
                      </div>
                    </Field>
                  )}
                />
              </FieldGroup>
            </FieldSet>
          ) : (
            <FieldSet>
              <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
                <Field className="gap-1.5 pt-2 px-2">
                  <FieldLabel>Location Type</FieldLabel>
                  <Input
                    value={selectedUserLocation ? `${selectedUserLocation.type_code} - ${selectedUserLocation.type_name}` : ""}
                    disabled readOnly
                  />
                </Field>
                <Controller
                  control={locationUpdateForm.control}
                  name="is_active"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <div className="flex items-center gap-2">
                        <Checkbox checked={field.value} onCheckedChange={field.onChange} />
                        <FieldLabel className="cursor-pointer">Active</FieldLabel>
                      </div>
                    </Field>
                  )}
                />
              </FieldGroup>
            </FieldSet>
          )}
        </DataDialog>
      </div>
    </Suspense>
  );
}
