"use client";

import {
  Activity, AlertTriangle, Award, BarChart3, CheckCircle, ClipboardCheck,
  ClipboardList, Container, CreditCard, FileCheck, FileText, Globe, GlobeLock,
  Key, LayoutDashboard, MapPin, Network, Radio, ReceiptText, Router, ScrollText,
  Settings, Shield, ShieldCheck, Users, Wrench,
  type LucideIcon,
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
    { id: "SM06", title: "Location", url: "/board/pages/SM06", newTab: false },
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
  if (key === "AM") return MapPin;
  if (key === "VL") return Shield;
  if (key === "IN") return AlertTriangle;
  if (key === "SI") return Activity;
  if (key === "RK") return BarChart3;
  if (key === "CP") return CheckCircle;
  if (key === "CA") return AlertTriangle;
  if (key === "IA") return ClipboardCheck;
  if (key === "MR") return Users;
  if (key === "QM") return Award;
  if (key === "CT") return Key;
  if (key === "DK") return Container;
  if (key === "DM") return FileText;
  if (key === "SR") return ClipboardList;
  if (key === "BL") return CreditCard;
  if (key === "PM") return Wrench;
  if (key === "DB") return LayoutDashboard;
  if (key === "XX") return Wrench;
  if (key === "NK") return Network;
  if (key === "MK") return Router;
  if (key === "DN") return Globe;
  if (key === "IP") return GlobeLock;
  if (key === "MO") return Radio;
  if (key === "MT") return Activity;
  if (key === "ML") return ScrollText;
  if (key === "MF") return FileCheck;
  if (key === "MS") return ShieldCheck;
  if (key === "MD") return Container;
  if (key === "IS") return Award;
  if (key === "WM") return Globe;
  if (key === "WS") return Shield;
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
  useEffect(() => {
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
