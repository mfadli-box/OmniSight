"use client";

import { Controller, useForm } from "react-hook-form";
import { Suspense, useCallback, useEffect, useMemo, useState } from "react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { storageKey, parseSession } from "@/lib/utility";
import { clientApi, ClientApiError } from "@/lib/client-api";
import { getErrorMessage } from "@/lib/error-message";
import { toast } from "sonner";
import DataTable, { Column, ActionConfig } from "@/uix/datatable";
import DataDialog from "@/uix/datadialog";
import { Field, FieldError, FieldGroup, FieldLabel, FieldSet } from "@/uix/field";
import { Input } from "@/uix/input";
import { Textarea } from "@/uix/textarea";
import { Checkbox } from "@/uix/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/uix/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/uix/tabs";

interface LocationTypeItem {
  id: string;
  code: string;
  name: string;
  description: string;
  icon: string;
  color: string;
  is_active: boolean;
  created_at: string;
}

interface LocationItem {
  id: string;
  location_type_id: string;
  parent_id: string;
  code: string;
  name: string;
  description: string;
  address: string;
  city: string;
  province: string;
  country: string;
  postal_code: string;
  latitude: string;
  longitude: string;
  phone: string;
  email: string;
  timezone: string;
  status: string;
  is_active: boolean;
  created_at: string;
}

interface LocationTypeSelectItem {
  id: string;
  code: string;
  name: string;
}

interface LocationSelectItem {
  id: string;
  parent_id: string;
  code: string;
  name: string;
}

const createTypeSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  icon: z.string().optional(),
  color: z.string().optional(),
  is_active: z.boolean().optional(),
});

const updateTypeSchema = z.object({
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  icon: z.string().optional(),
  color: z.string().optional(),
  is_active: z.boolean().optional(),
});

const createLocationSchema = z.object({
  location_type_id: z.string().min(1, "Location type is required"),
  parent_id: z.string().optional(),
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  address: z.string().optional(),
  city: z.string().optional(),
  province: z.string().optional(),
  country: z.string().optional(),
  postal_code: z.string().optional(),
  latitude: z.number().optional(),
  longitude: z.number().optional(),
  phone: z.string().optional(),
  email: z.string().optional(),
  timezone: z.string().optional(),
  status: z.string().optional(),
  is_active: z.boolean().optional(),
});

const updateLocationSchema = z.object({
  location_type_id: z.string().min(1, "Location type is required"),
  parent_id: z.string().optional(),
  code: z.string().min(1, "Code is required"),
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  address: z.string().optional(),
  city: z.string().optional(),
  province: z.string().optional(),
  country: z.string().optional(),
  postal_code: z.string().optional(),
  latitude: z.number().optional(),
  longitude: z.number().optional(),
  phone: z.string().optional(),
  email: z.string().optional(),
  timezone: z.string().optional(),
  status: z.string().optional(),
  is_active: z.boolean().optional(),
});

type CreateTypeForm = z.infer<typeof createTypeSchema>;
type UpdateTypeForm = z.infer<typeof updateTypeSchema>;
type CreateLocationForm = z.infer<typeof createLocationSchema>;
type UpdateLocationForm = z.infer<typeof updateLocationSchema>;

export default function Page() {
  const [activeTab, setActiveTab] = useState("locations");
  const [selectedRow, setSelectedRow] = useState<LocationTypeItem | LocationItem | null>(null);
  const [refresh, setRefresh] = useState(false);
  const [modal, setModal] = useState<{
    isOpen: boolean;
    mode: "create" | "update" | "detail";
    title: string;
    type: "type" | "location";
  }>({ isOpen: false, mode: "create", title: "", type: "location" });

  const [locationTypes, setLocationTypes] = useState<LocationTypeSelectItem[]>([]);
  const [locations, setLocations] = useState<LocationSelectItem[]>([]);

  const fetchLocationTypes = useCallback(async () => {
    try {
      const res = await clientApi<{ data: LocationTypeSelectItem[] }>("/SM06/type/select");
      setLocationTypes((res as any)?.data ?? []);
    } catch {
      setLocationTypes([]);
    }
  }, []);

  const fetchLocations = useCallback(async () => {
    try {
      const res = await clientApi<{ data: LocationSelectItem[] }>("/SM06/select");
      setLocations((res as any)?.data ?? []);
    } catch {
      setLocations([]);
    }
  }, []);

  useEffect(() => {
    fetchLocationTypes();
    fetchLocations();
  }, [fetchLocationTypes, fetchLocations]);

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
        data: LocationItem[];
        meta: { total: number; page: number; size: number };
      }>("/SM06", {
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

  const loadTypeData = useCallback(async (params: {
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
        data: LocationTypeItem[];
        meta: { total: number; page: number; size: number };
      }>("/SM06/type", {
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

  const typeColumns: Column<LocationTypeItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Code", accessor: "code", sortable: true },
    { header: "Name", accessor: "name", sortable: true },
    { header: "Description", accessor: "description", sortable: false },
    { header: "Icon", accessor: "icon", sortable: false },
    { header: "Color", accessor: "color", sortable: false },
    {
      header: "Active", accessor: "is_active", sortable: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-red-600 font-medium"}>
          {v ? "Active" : "Inactive"}
        </span>
      ),
    },
    { header: "Created At", accessor: "created_at", sortable: true },
  ], []);

  const locationColumns: Column<LocationItem>[] = useMemo(() => [
    { header: "ID", accessor: "id", sortable: false, hidden: true },
    { header: "Code", accessor: "code", sortable: true },
    { header: "Name", accessor: "name", sortable: true },
    { header: "City", accessor: "city", sortable: false },
    { header: "Province", accessor: "province", sortable: false },
    { header: "Country", accessor: "country", sortable: false },
    {
      header: "Status", accessor: "status", sortable: true,
      formatter: (v: string) => (
        <span className={`px-2 py-1 rounded-full text-xs font-medium ${
          v === "active" ? "bg-emerald-100 text-emerald-800" :
          v === "inactive" ? "bg-red-100 text-red-800" :
          v === "maintenance" ? "bg-yellow-100 text-yellow-800" :
          "bg-gray-100 text-gray-800"
        }`}>
          {v.charAt(0).toUpperCase() + v.slice(1)}
        </span>
      ),
    },
    {
      header: "Active", accessor: "is_active", sortable: false,
      formatter: (v: boolean) => (
        <span className={v ? "text-emerald-600 font-medium" : "text-red-600 font-medium"}>
          {v ? "Active" : "Inactive"}
        </span>
      ),
    },
    { header: "Created At", accessor: "created_at", sortable: true },
  ], []);

  const typeActions: ActionConfig<LocationTypeItem> = useMemo(() => ({
    onSearch: () => setRefresh((r) => !r),
    onCreate: () => { createTypeForm.reset(); setModal({ isOpen: true, mode: "create", title: "Create Location Type", type: "type" }); },
    onDetail: (row) => { setSelectedRow(row); setModal({ isOpen: true, mode: "detail", title: "Location Type Detail", type: "type" }); },
    onUpdate: (row) => { setSelectedRow(row); updateTypeForm.reset({ code: row.code, name: row.name, description: row.description, icon: row.icon, color: row.color, is_active: row.is_active }); setModal({ isOpen: true, mode: "update", title: "Update Location Type", type: "type" }); },
    onDelete: (row) => { handleDeleteType(row); },
  }), []);

  const locationActions: ActionConfig<LocationItem> = useMemo(() => ({
    onSearch: () => setRefresh((r) => !r),
    onCreate: () => { createLocationForm.reset(); setModal({ isOpen: true, mode: "create", title: "Create Location", type: "location" }); },
    onDetail: (row) => { setSelectedRow(row); setModal({ isOpen: true, mode: "detail", title: "Location Detail", type: "location" }); },
    onUpdate: (row) => { setSelectedRow(row); updateLocationForm.reset({ location_type_id: row.location_type_id, parent_id: row.parent_id, code: row.code, name: row.name, description: row.description, address: row.address, city: row.city, province: row.province, country: row.country, postal_code: row.postal_code, latitude: row.latitude ? Number(row.latitude) : undefined, longitude: row.longitude ? Number(row.longitude) : undefined, phone: row.phone, email: row.email, timezone: row.timezone, status: row.status, is_active: row.is_active }); setModal({ isOpen: true, mode: "update", title: "Update Location", type: "location" }); },
    onDelete: (row) => { handleDeleteLocation(row); },
  }), []);

  const createTypeForm = useForm<CreateTypeForm>({
    resolver: zodResolver(createTypeSchema),
    defaultValues: {
      code: "",
      name: "",
      description: "",
      icon: "",
      color: "#000000",
      is_active: true,
    },
  });
  const updateTypeForm = useForm<UpdateTypeForm>({
    resolver: zodResolver(updateTypeSchema),
    defaultValues: {
      code: "",
      name: "",
      description: "",
      icon: "",
      color: "#000000",
      is_active: true,
    },
  });
  const createLocationForm = useForm<CreateLocationForm>({
    resolver: zodResolver(createLocationSchema),
    defaultValues: {
      location_type_id: "",
      parent_id: "",
      code: "",
      name: "",
      description: "",
      address: "",
      city: "",
      province: "",
      country: "",
      postal_code: "",
      phone: "",
      email: "",
      timezone: "",
      status: "active",
      is_active: true,
    },
  });
  const updateLocationForm = useForm<UpdateLocationForm>({
    resolver: zodResolver(updateLocationSchema),
    defaultValues: {
      location_type_id: "",
      parent_id: "",
      code: "",
      name: "",
      description: "",
      address: "",
      city: "",
      province: "",
      country: "",
      postal_code: "",
      phone: "",
      email: "",
      timezone: "",
      status: "active",
      is_active: true,
    },
  });

  const handleCreateType = async (data: CreateTypeForm) => {
    try {
      await clientApi("/SM06/type", { method: "POST", body: data });
      toast.success("Location type created successfully");
      setModal({ ...modal, isOpen: false });
      setRefresh((r) => !r);
    } catch (err) {
      if (err instanceof ClientApiError) toast.error(getErrorMessage(err.code));
    }
  };

  const handleUpdateType = async (data: UpdateTypeForm) => {
    if (!selectedRow) return;
    try {
      await clientApi(`/SM06/type/${selectedRow.id}`, { method: "PUT", body: data });
      toast.success("Location type updated successfully");
      setModal({ ...modal, isOpen: false });
      setRefresh((r) => !r);
    } catch (err) {
      if (err instanceof ClientApiError) toast.error(getErrorMessage(err.code));
    }
  };

  const handleDeleteType = async (row: LocationTypeItem) => {
    if (!confirm(`Delete "${row.name}"?`)) return;
    try {
      await clientApi(`/SM06/type/${row.id}`, { method: "DELETE" });
      toast.success("Location type deleted successfully");
      setRefresh((r) => !r);
    } catch (err) {
      if (err instanceof ClientApiError) toast.error(getErrorMessage(err.code));
    }
  };

  const handleCreateLocation = async (data: CreateLocationForm) => {
    try {
      await clientApi("/SM06", { method: "POST", body: data });
      toast.success("Location created successfully");
      setModal({ ...modal, isOpen: false });
      setRefresh((r) => !r);
    } catch (err) {
      if (err instanceof ClientApiError) toast.error(getErrorMessage(err.code));
    }
  };

  const handleUpdateLocation = async (data: UpdateLocationForm) => {
    if (!selectedRow) return;
    try {
      await clientApi(`/SM06/${selectedRow.id}`, { method: "PUT", body: data });
      toast.success("Location updated successfully");
      setModal({ ...modal, isOpen: false });
      setRefresh((r) => !r);
    } catch (err) {
      if (err instanceof ClientApiError) toast.error(getErrorMessage(err.code));
    }
  };

  const handleDeleteLocation = async (row: LocationItem) => {
    if (!confirm(`Delete "${row.name}"?`)) return;
    try {
      await clientApi(`/SM06/${row.id}`, { method: "DELETE" });
      toast.success("Location deleted successfully");
      setRefresh((r) => !r);
    } catch (err) {
      if (err instanceof ClientApiError) toast.error(getErrorMessage(err.code));
    }
  };

  const getTypeName = (typeId: string) => {
    const type = locationTypes.find((t) => t.id === typeId);
    return type ? `${type.code} - ${type.name}` : "-";
  };

  const getParentName = (parentId: string) => {
    if (!parentId) return "-";
    const parent = locations.find((l) => l.id === parentId);
    return parent ? `${parent.code} - ${parent.name}` : "-";
  };

  const getTypeFormErrors = () => modal.mode === "create" ? createTypeForm.formState.errors : updateTypeForm.formState.errors;
  const getLocFormErrors = () => modal.mode === "create" ? createLocationForm.formState.errors : updateLocationForm.formState.errors;

  return (
    <Suspense fallback={<p>Loading...</p>}>
      <div className="p-3 max-w-8xl mx-auto space-y-3">
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList variant="line">
            <TabsTrigger value="locations">Locations</TabsTrigger>
            <TabsTrigger value="types">Location Types</TabsTrigger>
          </TabsList>

          <TabsContent value="locations">
            <DataTable<LocationItem>
              fetchData={loadData}
              columns={locationColumns}
              actions={locationActions}
              hideSelect={true}
              hideSearch={false}
              hidePaging={false}
              hideSort={false}
              hideColumnToggle={true}
              refreshTrigger={refresh}
            />
          </TabsContent>

          <TabsContent value="types">
            <DataTable<LocationTypeItem>
              fetchData={loadTypeData}
              columns={typeColumns}
              actions={typeActions}
              hideSelect={true}
              hideSearch={false}
              hidePaging={false}
              hideSort={false}
              hideColumnToggle={true}
              refreshTrigger={refresh}
            />
          </TabsContent>
        </Tabs>

        <DataDialog
          isOpen={modal.isOpen && modal.mode !== "detail" && modal.type === "type"}
          mode={modal.mode}
          title={modal.title}
          onClose={() => setModal({ ...modal, isOpen: false })}
          onSubmit={
            modal.mode === "create"
              ? createTypeForm.handleSubmit(handleCreateType)
              : updateTypeForm.handleSubmit(handleUpdateType)
          }
        >
          <FieldSet>
            <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
              <Controller
                control={modal.mode === "create" ? createTypeForm.control : updateTypeForm.control}
                name="code"
                render={({ field, fieldState }) => (
                  <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                    <FieldLabel>Code *</FieldLabel>
                    <Input {...field} maxLength={100} aria-invalid={fieldState.invalid} />
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </Field>
                )}
              />
              <Controller
                control={modal.mode === "create" ? createTypeForm.control : updateTypeForm.control}
                name="name"
                render={({ field, fieldState }) => (
                  <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                    <FieldLabel>Name *</FieldLabel>
                    <Input {...field} maxLength={255} aria-invalid={fieldState.invalid} />
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </Field>
                )}
              />
              <Controller
                control={modal.mode === "create" ? createTypeForm.control : updateTypeForm.control}
                name="description"
                render={({ field }) => (
                  <Field className="gap-1.5 pt-2 px-2">
                    <FieldLabel>Description</FieldLabel>
                    <Textarea {...field} value={field.value ?? ""} />
                  </Field>
                )}
              />
              <Controller
                control={modal.mode === "create" ? createTypeForm.control : updateTypeForm.control}
                name="icon"
                render={({ field }) => (
                  <Field className="gap-1.5 pt-2 px-2">
                    <FieldLabel>Icon</FieldLabel>
                    <Input {...field} value={field.value ?? ""} placeholder="Icon name (optional)" />
                  </Field>
                )}
              />
              <Controller
                control={modal.mode === "create" ? createTypeForm.control : updateTypeForm.control}
                name="color"
                render={({ field }) => (
                  <Field className="gap-1.5 pt-2 px-2">
                    <FieldLabel>Color</FieldLabel>
                    <Input {...field} value={field.value ?? ""} type="color" />
                  </Field>
                )}
              />
              <Controller
                control={modal.mode === "create" ? createTypeForm.control : updateTypeForm.control}
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
        </DataDialog>

        <DataDialog
          isOpen={modal.isOpen && modal.mode !== "detail" && modal.type === "location"}
          mode={modal.mode}
          title={modal.title}
          onClose={() => setModal({ ...modal, isOpen: false })}
          onSubmit={
            modal.mode === "create"
              ? createLocationForm.handleSubmit(handleCreateLocation)
              : updateLocationForm.handleSubmit(handleUpdateLocation)
          }
        >
          <FieldSet>
            <FieldGroup className="gap-3 p-3 max-h-[70vh] overflow-y-auto">
              <div className="grid grid-cols-2 gap-3">
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="code"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Code *</FieldLabel>
                      <Input {...field} maxLength={100} aria-invalid={fieldState.invalid} />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="name"
                  render={({ field, fieldState }) => (
                    <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                      <FieldLabel>Name *</FieldLabel>
                      <Input {...field} maxLength={255} aria-invalid={fieldState.invalid} />
                      {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                    </Field>
                  )}
                />
              </div>
              <Controller
                control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                name="location_type_id"
                render={({ field, fieldState }) => (
                  <Field className="gap-1.5 pt-2 px-2" data-invalid={fieldState.invalid}>
                    <FieldLabel>Location Type *</FieldLabel>
                    <Select
                      value={field.value ?? ""}
                      onValueChange={field.onChange}
                      items={locationTypes.map((type) => ({
                        value: type.id,
                        label: `${type.code} - ${type.name}`,
                      }))}
                    >
                      <SelectTrigger aria-invalid={fieldState.invalid}>
                        <SelectValue placeholder="Select location type" />
                      </SelectTrigger>
                      <SelectContent>
                        {locationTypes.map((type) => (
                          <SelectItem key={type.id} value={type.id}>
                            {type.code} - {type.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
                  </Field>
                )}
              />
              <Controller
                control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                name="parent_id"
                render={({ field }) => (
                  <Field className="gap-1.5 pt-2 px-2">
                    <FieldLabel>Parent Location</FieldLabel>
                    <Select
                      value={field.value ?? ""}
                      onValueChange={(v) => field.onChange(v || undefined)}
                      items={[
                        { value: "", label: "None (Root)" },
                        ...locations.filter((l) => l.id !== selectedRow?.id).map((loc) => ({
                          value: loc.id,
                          label: `${loc.code} - ${loc.name}`,
                        })),
                      ]}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Select parent location (optional)" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="">None (Root)</SelectItem>
                        {locations.filter((l) => l.id !== selectedRow?.id).map((loc) => (
                          <SelectItem key={loc.id} value={loc.id}>
                            {loc.code} - {loc.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </Field>
                )}
              />
              <Controller
                control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                name="description"
                render={({ field }) => (
                  <Field className="gap-1.5 pt-2 px-2">
                    <FieldLabel>Description</FieldLabel>
                    <Textarea {...field} value={field.value ?? ""} />
                  </Field>
                )}
              />
              <Controller
                control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                name="address"
                render={({ field }) => (
                  <Field className="gap-1.5 pt-2 px-2">
                    <FieldLabel>Address</FieldLabel>
                    <Textarea {...field} value={field.value ?? ""} />
                  </Field>
                )}
              />
              <div className="grid grid-cols-3 gap-3">
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="city"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>City</FieldLabel>
                      <Input {...field} value={field.value ?? ""} />
                    </Field>
                  )}
                />
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="province"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Province</FieldLabel>
                      <Input {...field} value={field.value ?? ""} />
                    </Field>
                  )}
                />
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="country"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Country</FieldLabel>
                      <Input {...field} value={field.value ?? ""} />
                    </Field>
                  )}
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="postal_code"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Postal Code</FieldLabel>
                      <Input {...field} value={field.value ?? ""} />
                    </Field>
                  )}
                />
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="timezone"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Timezone</FieldLabel>
                      <Input {...field} value={field.value ?? ""} />
                    </Field>
                  )}
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="phone"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Phone</FieldLabel>
                      <Input {...field} value={field.value ?? ""} />
                    </Field>
                  )}
                />
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="email"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Email</FieldLabel>
                      <Input {...field} value={field.value ?? ""} type="email" />
                    </Field>
                  )}
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="latitude"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Latitude</FieldLabel>
                      <Input {...field} type="number" step="any" value={field.value ?? ""} onChange={(e) => field.onChange(e.target.value ? Number(e.target.value) : undefined)} />
                    </Field>
                  )}
                />
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="longitude"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Longitude</FieldLabel>
                      <Input {...field} type="number" step="any" value={field.value ?? ""} onChange={(e) => field.onChange(e.target.value ? Number(e.target.value) : undefined)} />
                    </Field>
                  )}
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="status"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <FieldLabel>Status</FieldLabel>
                      <Select
                        value={field.value ?? "active"}
                        onValueChange={field.onChange}
                        items={[
                          { value: "active", label: "Active" },
                          { value: "inactive", label: "Inactive" },
                          { value: "maintenance", label: "Maintenance" },
                          { value: "closed", label: "Closed" },
                        ]}
                      >
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="active">Active</SelectItem>
                          <SelectItem value="inactive">Inactive</SelectItem>
                          <SelectItem value="maintenance">Maintenance</SelectItem>
                          <SelectItem value="closed">Closed</SelectItem>
                        </SelectContent>
                      </Select>
                    </Field>
                  )}
                />
                <Controller
                  control={modal.mode === "create" ? createLocationForm.control : updateLocationForm.control}
                  name="is_active"
                  render={({ field }) => (
                    <Field className="gap-1.5 pt-2 px-2">
                      <div className="flex items-center gap-2 pt-6">
                        <Checkbox checked={field.value ?? true} onCheckedChange={field.onChange} />
                        <FieldLabel className="cursor-pointer">Active</FieldLabel>
                      </div>
                    </Field>
                  )}
                />
              </div>
            </FieldGroup>
          </FieldSet>
        </DataDialog>

        <DataDialog
          isOpen={modal.isOpen && modal.mode === "detail"}
          mode="detail"
          title={modal.title}
          onClose={() => setModal({ ...modal, isOpen: false })}
        >
          <FieldGroup className="p-3 max-h-[70vh] overflow-y-auto">
            {selectedRow && modal.type === "type" && (
              <div className="grid grid-cols-2 gap-3">
                {([
                  ["Code", (selectedRow as LocationTypeItem).code],
                  ["Name", (selectedRow as LocationTypeItem).name],
                  ["Description", (selectedRow as LocationTypeItem).description || "-"],
                  ["Icon", (selectedRow as LocationTypeItem).icon || "-"],
                  ["Color", (selectedRow as LocationTypeItem).color || "-"],
                  ["Active", (selectedRow as LocationTypeItem).is_active ? "Active" : "Inactive"],
                  ["Created At", (selectedRow as LocationTypeItem).created_at],
                ] as [string, string][]).map(([label, value]) => (
                  <div key={label} className="flex flex-col gap-1 border-b px-1 py-1.5 last:border-b-0 sm:flex-row sm:items-start sm:gap-3">
                    <FieldLabel className="w-32 shrink-0 text-sm text-muted-foreground">{label}</FieldLabel>
                    <span className="text-sm break-all">{value}</span>
                  </div>
                ))}
              </div>
            )}
            {selectedRow && modal.type === "location" && (
              <div className="grid grid-cols-2 gap-3">
                {([
                  ["Code", (selectedRow as LocationItem).code],
                  ["Name", (selectedRow as LocationItem).name],
                  ["Type", getTypeName((selectedRow as LocationItem).location_type_id)],
                  ["Parent", getParentName((selectedRow as LocationItem).parent_id)],
                  ["Description", (selectedRow as LocationItem).description || "-"],
                  ["Address", (selectedRow as LocationItem).address || "-"],
                  ["City", (selectedRow as LocationItem).city || "-"],
                  ["Province", (selectedRow as LocationItem).province || "-"],
                  ["Country", (selectedRow as LocationItem).country],
                  ["Postal Code", (selectedRow as LocationItem).postal_code || "-"],
                  ["Phone", (selectedRow as LocationItem).phone || "-"],
                  ["Email", (selectedRow as LocationItem).email || "-"],
                  ["Timezone", (selectedRow as LocationItem).timezone || "-"],
                  ["Status", (selectedRow as LocationItem).status ? (selectedRow as LocationItem).status.charAt(0).toUpperCase() + (selectedRow as LocationItem).status.slice(1) : "-"],
                  ["Active", (selectedRow as LocationItem).is_active ? "Active" : "Inactive"],
                  ["Created At", (selectedRow as LocationItem).created_at],
                ] as [string, string][]).map(([label, value]) => (
                  <div key={label} className="flex flex-col gap-1 border-b px-1 py-1.5 last:border-b-0 sm:flex-row sm:items-start sm:gap-3">
                    <FieldLabel className="w-32 shrink-0 text-sm text-muted-foreground">{label}</FieldLabel>
                    <span className="text-sm break-all">{value}</span>
                  </div>
                ))}
              </div>
            )}
          </FieldGroup>
        </DataDialog>
      </div>
    </Suspense>
  );
}
