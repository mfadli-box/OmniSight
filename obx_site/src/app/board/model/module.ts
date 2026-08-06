"use client";

import {
  type LucideIcon,
  Activity, Container, Cpu, FileText, Globe, ReceiptText,
  Settings, ServerCog, Shield, Users,
} from "lucide-react";
import { useEffect, useState } from "react";
import { usePreferencesStore } from "@/app/theme";
import { parseSession, storageKey } from "@/lib/utility";

export type NavBadge = "new" | "soon";

export interface NavSubItem {
  id: string;
  title: string;
  url: string;
  icon?: LucideIcon;
  badge?: NavBadge;
  disabled?: boolean;
  newTab?: boolean;
}
interface NavItemBase {
  id: string;
  title: string;
  icon?: LucideIcon;
  badge?: NavBadge;
  disabled?: boolean;
  newTab?: boolean;
}
export interface NavMainLinkItem extends NavItemBase {
  url: string;
  subItems?: never;
}
export interface NavMainParentItem extends NavItemBase {
  subItems: NavSubItem[];
}
export type NavMainItem = NavMainLinkItem | NavMainParentItem;
export interface NavGroup {
  id: number;
  label?: string;
  items: NavMainItem[];
}

type ModuleTreeNode = {
  id: string;
  code: string;
  name: string;
  path: string;
  is_page: boolean;
  children?: ModuleTreeNode[];
};

const staticAdminItem: NavMainParentItem = {
  id: "SM",
  title: "System Manager",
  icon: Settings,
  subItems: [
    { id: "SM01", title: "User", url: "/board/pages/SM01", newTab: false },
    { id: "SM02", title: "Module", url: "/board/pages/SM02", newTab: false },
    { id: "SM03", title: "Company", url: "/board/pages/SM03", newTab: false },
    { id: "SM04", title: "Signature", url: "/board/pages/SM04", newTab: false },
    { id: "SM05", title: "Session", url: "/board/pages/SM05", newTab: false },
  ],
};

const staticSessionItem: NavMainParentItem = {
  id: "SP",
  title: "Session Profile",
  icon: Users,
  subItems: [
    { id: "SP01", title: "Profile", url: "/board/pages/SP01", newTab: false },
    { id: "SP02", title: "Password", url: "/board/pages/SP02", newTab: false },
    { id: "SP03", title: "History", url: "/board/pages/SP03", newTab: false },
  ],
};

export const moduleItem: NavGroup[] = [
  {
    id: 1,
    label: "",
    items: [],
  },
];

function mapIconByCode(code: string): LucideIcon | undefined {
  const key = code.slice(0, 2).toUpperCase();
  if (key === "SP") return Users;
  if (key === "SM") return Settings;
  if (key === "DK") return Container;
  if (key === "DM") return FileText;
  if (key === "IM") return Activity;
  if (key === "NX") return Globe;
  if (key === "PL") return ServerCog;
  if (key === "SC") return Shield;
  if (key === "VM") return Cpu;
  return ReceiptText;
}

function collectSubItems(nodes: ModuleTreeNode[]): NavSubItem[] {
  const out: NavSubItem[] = [];
  for (const node of nodes) {
    if (node.path) {
      out.push({
        id: node.code,
        title: node.name,
        url: node.path,
        newTab: !node.is_page,
      });
    }
    if (Array.isArray(node.children) && node.children.length > 0) {
      out.push(...collectSubItems(node.children));
    }
  }
  return out;
}

function toNavItem(node: ModuleTreeNode): NavMainItem | null {
  const icon = mapIconByCode(node.code);
  if (node.is_page) {
    if (!node.path) return null;
    return { id: node.code, title: node.name, icon, url: node.path, newTab: false };
  }
  const subItems = Array.isArray(node.children) ? collectSubItems(node.children) : [];
  return { id: node.code, title: node.name, icon, subItems };
}

function dedupeById(items: NavMainItem[]): NavMainItem[] {
  const map = new Map<string, NavMainItem>();
  for (const item of items) {
    map.set(item.id, item);
  }
  return Array.from(map.values());
}

function buildNavGroups(nodes: ModuleTreeNode[], isLoggedIn: boolean, isAdmin: boolean): NavGroup[] {
  const items: NavMainItem[] = [];
  for (const root of nodes) {
    const navItem = toNavItem(root);
    if (navItem) items.push(navItem);
  }
  if (isAdmin) items.push(staticAdminItem);
  if (isLoggedIn) items.push(staticSessionItem);
  return [
    {
      id: 1,
      label: "",
      items: dedupeById(items),
    },
  ];
}

export function useModuleItem(): NavGroup[] {
  const companyId = usePreferencesStore((s) => s.companyId);
  const [items, setItems] = useState<NavGroup[]>(moduleItem);
  const [isHydrated, setIsHydrated] = useState(false);
  useEffect(() => {
    setIsHydrated(true);
  }, []);
  useEffect(() => {
    if (!isHydrated) return;
    let active = true;
    const load = async () => {
      const session = parseSession(window.localStorage.getItem(storageKey));
      const isLoggedIn = Boolean(session?.token);
      const isAdmin = Boolean(session?.user_profile.is_admin);
      if (!isLoggedIn || !session) {
        if (active) setItems(buildNavGroups([], false, false));
        return;
      }
      try {
        const response = await fetch(`/proxy/pages/SP01/module?company_id=${companyId}`, {
          method: "GET",
          headers: { Authorization: `Bearer ${session.token}` },
        });
        const data = await response.json();
        const tree = Array.isArray(data?.data) ? data.data : [];
        if (active) setItems(buildNavGroups(tree, true, isAdmin));
      } catch {
        if (active) setItems(buildNavGroups([], true, isAdmin));
      }
    };
    void load();
    return () => { active = false; };
  }, [companyId]);
  return items;
}
