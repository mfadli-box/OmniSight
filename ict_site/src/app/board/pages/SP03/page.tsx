"use client";

import { storageKey, parseSession } from "@/lib/utility";
import { clientApi } from "@/lib/client-api";
import { Suspense, useCallback, useMemo, useState } from "react";
import DataTable, { Column, ActionConfig } from "@/uix/datatable";
import DataDialog from "@/uix/datadialog";
import { FieldGroup, FieldLabel } from "@/uix/field";

interface MainData {
  id: string;
  user_id: string;
  company_id: string;
  module_code: string;
  action: string;
  path: string;
  ip_address: string;
  user_agent: string;
  created_at: string;
}
export default function Page() {
  const [mainSelectedRow, setMainSelectedRow] = useState<any>(null);
  const [refresh, setRefresh] = useState(false);
  const loadData = useCallback(async (params: {
    search: string;
    page: number;
    size: number;
    sort_by: string;
    sort_order: 'asc' | 'desc';
  }) => {
    const session = parseSession(window.localStorage.getItem(storageKey));
    if (!session) {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
    try {
      const data = await clientApi<{ data: MainData[]; meta: { total: number; page: number; size: number } }>("/SP03", {
        params: {
          page: String(params.page),
          size: String(params.size),
          sort_by: params.sort_by,
          sort_order: params.sort_order,
          search: params.search,
        },
      });
      return { data: data?.data ?? [], meta: data?.meta ?? { total: 0, page: 1, size: 10 } };
    } catch {
      return { data: [], meta: { total: 0, page: 1, size: 10 } };
    }
  }, []);
  const mainColumns: Column<MainData>[] = useMemo(() => [
    { header: 'ID',         accessor: 'id',          sortable: false, hidden: true  },
    { header: 'User',       accessor: 'user_id',     sortable: false, hidden: true  },
    { header: 'Company',    accessor: 'company_id',  sortable: false, hidden: true  },
    { header: 'Module',     accessor: 'module_code', sortable: true,  hidden: false },
    { header: 'Action',     accessor: 'action',      sortable: true,  hidden: false },
    { header: 'Path',       accessor: 'path',        sortable: false, hidden: true  },
    { header: 'IP Address', accessor: 'ip_address',  sortable: true,  hidden: false },
    { header: 'User Agent', accessor: 'user_agent',  sortable: false, hidden: true  },
    { header: 'Created At', accessor: 'created_at',  sortable: true,  hidden: false,
      formatter: (value: string) => new Date(value).toLocaleString('id-ID'),
    },
  ], []);
  const goMainSearch = () => {
    setRefresh(!refresh);
  };
  const goMainDetail = (row: MainData) => {
    setMainSelectedRow(row);
    setMainModal({
      ...mainModal,
      isOpen: true,
      mode: "detail",
      title: 'Record Details',
      description: null,
      onSubmit: undefined,
    });
  };
  const mainActions: ActionConfig<MainData> = {
    onSearch: goMainSearch,
    onDetail: goMainDetail,
    pickSelect: false,
    hideDetail: false,
    hideCreate: true,
    hideUpdate: true,
    hideDelete: true,
  };
  const [mainModal, setMainModal] = useState({} as any);
  return (
    <Suspense fallback={<p>Loading user profile...</p>}>
      <div className="p-3 max-w-8xl mx-auto">
        <DataTable
          fetchData={loadData}
          columns={mainColumns}
          actions={mainActions}
          hideSearch={false}
          hideSelect={true}
          hidePaging={false}
          hideSort={false}
          hideColumnToggle={true}
          refreshTrigger={refresh}
        />
        <DataDialog
          isOpen={mainModal.isOpen}
          mode={mainModal.mode}
          title={mainModal.title}
          description={mainModal.description}
          onClose={() => setMainModal({ ...mainModal, isOpen: false })}
        >
          {mainSelectedRow && (
            <FieldGroup className="p-3 max-h-[70vh] overflow-y-auto">
              {([
                ["Module", mainSelectedRow.module_code],
                ["Action", mainSelectedRow.action],
                ["Path", mainSelectedRow.path],
                ["IP Address", mainSelectedRow.ip_address],
                ["Created At", mainSelectedRow.created_at ? new Date(mainSelectedRow.created_at).toLocaleString("id-ID") : "-"],
                ["User Agent", mainSelectedRow.user_agent],
              ] as [string, string][]).map(([label, value]) => (
                <div key={label} className="flex flex-col gap-1 border-b px-1 py-1.5 last:border-b-0 sm:flex-row sm:items-start sm:gap-3">
                  <FieldLabel className="w-32 shrink-0 text-sm text-muted-foreground">{label}</FieldLabel>
                  <span className="text-sm break-all">{value || "-"}</span>
                </div>
              ))}
            </FieldGroup>
          )}
        </DataDialog>
      </div>
    </Suspense>
  );
}
