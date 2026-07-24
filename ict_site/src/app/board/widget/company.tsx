"use client";

import { useEffect, useMemo, useState } from "react";
import { useShallow } from "zustand/react/shallow";
import { storageKey, parseSession } from "@/lib/utility";
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/uix/combobox";
import { usePreferencesStore } from "@/app/theme";
import type { CompanyOption } from "../../theme/stores";

type CompanyResponse = {
  data?: { company_id: string; name: string }[];
};

const upsertSessionCompanyName = (companyName: string) => {
  const raw = window.localStorage.getItem(storageKey);
  const session = parseSession(raw || "");
  if (!session) {
    return;
  }
  window.localStorage.setItem(
    storageKey,
    JSON.stringify({
      ...session,
      user_profile: {
        ...session.user_profile,
        company_name: companyName,
      },
    }),
  );
};

export function CompanyCombobox() {
  const { companyId, companyList, setCompanyId, setCompanyList } = usePreferencesStore(
    useShallow((s) => ({
      companyId: s.companyId,
      companyList: s.companyList,
      setCompanyId: s.setCompanyId,
      setCompanyList: s.setCompanyList,
    })),
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lockedCompanyId, setLockedCompanyId] = useState("");
  useEffect(() => {
    let isDisposed = false;
    setLoading(true);
    setError(null);
    const session = parseSession(window.localStorage.getItem(storageKey) || "");
    if (!session) {
      setLockedCompanyId("");
      setCompanyId("");
      setCompanyList([]);
      setLoading(false);
      return;
    }
    const sessionCompanyId = (session.user_profile.company_id || "").trim();
    setLockedCompanyId(sessionCompanyId);
    const initialCompany = (sessionCompanyId || companyId || "").trim();
    if (initialCompany) {
      setCompanyId(initialCompany);
    }
      fetch(`/proxy/pages/SP01/company`, {
        headers: {
          Authorization: `Bearer ${session.token}`,
        },
      })
      .then((res) => res.json())
      .then((response: CompanyResponse) => {
        if (isDisposed) {
          return;
        }
        const raw = Array.isArray(response.data) ? response.data : [];
        const nextCompanies: CompanyOption[] = raw.map((c) => ({
          id: c.company_id,
          code: c.company_id,
          name: c.name,
        }));
        setCompanyList(nextCompanies);
        if (!nextCompanies.length) {
          return;
        }
        if (sessionCompanyId) {
          const matchedLockedCompany = nextCompanies.find(
            (company) => String(company.id).trim() === sessionCompanyId,
          );
          setCompanyId(matchedLockedCompany?.id ?? "");
          upsertSessionCompanyName(matchedLockedCompany?.name ?? "");
          return;
        }
        const normalizedInitialCompany = String(initialCompany).trim();
        const matchedInitialCompany = nextCompanies.find(
          (company) => String(company.id).trim() === normalizedInitialCompany,
        );
        if (normalizedInitialCompany && !matchedInitialCompany) {
          setCompanyId("");
          upsertSessionCompanyName("");
        }
        if (matchedInitialCompany && companyId !== matchedInitialCompany.id) {
          setCompanyId(matchedInitialCompany.id);
        }
        if (matchedInitialCompany) {
          upsertSessionCompanyName(matchedInitialCompany.name);
        }
      })
      .catch(() => {
        if (isDisposed) {
          return;
        }
        setError("Failed to load company list.");
      })
      .finally(() => {
        if (isDisposed) {
          return;
        }
        setLoading(false);
      });
    return () => {
      isDisposed = true;
    };
  }, [setCompanyId, setCompanyList]);
  const placeholder = useMemo(() => {
    if (loading) {
      return "Loading company...";
    }
    return "Select company";
  }, [loading]);
  const selectedCompanyId = lockedCompanyId || companyId;
  return (
    <div className="px-2 pt-2 pb-1 group-data-[collapsible=icon]:hidden">
      <Combobox
        key={`${selectedCompanyId}-${companyList.length}`}
        value={selectedCompanyId || ""}
        onValueChange={(value) => {
          if (lockedCompanyId) {
            return;
          }
          if (value !== null && value !== undefined) {
            const id = String(value).trim();
            if (!id) {
              setCompanyId("");
              upsertSessionCompanyName("");
              return;
            }
            setCompanyId(id);
            const selected = companyList.find((company) => company.id === id);
            upsertSessionCompanyName(selected?.name ?? "");
          }
        }}
        itemToStringLabel={(item) => {
          const id = typeof item === "string" ? item : String(item ?? "");
          if (!id) {
            return "- Select Company -";
          }
          const selected = companyList.find((company) => String(company.id).trim() === id.trim());
          if (selected?.name) {
            return selected.name;
          }
          return loading ? id : "- Select Company -";
        }}
      >
        <ComboboxInput
          placeholder={placeholder}
          disabled={Boolean(lockedCompanyId) || loading || !companyList.length}
          showClear={false}
        />
        <ComboboxContent>
          <ComboboxEmpty>
            {error ?? (companyList.length === 0 ? "No company available." : "")}
          </ComboboxEmpty>
          <ComboboxList>
            {!Boolean(lockedCompanyId) ? <ComboboxItem value="">- Select Company -</ComboboxItem> : null}
            {companyList.map((company) => (
              <ComboboxItem key={company.id} value={company.id}>
                {company.name}
              </ComboboxItem>
            ))}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
    </div>
  );
}
